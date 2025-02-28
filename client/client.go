package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)


const Status441EntityNotFound = 441
const Status442EntityAlreadyExists = 441


type HttpError struct {
	StatusCode int
	Error      error
}

type Client struct {
	// AccessTokenScope           string
	// AccessToken                azcore.AccessToken
	AccessToken                string
	Url                        string
	InsecureSkipVerify         bool
	HealthCheckIntervalSeconds int
	HealthCheckTimeoutSeconds  int
}

var (
	IsHealthy       = false
	IsTimedOut      = false
	mu              sync.Mutex
	healthCheckOnce sync.Once
)

// Initialize bearer token once and reuse
func CreateClient(url string, insecure_skip_verify bool, healthCheckIntervalSeconds int, healthCheckTimeoutSeconds int, accessToken string, clientId string, clientSecret string) (*Client, error) {

	if accessToken == "" {
		tokenResp, err := getAccessToken(url, clientId, clientSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to get access token: %w", err)
		}
		accessToken = tokenResp.AccessToken
	}

	client := &Client{
		Url:                        url,
		InsecureSkipVerify:         insecure_skip_verify,
		HealthCheckIntervalSeconds: healthCheckIntervalSeconds,
		HealthCheckTimeoutSeconds:  healthCheckTimeoutSeconds,
		AccessToken:                accessToken,
	}

	// Start the background health check loop
	healthCheckOnce.Do(func() {
		go client.startHealthCheckLoop()
	})

	return client, nil
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func getAccessToken(serverUrl string, clientId string, clientSecret string) (*TokenResponse, *HttpError) {
	tokenURL := fmt.Sprintf("%s/connect/token", serverUrl)

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "snapcd_scope")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Create a new request with form-encoded data
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, &HttpError{StatusCode: 0, Error: err}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request with context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, &HttpError{StatusCode: 0, Error: err}
	}
	defer resp.Body.Close()

	// Ensure success status code
	if resp.StatusCode != http.StatusOK {
		return nil, &HttpError{StatusCode: 0, Error: fmt.Errorf("unexpected status code: %d", resp.StatusCode)}
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &HttpError{StatusCode: 0, Error: err}
	}

	// Parse JSON response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, &HttpError{StatusCode: 0, Error: err}
	}

	return &tokenResp, nil
}

// func getAccessToken(ctx context.Context, scope string) (azcore.AccessToken, error) {
// 	credentialOptions := azidentity.DefaultAzureCredentialOptions{}
// 	var nullToken azcore.AccessToken

// 	credential, err := azidentity.NewDefaultAzureCredential(&credentialOptions)
// 	if err != nil {
// 		return nullToken, err
// 	}

// 	token, err := credential.GetToken(ctx, policy.TokenRequestOptions{
// 		Scopes: []string{scope},
// 	})
// 	if err != nil {
// 		return nullToken, err
// 	}

// 	return azcore.AccessToken{
// 		Token:     token.Token,
// 		ExpiresOn: token.ExpiresOn,
// 	}, nil
// }

func (client *Client) makeRequest(method string, path string, body []byte) (map[string]interface{}, *HttpError) {
	ticker := time.Tick(1 * time.Second)

	// if IsHealthy, immediately go to request
	mu.Lock()
	if IsHealthy {
		mu.Unlock()
		goto Request
	}
	mu.Unlock()

	// Otherwise start a loop continues until either IsHealthy (continue to request) or IsTimedOut (return error)
	for {
		select {
		case <-ticker:
			mu.Lock()
			if IsHealthy {
				mu.Unlock()
				goto Request
			}
			if IsTimedOut {
				mu.Unlock()
				return nil, &HttpError{StatusCode: 0, Error: errors.New("Health check timed out after " + strconv.Itoa(client.HealthCheckTimeoutSeconds) + " seconds. Service is not healthy.")}
			}
			mu.Unlock()
		}
	}

Request:
	req, err := http.NewRequest(method, client.Url+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, &HttpError{StatusCode: 441, Error: errors.New("Status441EntityNotFound")}
	}

	req.Header.Set("Authorization", "Bearer "+client.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: client.InsecureSkipVerify},
	}
	httpClient := &http.Client{Timeout: 120 * time.Second, Transport: tr}

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, &HttpError{StatusCode: 0, Error: err}
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &HttpError{StatusCode: 0, Error: err}
	}

	if response.StatusCode == 441 {
		return nil, &HttpError{StatusCode: 441, Error: errors.New("Status441EntityNotFound")}
	}

	if response.StatusCode == 442 {
		return nil, &HttpError{StatusCode: 442, Error: errors.New("Status442EntityAlreadyExists")}
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, &HttpError{StatusCode: 441, Error: fmt.Errorf("unexpected status code: %d (%s). Response body: %s", response.StatusCode, http.StatusText(response.StatusCode), string(responseBody))}
	}

	var result map[string]interface{}

	if len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, &result)
		if err != nil {
			return nil, &HttpError{StatusCode: 0, Error: err}
		}
	}

	return result, nil
}

func (client *Client) Post(path string, data interface{}) (map[string]interface{}, *HttpError) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, &HttpError{StatusCode: 0, Error: err}
	}

	return client.makeRequest(http.MethodPost, path, body)
}

func (client *Client) Get(path string) (map[string]interface{}, *HttpError) {
	return client.makeRequest(http.MethodGet, path, nil)
}

func (client *Client) Put(path string, data interface{}) (map[string]interface{}, *HttpError) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, &HttpError{StatusCode: 0, Error: err}
	}

	return client.makeRequest(http.MethodPut, path, body)
}

func (client *Client) Delete(path string) (map[string]interface{}, *HttpError) {
	return client.makeRequest(http.MethodDelete, path, nil)
}

func (client *Client) checkHealth() bool {
	req, err := http.NewRequest(http.MethodGet, client.Url+"/health", nil)
	if err != nil {
		return false
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: client.InsecureSkipVerify},
	}
	httpClient := &http.Client{Timeout: time.Duration(5) * time.Second, Transport: tr}

	response, err := httpClient.Do(req)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusOK
}

func (client *Client) startHealthCheckLoop() {

	// poll the /health endpoint and immediately return if it is healthy
	if client.checkHealth() {
		mu.Lock()
		IsHealthy = true
		mu.Unlock()
		return
	}

	// Otherwise start a loop that polls /health every client.HealthCheckIntervalSeconds.

	timeout := time.After(time.Duration(client.HealthCheckTimeoutSeconds) * time.Second)
	ticker := time.Tick(time.Duration(client.HealthCheckIntervalSeconds) * time.Second)

	for {
		select {
		case <-timeout:
			mu.Lock()
			IsTimedOut = true
			mu.Unlock()
			return
		case <-ticker:
			if client.checkHealth() {
				mu.Lock()
				IsHealthy = true
				mu.Unlock()
				return
			}
		}
	}
}

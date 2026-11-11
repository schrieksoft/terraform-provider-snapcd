package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// CliToken is what `snapcd account get-access-token --output json` prints: a user's access token
// and the Server and organization of the CLI context it belongs to.
type CliToken struct {
	AccessToken        string    `json:"accessToken"`
	ExpiresOn          time.Time `json:"expiresOn"`
	Url                string    `json:"url"`
	OrganizationId     string    `json:"organizationId"`
	InsecureSkipVerify bool      `json:"insecureSkipVerify"`
}

// The CLI is asked again this long before its token expires, so a request never carries one that
// lapses in flight.
const cliTokenRenewBefore = 5 * time.Minute

// GetCliToken asks the snapcd CLI for an access token, the way the azurerm provider asks the Azure
// CLI. SNAPCD_CLI_PATH names the binary when it is not `snapcd` on the PATH, and SNAPCD_CLI_CONTEXT
// picks a context other than the current one. The CLI never prompts, so this cannot hang.
func GetCliToken(ctx context.Context) (*CliToken, error) {
	binary := os.Getenv("SNAPCD_CLI_PATH")
	if binary == "" {
		binary = "snapcd"
	}

	args := []string{"account", "get-access-token", "--output", "json"}
	if cliContext := os.Getenv("SNAPCD_CLI_CONTEXT"); cliContext != "" {
		args = append(args, "--context", cliContext)
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) || errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("no credentials are configured and the snapcd CLI (%s) was not found: set client_id with client_secret or federated_token_file, set access_token, or install the CLI and run `snapcd login`", binary)
		}
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return nil, fmt.Errorf("the snapcd CLI could not provide an access token: %s", message)
	}

	var token CliToken
	if err := json.Unmarshal(stdout.Bytes(), &token); err != nil {
		return nil, fmt.Errorf("the snapcd CLI returned output that is not a token: %v", err)
	}
	if token.AccessToken == "" {
		return nil, errors.New("the snapcd CLI returned no access token")
	}
	return &token, nil
}

// cliTokenSource keeps the CLI's token and asks for a new one before it expires.
type cliTokenSource struct {
	mu      sync.Mutex
	current *CliToken
	fetch   func(context.Context) (*CliToken, error)
	now     func() time.Time
}

func (s *cliTokenSource) token() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current != nil && s.now().Add(cliTokenRenewBefore).Before(s.current.ExpiresOn) {
		return s.current.AccessToken, nil
	}

	next, err := s.fetch(context.Background())
	if err != nil {
		return "", err
	}
	s.current = next
	return next.AccessToken, nil
}

// newCliTokenSource starts from a token already fetched, so the provider could read the Server and
// organization from it before the client existed.
func newCliTokenSource(initial *CliToken) *cliTokenSource {
	return &cliTokenSource{current: initial, fetch: GetCliToken, now: time.Now}
}

func (client *Client) bearerToken() (string, error) {
	if client.cli == nil {
		return client.AccessToken, nil
	}
	return client.cli.token()
}

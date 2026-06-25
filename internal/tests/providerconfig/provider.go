// SPDX-License-Identifier: MPL-2.0

package providerconfig

import (
	"fmt"
	"math/rand"
	"os"

	provider "terraform-provider-snapcd/internal/provider"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ProviderConfig() string {
	url := getEnvOrDefault("SNAPCD_URL", "https://localhost:20002")
	insecureSkipVerify := getEnvOrDefault("SNAPCD_INSECURE_SKIP_VERIFY", "true")
	clientID := getEnvOrDefault("SNAPCD_CLIENT_ID", "debugTerraformer")
	clientSecret := getEnvOrDefault("SNAPCD_CLIENT_SECRET", "debugTerraformer")
	organizationID := getEnvOrDefault("SNAPCD_ORGANIZATION_ID", "10000000-0000-0000-0000-000000000000")

	return fmt.Sprintf(`
provider "snapcd" {
	url                           = %q
	insecure_skip_verify          = %s
	health_check_interval_seconds = 5
	health_check_timeout_seconds  = 20
	client_id                     = %q
	client_secret                 = %q
	organization_id               = %q
}
`, url, insecureSkipVerify, clientID, clientSecret, organizationID)
}

var (
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"snapcd": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
	NamePostfix = GenerateRandomString(16)
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

func AppendRandomString(input string) string {
	return fmt.Sprintf(input, NamePostfix)
}

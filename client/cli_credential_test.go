package client

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func fakeCli(t *testing.T, script string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the stand-in CLI is a shell script")
	}
	path := filepath.Join(t.TempDir(), "snapcd")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SNAPCD_CLI_PATH", path)
}

func TestGetCliToken_ReadsTheTokenAndItsContext(t *testing.T) {
	fakeCli(t, `[ "$*" = "account get-access-token --output json --context staging" ] || { echo "unexpected args: $*" >&2; exit 2; }
echo '{"accessToken":"tok","expiresOn":"2026-11-11T10:00:00Z","url":"https://snapcd.example","organizationId":"10000000-0000-0000-0000-000000000000","insecureSkipVerify":true}'`)
	t.Setenv("SNAPCD_CLI_CONTEXT", "staging")

	token, err := GetCliToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "tok" || token.Url != "https://snapcd.example" ||
		token.OrganizationId != "10000000-0000-0000-0000-000000000000" || !token.InsecureSkipVerify {
		t.Fatalf("unexpected token: %+v", token)
	}
	if !token.ExpiresOn.Equal(time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)) {
		t.Fatalf("expiresOn = %v", token.ExpiresOn)
	}
}

func TestGetCliToken_ReportsTheCliError(t *testing.T) {
	fakeCli(t, `echo "Error: Not logged in to context 'prod'. Run: snapcd login --context prod" >&2; exit 1`)

	_, err := GetCliToken(context.Background())
	if err == nil || !strings.Contains(err.Error(), "snapcd login --context prod") {
		t.Fatalf("err = %v", err)
	}
}

func TestGetCliToken_NamesTheAlternativesWhenTheCliIsMissing(t *testing.T) {
	t.Setenv("SNAPCD_CLI_PATH", filepath.Join(t.TempDir(), "no-such-snapcd"))

	_, err := GetCliToken(context.Background())
	if err == nil || !strings.Contains(err.Error(), "snapcd login") || !strings.Contains(err.Error(), "client_secret") {
		t.Fatalf("err = %v", err)
	}
}

func TestCliTokenSource_AsksAgainShortlyBeforeExpiry(t *testing.T) {
	now := time.Date(2026, 11, 11, 9, 0, 0, 0, time.UTC)
	calls := 0
	source := &cliTokenSource{
		current: &CliToken{AccessToken: "first", ExpiresOn: now.Add(30 * time.Minute)},
		fetch: func(context.Context) (*CliToken, error) {
			calls++
			return &CliToken{AccessToken: "second", ExpiresOn: now.Add(time.Hour)}, nil
		},
		now: func() time.Time { return now },
	}

	if token, _ := source.token(); token != "first" || calls != 0 {
		t.Fatalf("token = %s, calls = %d", token, calls)
	}

	now = now.Add(26 * time.Minute)
	if token, _ := source.token(); token != "second" || calls != 1 {
		t.Fatalf("token = %s, calls = %d", token, calls)
	}

	if token, _ := source.token(); token != "second" || calls != 1 {
		t.Fatalf("token = %s, calls = %d", token, calls)
	}
}

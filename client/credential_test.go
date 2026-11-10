package client

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClientCredential_SecretWhenNoTokenFile(t *testing.T) {
	data, err := clientCredential("s3cret", "")
	if err != nil {
		t.Fatal(err)
	}
	if data.Get("client_secret") != "s3cret" || data.Has("snapcd_federated_assertion") {
		t.Fatalf("unexpected form: %v", data)
	}
}

func TestClientCredential_SecretTakesPrecedenceOverTokenFile(t *testing.T) {
	data, err := clientCredential("s3cret", "/does/not/exist")
	if err != nil {
		t.Fatal(err)
	}
	if data.Get("client_secret") != "s3cret" || data.Has("snapcd_federated_assertion") {
		t.Fatalf("unexpected form: %v", data)
	}
}

func TestClientCredential_TokenFileIsReadEachCall(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(path, []byte("first\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	data, err := clientCredential("", path)
	if err != nil {
		t.Fatal(err)
	}
	if data.Has("client_secret") {
		t.Fatalf("client_secret sent alongside a federated token: %v", data)
	}
	if data.Get("snapcd_federated_assertion") != "first" {
		t.Fatalf("assertion = %q", data.Get("snapcd_federated_assertion"))
	}
	if data.Get("snapcd_federated_assertion_type") != "urn:ietf:params:oauth:client-assertion-type:jwt-bearer" {
		t.Fatalf("assertion type = %q", data.Get("snapcd_federated_assertion_type"))
	}

	if err := os.WriteFile(path, []byte("second"), 0o600); err != nil {
		t.Fatal(err)
	}
	data, err = clientCredential("", path)
	if err != nil {
		t.Fatal(err)
	}
	if data.Get("snapcd_federated_assertion") != "second" {
		t.Fatalf("assertion = %q", data.Get("snapcd_federated_assertion"))
	}
}

func TestClientCredential_MissingOrEmptyTokenFileFails(t *testing.T) {
	dir := t.TempDir()
	if _, err := clientCredential("", filepath.Join(dir, "missing")); err == nil {
		t.Fatal("expected an error for a missing file")
	}
	empty := filepath.Join(dir, "empty")
	if err := os.WriteFile(empty, []byte("  \n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := clientCredential("", empty); err == nil {
		t.Fatal("expected an error for an empty file")
	}
}

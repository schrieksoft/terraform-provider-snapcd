package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type holdCall struct {
	Path string
	Body map[string]interface{}
}

type holdServer struct {
	mu        sync.Mutex
	calls     []holdCall
	status    int
	body      string
	delay     time.Duration
	effective int
	warning   string
}

func (s *holdServer) handler(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if s.delay > 0 {
		time.Sleep(s.delay)
	}
	s.mu.Lock()
	s.calls = append(s.calls, holdCall{Path: r.URL.Path, Body: body})
	status, text, effective, warning := s.status, s.body, s.effective, s.warning
	s.mu.Unlock()

	if status != 0 {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(text))
		return
	}
	if strings.HasSuffix(r.URL.Path, "/release") {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"effectiveHeartbeatIntervalSeconds": effective, "effectiveExpirySeconds": 20, "warning": warning, "notHeldModuleIds": []string{}})
}

func (s *holdServer) setStatus(status int) {
	s.mu.Lock()
	s.status = status
	s.mu.Unlock()
}

func (s *holdServer) callsTo(suffix string) []holdCall {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []holdCall
	for _, c := range s.calls {
		if strings.HasSuffix(c.Path, suffix) {
			out = append(out, c)
		}
	}
	return out
}

func moduleIds(c holdCall) []string {
	raw, _ := c.Body["ModuleIds"].([]interface{})
	ids := make([]string, 0, len(raw))
	for _, v := range raw {
		ids = append(ids, v.(string))
	}
	return ids
}

func newHoldTestClient(t *testing.T, settings ModuleHoldSettings) (*Client, *holdServer) {
	t.Helper()
	mu.Lock()
	IsHealthy = true
	mu.Unlock()

	server := &holdServer{}
	srv := httptest.NewServer(http.HandlerFunc(server.handler))
	client := &Client{Url: srv.URL, OrganizationId: "org", AccessToken: "token"}
	client.holder = newModuleHolder(client, settings)
	t.Cleanup(func() {
		client.holder.release(time.Second)
		srv.Close()
	})
	return client, server
}

func enabled(interval int) ModuleHoldSettings {
	return ModuleHoldSettings{Enabled: true, HeartbeatIntervalSeconds: &interval}
}

func holdField(fields map[string]interface{}) map[string]interface{} {
	hold, _ := fields["Hold"].(map[string]interface{})
	return hold
}

func TestHoldModuleSendsOneHeartbeatPerModule(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))

	client.HoldModule("a")
	client.HoldModule("a")

	calls := server.callsTo("/ModuleHold/heartbeat")
	if len(calls) != 1 {
		t.Fatalf("expected 1 heartbeat, got %d", len(calls))
	}
	if ids := moduleIds(calls[0]); len(ids) != 1 || ids[0] != "a" {
		t.Fatalf("expected heartbeat for [a], got %v", ids)
	}
	if calls[0].Path != "/api/org/ModuleHold/heartbeat" {
		t.Fatalf("unexpected path %s", calls[0].Path)
	}
	if calls[0].Body["HeartbeatIntervalSeconds"].(float64) != 60 {
		t.Fatalf("expected the heartbeat interval in the request, got %v", calls[0].Body)
	}
}

func TestParallelWritesForOneModuleWaitForTheFirstHeartbeat(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))
	server.delay = 200 * time.Millisecond

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.HoldModule("a")
		}()
	}
	wg.Wait()

	if n := len(server.callsTo("/heartbeat")); n != 1 {
		t.Fatalf("expected 1 heartbeat, got %d", n)
	}
}

func TestAHeartbeatCarriesTheModulesAlreadyHeld(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))

	client.HoldModule("a")
	client.HoldModule("b")

	calls := server.callsTo("/heartbeat")
	if ids := moduleIds(calls[1]); len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Fatalf("expected second heartbeat for [a b], got %v", ids)
	}
}

func TestAFailedHoldIsRetriedOnTheNextWrite(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))
	server.setStatus(http.StatusInternalServerError)

	client.HoldModule("a")
	server.setStatus(0)
	client.HoldModule("b")

	calls := server.callsTo("/heartbeat")
	if ids := moduleIds(calls[1]); len(ids) != 1 || ids[0] != "b" {
		t.Fatalf("a failed hold must not be sent as held, got %v", ids)
	}

	client.HoldModule("a")
	calls = server.callsTo("/heartbeat")
	if len(calls) != 3 {
		t.Fatalf("expected the failed module to be held on its next write, got %d heartbeats", len(calls))
	}
}

func TestHoldingDisabledSendsNothing(t *testing.T) {
	client, server := newHoldTestClient(t, ModuleHoldSettings{Enabled: false})

	client.HoldModule("a")
	client.ModuleCreated("b")
	if fields, _ := client.ModuleCreateHoldFields(); fields != nil {
		t.Fatalf("expected no create fields, got %v", fields)
	}
	time.Sleep(1500 * time.Millisecond)
	client.holder.release(time.Second)

	if n := len(server.callsTo("")); n != 0 {
		t.Fatalf("expected no requests, got %d", n)
	}
}

func TestCreateFieldsCarryTheIntervalAndBufferInTheHoldObject(t *testing.T) {
	buffer := 5
	interval := 10
	client, _ := newHoldTestClient(t, ModuleHoldSettings{Enabled: true, HeartbeatIntervalSeconds: &interval, HeartbeatBufferSeconds: &buffer})

	fields, _ := client.ModuleCreateHoldFields()
	if hold := holdField(fields); hold["HeartbeatIntervalSeconds"] != 10 || hold["HeartbeatBufferSeconds"] != 5 {
		t.Fatalf("unexpected create fields %v", fields)
	}
}

func TestACreatedModuleIsKeptHeldByTheHeartbeatAndReleasedAtTheEnd(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(1))

	client.ModuleCreated("a")
	if n := len(server.callsTo("/heartbeat")); n != 0 {
		t.Fatalf("a module held on creation needs no immediate heartbeat, got %d", n)
	}

	time.Sleep(1500 * time.Millisecond)
	beats := server.callsTo("/heartbeat")
	if len(beats) == 0 {
		t.Fatalf("expected a heartbeat within one interval")
	}
	if ids := moduleIds(beats[0]); len(ids) != 1 || ids[0] != "a" {
		t.Fatalf("expected heartbeat for [a], got %v", ids)
	}

	client.holder.release(time.Second)
	releases := server.callsTo("/ModuleHold/release")
	if len(releases) != 1 || len(moduleIds(releases[0])) != 1 {
		t.Fatalf("expected one release for [a], got %v", releases)
	}

	after := len(server.callsTo("/heartbeat"))
	time.Sleep(1500 * time.Millisecond)
	if n := len(server.callsTo("/heartbeat")); n != after {
		t.Fatalf("heartbeats continued after release")
	}
}

func TestTheServersHeartbeatIntervalIsUsedAndItsWarningShownOnce(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))
	server.effective = 1
	server.warning = "held with a heartbeat every 1 seconds"

	if w := client.HoldModule("a"); w != "held with a heartbeat every 1 seconds" {
		t.Fatalf("expected the server's warning on the first hold, got %q", w)
	}
	if w := client.HoldModule("b"); w != "" {
		t.Fatalf("expected the warning only once, got %q", w)
	}
	if _, w := client.ModuleCreateHoldFields(); w != "" {
		t.Fatalf("expected the warning only once, got %q", w)
	}

	time.Sleep(1500 * time.Millisecond)
	if n := len(server.callsTo("/heartbeat")); n < 3 {
		t.Fatalf("expected the loop to beat at the server's 1-second interval, got %d heartbeats", n)
	}
}

func TestAModuleCreateLearnsTheServersSettingsFirst(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))
	server.effective = 25
	server.warning = "adjusted"

	fields, w := client.ModuleCreateHoldFields()
	if holdField(fields)["HeartbeatIntervalSeconds"] != 60 {
		t.Fatalf("expected the configured interval in the create fields, got %v", fields)
	}
	if client.holder.currentInterval() != 25*time.Second {
		t.Fatalf("expected the heartbeat to use the server's interval, got %v", client.holder.currentInterval())
	}
	if w != "adjusted" {
		t.Fatalf("expected the warning with the create fields, got %q", w)
	}
	calls := server.callsTo("/heartbeat")
	if len(calls) != 1 || len(moduleIds(calls[0])) != 0 {
		t.Fatalf("expected one heartbeat with no modules to learn the settings, got %v", calls)
	}

	client.ModuleCreated("a")
	client.ModuleCreateHoldFields()
	if len(server.callsTo("/heartbeat")) != 1 {
		t.Fatalf("expected the settings to be learned only once")
	}
}

func TestHoldFailuresDoNotFailTheWrite(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))
	server.setStatus(http.StatusInternalServerError)

	if w := client.HoldModule("a"); w != "" {
		t.Fatalf("expected no warning, got %q", w)
	}
	if fields, _ := client.ModuleCreateHoldFields(); holdField(fields)["HeartbeatIntervalSeconds"] != 60 {
		t.Fatalf("expected the configured interval when the server cannot be reached, got %v", fields)
	}
}

func TestUnsetTimingIsLeftForTheServerToDefault(t *testing.T) {
	client, server := newHoldTestClient(t, ModuleHoldSettings{Enabled: true})
	server.effective = 7

	fields, _ := client.ModuleCreateHoldFields()
	hold := holdField(fields)
	if hold == nil || len(hold) != 0 {
		t.Fatalf("expected an empty hold object, got %v", fields)
	}
	body := server.callsTo("/heartbeat")[0].Body
	if _, ok := body["HeartbeatIntervalSeconds"]; ok {
		t.Fatalf("expected no heartbeat interval in the request, got %v", body)
	}
	if _, ok := body["HeartbeatBufferSeconds"]; ok {
		t.Fatalf("expected no buffer in the request, got %v", body)
	}
	if client.holder.currentInterval() != 7*time.Second {
		t.Fatalf("expected the server's interval, got %v", client.holder.currentInterval())
	}
}

func idsIn(c holdCall, list string) []string {
	raw, _ := c.Body[list].([]interface{})
	ids := make([]string, 0, len(raw))
	for _, v := range raw {
		ids = append(ids, v.(string))
	}
	return ids
}

func TestOneHeartbeatCarriesModulesNamespacesAndStacks(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))

	client.HoldModule("m")
	client.HoldNamespace("n")
	client.HoldStack("s")
	client.HoldNamespace("n")

	calls := server.callsTo("/heartbeat")
	if len(calls) != 3 {
		t.Fatalf("expected one heartbeat per new hold, got %d", len(calls))
	}
	last := calls[2]
	if m, n, s := idsIn(last, "ModuleIds"), idsIn(last, "NamespaceIds"), idsIn(last, "StackIds"); len(m) != 1 || len(n) != 1 || len(s) != 1 {
		t.Fatalf("expected the last heartbeat to carry all three, got modules %v namespaces %v stacks %v", m, n, s)
	}

	client.holder.release(time.Second)
	release := server.callsTo("/release")[0]
	if len(idsIn(release, "ModuleIds")) != 1 || len(idsIn(release, "NamespaceIds")) != 1 || len(idsIn(release, "StackIds")) != 1 {
		t.Fatalf("expected the release to carry all three, got %v", release.Body)
	}
}

func TestAModuleAndANamespaceWithTheSameIdAreHeldSeparately(t *testing.T) {
	client, server := newHoldTestClient(t, enabled(60))

	client.HoldModule("x")
	client.HoldNamespace("x")

	if n := len(server.callsTo("/heartbeat")); n != 2 {
		t.Fatalf("expected two heartbeats, got %d", n)
	}
}

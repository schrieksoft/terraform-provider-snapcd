package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

const moduleHoldEndpoint = "/ModuleHold"

// ModuleHoldWarningSummary is the summary of the warning shown when the server changes the hold settings.
const ModuleHoldWarningSummary = "Module hold settings adjusted by the server"

// ModuleHoldSettings configures holding the modules this process writes until it exits.
type ModuleHoldSettings struct {
	Enabled                  bool
	HeartbeatIntervalSeconds *int
	HeartbeatBufferSeconds   *int
}

// holdTarget is one thing a run holds: a module, or every module in a namespace or a stack. Its list is the
// request field that carries its id.
type holdTarget struct {
	list string
	id   string
}

const (
	moduleIdsList    = "ModuleIds"
	namespaceIdsList = "NamespaceIds"
	stackIdsList     = "StackIds"
)

var holdLists = []string{moduleIdsList, namespaceIdsList, stackIdsList}

type heldEntry struct {
	mu   sync.Mutex
	held bool
}

type moduleHolder struct {
	client   *Client
	settings ModuleHoldSettings

	mu         sync.Mutex
	entries    map[holdTarget]*heldEntry
	interval   int
	negotiated bool
	warned     bool
	warning    string

	heartbeatMu sync.Mutex
	loopOnce    sync.Once
	stop        chan struct{}
	stopped     bool
}

var (
	holdersMu sync.Mutex
	holders   []*moduleHolder
)

func newModuleHolder(client *Client, settings ModuleHoldSettings) *moduleHolder {
	h := &moduleHolder{client: client, settings: settings, entries: map[holdTarget]*heldEntry{}, interval: fallbackInterval(settings), stop: make(chan struct{})}
	holdersMu.Lock()
	holders = append(holders, h)
	holdersMu.Unlock()
	return h
}

// fallbackInterval is the heartbeat interval used until the server has said which one applies.
func fallbackInterval(settings ModuleHoldSettings) int {
	if settings.HeartbeatIntervalSeconds != nil && *settings.HeartbeatIntervalSeconds > 0 {
		return *settings.HeartbeatIntervalSeconds
	}
	return 10
}

// timing returns the heartbeat settings the user set, leaving out the ones the server should default.
func (h *moduleHolder) timing() map[string]interface{} {
	timing := map[string]interface{}{}
	if h.settings.HeartbeatIntervalSeconds != nil {
		timing["HeartbeatIntervalSeconds"] = *h.settings.HeartbeatIntervalSeconds
	}
	if h.settings.HeartbeatBufferSeconds != nil {
		timing["HeartbeatBufferSeconds"] = *h.settings.HeartbeatBufferSeconds
	}
	return timing
}

func (h *moduleHolder) active() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.settings.Enabled && !h.stopped
}

func (h *moduleHolder) entry(target holdTarget) *heldEntry {
	h.mu.Lock()
	defer h.mu.Unlock()
	e, ok := h.entries[target]
	if !ok {
		e = &heldEntry{}
		h.entries[target] = e
	}
	return e
}

// heldIds returns the ids held so far, plus adding, per request list; empty when nothing is held.
func (h *moduleHolder) heldIds(adding *holdTarget) map[string][]string {
	h.mu.Lock()
	defer h.mu.Unlock()
	lists := map[string][]string{}
	for target, e := range h.entries {
		if e.held || (adding != nil && target == *adding) {
			lists[target.list] = append(lists[target.list], target.id)
		}
	}
	for _, ids := range lists {
		sort.Strings(ids)
	}
	return lists
}

func withAllLists(body map[string]interface{}, lists map[string][]string) map[string]interface{} {
	for _, list := range holdLists {
		ids := lists[list]
		if ids == nil {
			ids = []string{}
		}
		body[list] = ids
	}
	return body
}

// takeWarning returns the server's warning about the hold settings the first time it is asked after the
// warning arrives, and an empty string otherwise, so it is shown once per process.
func (h *moduleHolder) takeWarning() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	w := h.warning
	h.warning = ""
	return w
}

func (h *moduleHolder) currentInterval() time.Duration {
	h.mu.Lock()
	defer h.mu.Unlock()
	return time.Duration(h.interval) * time.Second
}

// HoldModule holds a module before the first write concerning it. It returns the server's warning about the
// hold settings, once.
func (client *Client) HoldModule(moduleId string) string {
	return client.hold(holdTarget{moduleIdsList, moduleId})
}

// HoldNamespace holds every module in a namespace before the first write concerning the namespace. It returns
// the server's warning about the hold settings, once.
func (client *Client) HoldNamespace(namespaceId string) string {
	return client.hold(holdTarget{namespaceIdsList, namespaceId})
}

// HoldStack holds every module in a stack's namespaces before the first write concerning the stack. It returns
// the server's warning about the hold settings, once.
func (client *Client) HoldStack(stackId string) string {
	return client.hold(holdTarget{stackIdsList, stackId})
}

// hold sends a heartbeat for a target not yet held; parallel writes for the same target wait for it. A failure
// is logged and retried on the next write, and never fails the write.
func (client *Client) hold(target holdTarget) string {
	h := client.holder
	if h == nil || target.id == "" || !h.active() {
		return ""
	}

	e := h.entry(target)
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.held {
		return ""
	}

	if httpError := h.heartbeat(&target, false); httpError != nil {
		log.Printf("[WARN] Could not hold %s %s, jobs may start before this run ends: %v", target.list, target.id, httpError.Error)
		return ""
	}

	h.mu.Lock()
	e.held = true
	h.mu.Unlock()
	h.startLoop()
	return h.takeWarning()
}

// ModuleCreateHoldFields returns the field that asks the server to hold a module it is creating, or nil, and the
// server's warning about the hold settings, once. Before the first hold it asks the server which settings apply.
func (client *Client) ModuleCreateHoldFields() (map[string]interface{}, string) {
	h := client.holder
	if h == nil || !h.active() {
		return nil, ""
	}

	h.mu.Lock()
	negotiated := h.negotiated
	h.mu.Unlock()
	if !negotiated {
		if httpError := h.heartbeat(nil, true); httpError != nil {
			log.Printf("[WARN] Could not reach the module hold endpoint, the new module may not be held: %v", httpError.Error)
		}
	}

	return map[string]interface{}{"Hold": h.timing()}, h.takeWarning()
}

// ModuleCreated records a module the server held on creation, so the heartbeat keeps it held.
func (client *Client) ModuleCreated(moduleId string) {
	h := client.holder
	if h == nil || moduleId == "" || !h.active() {
		return
	}
	e := h.entry(holdTarget{moduleIdsList, moduleId})
	h.mu.Lock()
	e.held = true
	h.mu.Unlock()
	h.startLoop()
}

func (h *moduleHolder) startLoop() {
	h.loopOnce.Do(func() {
		go func() {
			for {
				select {
				case <-h.stop:
					return
				case <-time.After(h.currentInterval()):
					if !h.active() {
						return
					}
					if httpError := h.heartbeat(nil, false); httpError != nil {
						log.Printf("[WARN] Module hold heartbeat failed: %v", httpError.Error)
					}
					if w := h.takeWarning(); w != "" {
						log.Printf("[WARN] %s: %s", ModuleHoldWarningSummary, w)
					}
				}
			}
		}()
	})
}

type moduleHoldHeartbeatResponse struct {
	EffectiveHeartbeatIntervalSeconds int      `json:"effectiveHeartbeatIntervalSeconds"`
	EffectiveExpirySeconds            int      `json:"effectiveExpirySeconds"`
	Warning                           string   `json:"warning"`
	NotHeldModuleIds                  []string `json:"notHeldModuleIds"`
	NotHeldNamespaceIds               []string `json:"notHeldNamespaceIds"`
	NotHeldStackIds                   []string `json:"notHeldStackIds"`
}

func (h *moduleHolder) heartbeat(adding *holdTarget, negotiate bool) *HttpError {
	h.heartbeatMu.Lock()
	defer h.heartbeatMu.Unlock()

	ids := h.heldIds(adding)
	if len(ids) == 0 && !negotiate {
		return nil
	}

	body := withAllLists(h.timing(), ids)

	result, httpError := h.client.Post(moduleHoldEndpoint+"/heartbeat", body)
	if httpError != nil {
		return httpError
	}

	var response moduleHoldHeartbeatResponse
	if raw, err := json.Marshal(result); err == nil {
		_ = json.Unmarshal(raw, &response)
	}
	h.mu.Lock()
	if response.EffectiveHeartbeatIntervalSeconds > 0 {
		h.interval = response.EffectiveHeartbeatIntervalSeconds
	}
	h.negotiated = true
	if response.Warning != "" && !h.warned {
		h.warned = true
		h.warning = response.Warning
	}
	h.mu.Unlock()

	for _, id := range response.NotHeldModuleIds {
		log.Printf("[WARN] Module %s is held by another principal or may not be updated; this run does not hold it", id)
	}
	for _, id := range response.NotHeldNamespaceIds {
		log.Printf("[WARN] Namespace %s is held by another principal or may not be updated; this run does not hold it", id)
	}
	for _, id := range response.NotHeldStackIds {
		log.Printf("[WARN] Stack %s is held by another principal or may not be updated; this run does not hold it", id)
	}
	return nil
}

// ReleaseModuleHolds stops every heartbeat and releases the modules this process holds.
func ReleaseModuleHolds(timeout time.Duration) {
	holdersMu.Lock()
	current := append([]*moduleHolder(nil), holders...)
	holdersMu.Unlock()

	var wg sync.WaitGroup
	for _, h := range current {
		wg.Add(1)
		go func(h *moduleHolder) {
			defer wg.Done()
			h.release(timeout)
		}(h)
	}
	wg.Wait()
}

func (h *moduleHolder) release(timeout time.Duration) {
	h.mu.Lock()
	if h.stopped {
		h.mu.Unlock()
		return
	}
	h.stopped = true
	close(h.stop)
	h.mu.Unlock()

	h.heartbeatMu.Lock()
	defer h.heartbeatMu.Unlock()

	ids := h.heldIds(nil)
	if len(ids) == 0 {
		return
	}

	body, err := json.Marshal(withAllLists(map[string]interface{}{}, ids))
	if err != nil {
		return
	}
	if _, httpError := h.client.doRequest(http.MethodPost, h.client.BuildEndpoint(moduleHoldEndpoint+"/release"), body, timeout); httpError != nil {
		log.Printf("[WARN] Could not release module holds, they lapse at their expiry instead: %v", httpError.Error)
		return
	}
	log.Printf("[DEBUG] Released holds on %s", fmt.Sprint(ids))
}

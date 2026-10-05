package flow

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	codecpb "go.unistack.org/micro-proto/v5/codec"
	"go.unistack.org/micro/v5/store"
)

// WorkflowState is the typed record of a single workflow execution.
//
// All records are keyed by the execution id (EID) — the value returned by
// Workflow.Execute. The workflow definition id is stored inside the record
// (WorkflowID), not used as a key, so a workflow definition can have many
// executions.
type WorkflowState struct {
	// EID is the execution (run) id, the primary key of the record.
	EID string `json:"eid,omitempty"`
	// WorkflowID references the workflow definition (Workflow.ID).
	WorkflowID string `json:"workflow_id,omitempty"`
	// Status is the execution status.
	Status Status `json:"status"`
	// Graph is the step dependency graph: step id -> required step ids.
	Graph map[string][]string `json:"graph,omitempty"`
	// LastStep is the last completed step id.
	LastStep string `json:"last_step,omitempty"`
	// StartedAt is the execution start time (UTC).
	StartedAt time.Time `json:"started_at,omitempty"`
	// FinishedAt is the execution finish time (UTC). Zero until terminal.
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

// StepState is the typed record of a single step of a workflow execution.
type StepState struct {
	// EID is the execution id the step belongs to.
	EID string `json:"eid,omitempty"`
	// StepID is the step id within the workflow.
	StepID string `json:"step_id,omitempty"`
	// Status is the step status.
	Status Status `json:"status"`
	// Req is the last request the step was executed with.
	Req []byte `json:"req,omitempty"`
	// Rsp is the step response (set on success).
	Rsp []byte `json:"rsp,omitempty"`
	// Error is the last attempt error text (set on failure).
	Error string `json:"error,omitempty"`
	// Attempt is the last attempt number (1-based).
	Attempt int `json:"attempt,omitempty"`
	// StartedAt is the last attempt start time (UTC).
	StartedAt time.Time `json:"started_at,omitempty"`
	// FinishedAt is the last attempt finish time (UTC).
	FinishedAt time.Time `json:"finished_at,omitempty"`
	// Checkpoint is the step checkpoint data (see flow.StepCheckpoint).
	// It survives retries and is seeded back into the step context, so a
	// restarted step can continue from the last checkpoint.
	Checkpoint []byte `json:"checkpoint,omitempty"`
}

// WorkflowFilter optionally filters StateStore.WorkflowList results.
// Zero values mean "no filter" for the corresponding field.
type WorkflowFilter struct {
	// WorkflowID filters by workflow definition id.
	WorkflowID string
	// Statuses filters by execution status. Empty means all.
	Statuses []Status
	// Offset is the pagination offset.
	Offset uint
	// Limit is the pagination limit (0 = unlimited).
	Limit uint
}

// WorkflowListEntry is a single result of StateStore.WorkflowList.
type WorkflowListEntry struct {
	// EID is the execution id.
	EID string
	// State is the workflow execution record.
	State *WorkflowState
}

// StateStore is the typed persistence interface for workflow execution state.
//
// Record-based methods (Load/Save/Delete/List) map naturally to table rows
// (e.g. Postgres: upserts by eid), while the narrow hot-path methods
// (WorkflowSetStatus/WorkflowStatus/WorkflowSetLastStep/StepSetCheckpoint)
// allow cheap single-field updates without loading the full record.
//
// Save methods are full upserts: the provided record replaces the stored one.
// Load methods return store.ErrNotFound when the record does not exist.
type StateStore interface {
	// WorkflowLoad loads the workflow execution record by eid.
	WorkflowLoad(ctx context.Context, eid string) (*WorkflowState, error)
	// WorkflowSave upserts the workflow execution record.
	WorkflowSave(ctx context.Context, eid string, st *WorkflowState) error
	// WorkflowDelete deletes the workflow record and all its step records.
	WorkflowDelete(ctx context.Context, eid string) error
	// WorkflowList lists workflow execution records, optionally filtered,
	// newest first.
	WorkflowList(ctx context.Context, f *WorkflowFilter) ([]WorkflowListEntry, error)

	// WorkflowSetStatus updates only the execution status.
	WorkflowSetStatus(ctx context.Context, eid string, s Status) error
	// WorkflowStatus returns only the execution status.
	WorkflowStatus(ctx context.Context, eid string) (Status, error)
	// WorkflowSetLastStep updates only the last completed step id.
	WorkflowSetLastStep(ctx context.Context, eid, sid string) error

	// StepLoad loads the step record of an execution.
	StepLoad(ctx context.Context, eid, sid string) (*StepState, error)
	// StepSave upserts the step record of an execution.
	StepSave(ctx context.Context, eid, sid string, st *StepState) error
	// StepList loads all step records of an execution, keyed by step id.
	StepList(ctx context.Context, eid string) (map[string]*StepState, error)
	// StepSetCheckpoint updates only the checkpoint data of a step.
	StepSetCheckpoint(ctx context.Context, eid, sid string, data []byte) error
}

// WorkflowWatcher is an optional StateStore extension for pushing workflow
// status changes to subscribers instead of relying on polling.
//
// When the state store implements this interface, the engine subscribes to
// every execution running in this process (on registration) and cancels it
// as soon as the store reports a StatusAborted or StatusSuspend status.
type WorkflowWatcher interface {
	// WatchWorkflow subscribes to the status changes of the execution eid.
	// The returned channel first receives the current status (when the
	// record exists) and then every subsequently set status; it is closed
	// when ctx is done.
	WatchWorkflow(ctx context.Context, eid string) (<-chan Status, error)
}

// kvStateStore adapts a plain store.Store to StateStore.
//
// State is stored as flat keys (only key paths, no key parts):
//
//	workflows/<eid>/status         execution status
//	workflows/<eid>/workflow_id    workflow definition id
//	workflows/<eid>/steps          step dependency graph (JSON)
//	workflows/<eid>/last_step      last completed step id
//	workflows/<eid>/started_at     RFC3339 start time
//	workflows/<eid>/finished_at    RFC3339 finish time
//	steps/<eid>/<sid>/status       step status
//	steps/<eid>/<sid>/req          step request
//	steps/<eid>/<sid>/rsp          step response
//	steps/<eid>/<sid>/error        last attempt error
//	steps/<eid>/<sid>/attempt      attempt number
//	steps/<eid>/<sid>/started_at   RFC3339 attempt start time
//	steps/<eid>/<sid>/finished_at  RFC3339 attempt finish time
//	steps/<eid>/<sid>/checkpoint   step checkpoint data
//
// Values are *codecpb.Frame encoded with the store codec, which keeps the
// adapter compatible with data written by the previous key-based engine.
type kvStateStore struct {
	s store.Store

	// In-process status push (WorkflowWatcher). Subscribers belong to this
	// process, so statuses are delivered to in-process subscribers only;
	// for cross-process push use a StateStore implementation backed by a
	// pub/sub store (e.g. Redis).
	subMu sync.Mutex
	subs  map[string]map[chan Status]struct{}
}

// NewKVStateStore creates a StateStore backed by a store.Store.
func NewKVStateStore(s store.Store) StateStore {
	return &kvStateStore{s: s}
}

const (
	kvWorkflowPrefix = "workflows"
	kvStepPrefix     = "steps"
)

func kvWorkflowKey(eid, field string) string {
	return kvWorkflowPrefix + "/" + eid + "/" + field
}

func kvStepKey(eid, sid, field string) string {
	return kvStepPrefix + "/" + eid + "/" + sid + "/" + field
}

// splitKey splits a full key into its namespace (all but the last segment)
// and the last segment. The namespace is passed to the store as a namespace
// option, which keeps the layout compatible with data written through
// store.NewNamespaceStore by the previous engine.
func splitKey(key string) (ns, field string) {
	i := strings.LastIndexByte(key, '/')
	return key[:i], key[i+1:]
}

func (k *kvStateStore) read(ctx context.Context, key string) ([]byte, error) {
	ns, field := splitKey(key)
	f := &codecpb.Frame{}
	if err := k.s.Read(ctx, field, f, store.ReadNamespace(ns)); err != nil {
		return nil, err
	}
	return f.Data, nil
}

func (k *kvStateStore) write(ctx context.Context, key string, data []byte) error {
	ns, field := splitKey(key)
	return k.s.Write(ctx, field, &codecpb.Frame{Data: data}, store.WriteNamespace(ns))
}

func (k *kvStateStore) delete(ctx context.Context, key string) error {
	ns, field := splitKey(key)
	err := k.s.Delete(ctx, field, store.DeleteNamespace(ns))
	if errors.Is(err, store.ErrNotFound) {
		return nil
	}
	return err
}

// readOpt reads a key, reporting ok=false (and no error) when it does not
// exist and propagating other errors.
func (k *kvStateStore) readOpt(ctx context.Context, key string) (data []byte, ok bool, err error) {
	data, err = k.read(ctx, key)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return data, true, nil
}

// listKeys lists the keys below namespace. It works whether the underlying
// store returns namespace-stripped keys or full keys.
func (k *kvStateStore) listKeys(ctx context.Context, ns string) ([]string, error) {
	keys, err := k.s.List(ctx, store.ListNamespace(ns))
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimPrefix(key, ns)
		key = strings.TrimPrefix(key, "/")
		if key != "" {
			out = append(out, key)
		}
	}
	return out, nil
}

func statusFromString(s string) Status {
	s = strings.ToLower(strings.TrimSpace(s))
	// Accept both "StatusRunning" (Status.String) and "Running".
	s = strings.TrimPrefix(s, "status")
	switch s {
	case "running":
		return StatusRunning
	case "success":
		return StatusSuccess
	case "failure", "failed":
		return StatusFailure
	case "aborted":
		return StatusAborted
	case "suspend", "suspended":
		return StatusSuspend
	default:
		return StatusPending
	}
}

func timeFromData(data []byte) time.Time {
	if len(data) == 0 {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return time.Time{}
	}
	return t
}

func (k *kvStateStore) WorkflowLoad(ctx context.Context, eid string) (*WorkflowState, error) {
	data, err := k.read(ctx, kvWorkflowKey(eid, "status"))
	if err != nil {
		return nil, err
	}
	st := &WorkflowState{
		EID:    eid,
		Status: statusFromString(string(data)),
	}

	if data, ok, err := k.readOpt(ctx, kvWorkflowKey(eid, "workflow_id")); err != nil {
		return nil, err
	} else if ok {
		st.WorkflowID = string(data)
	}

	if data, ok, err := k.readOpt(ctx, kvWorkflowKey(eid, "last_step")); err != nil {
		return nil, err
	} else if ok {
		st.LastStep = string(data)
	}

	if data, ok, err := k.readOpt(ctx, kvWorkflowKey(eid, "steps")); err != nil {
		return nil, err
	} else if ok && len(data) > 0 {
		graph := make(map[string][]string)
		if err := json.Unmarshal(data, &graph); err != nil {
			return nil, err
		}
		st.Graph = graph
	}

	if data, ok, err := k.readOpt(ctx, kvWorkflowKey(eid, "started_at")); err != nil {
		return nil, err
	} else if ok {
		st.StartedAt = timeFromData(data)
	}

	if data, ok, err := k.readOpt(ctx, kvWorkflowKey(eid, "finished_at")); err != nil {
		return nil, err
	} else if ok {
		st.FinishedAt = timeFromData(data)
	}

	return st, nil
}

func (k *kvStateStore) WorkflowSave(ctx context.Context, eid string, st *WorkflowState) error {
	if err := k.write(ctx, kvWorkflowKey(eid, "status"), []byte(st.Status.String())); err != nil {
		return err
	}
	if st.WorkflowID != "" {
		if err := k.write(ctx, kvWorkflowKey(eid, "workflow_id"), []byte(st.WorkflowID)); err != nil {
			return err
		}
	}
	if len(st.Graph) > 0 {
		data, err := json.Marshal(st.Graph)
		if err != nil {
			return err
		}
		if err := k.write(ctx, kvWorkflowKey(eid, "steps"), data); err != nil {
			return err
		}
	}
	if st.LastStep != "" {
		if err := k.write(ctx, kvWorkflowKey(eid, "last_step"), []byte(st.LastStep)); err != nil {
			return err
		}
	}
	if !st.StartedAt.IsZero() {
		if err := k.write(ctx, kvWorkflowKey(eid, "started_at"), []byte(st.StartedAt.UTC().Format(time.RFC3339))); err != nil {
			return err
		}
	}
	if !st.FinishedAt.IsZero() {
		if err := k.write(ctx, kvWorkflowKey(eid, "finished_at"), []byte(st.FinishedAt.UTC().Format(time.RFC3339))); err != nil {
			return err
		}
	}
	k.notify(eid, st.Status)
	return nil
}

func (k *kvStateStore) WorkflowDelete(ctx context.Context, eid string) error {
	for _, prefix := range []string{kvWorkflowPrefix + "/" + eid + "/", kvStepPrefix + "/" + eid + "/"} {
		keys, err := k.listKeys(ctx, prefix)
		if err != nil {
			return err
		}
		for _, key := range keys {
			if err := k.delete(ctx, prefix+key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *kvStateStore) WorkflowSetStatus(ctx context.Context, eid string, s Status) error {
	if err := k.write(ctx, kvWorkflowKey(eid, "status"), []byte(s.String())); err != nil {
		return err
	}
	k.notify(eid, s)
	return nil
}

func (k *kvStateStore) WorkflowStatus(ctx context.Context, eid string) (Status, error) {
	data, err := k.read(ctx, kvWorkflowKey(eid, "status"))
	if err != nil {
		return StatusPending, err
	}
	return statusFromString(string(data)), nil
}

func (k *kvStateStore) WorkflowSetLastStep(ctx context.Context, eid, sid string) error {
	return k.write(ctx, kvWorkflowKey(eid, "last_step"), []byte(sid))
}

// WatchWorkflow implements WorkflowWatcher with an in-process pub/sub
// registry: statuses written through this store instance are pushed to its
// subscribers.
func (k *kvStateStore) WatchWorkflow(ctx context.Context, eid string) (<-chan Status, error) {
	ch := make(chan Status, 1)

	// Register before reading the current status so that no write is lost
	// between the two: a concurrent write is either captured by the read or
	// delivered by notify.
	k.subMu.Lock()
	if k.subs == nil {
		k.subs = make(map[string]map[chan Status]struct{})
	}
	if k.subs[eid] == nil {
		k.subs[eid] = make(map[chan Status]struct{})
	}
	k.subs[eid][ch] = struct{}{}
	st, err := k.WorkflowStatus(ctx, eid)
	k.subMu.Unlock()

	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			k.unsub(eid, ch)
			return nil, err
		}
	} else {
		ch <- st
	}

	go func() {
		<-ctx.Done()
		k.unsub(eid, ch)
		close(ch)
	}()

	return ch, nil
}

// unsub removes a subscriber from the in-process pub/sub registry.
func (k *kvStateStore) unsub(eid string, ch chan Status) {
	k.subMu.Lock()
	delete(k.subs[eid], ch)
	k.subMu.Unlock()
}

// notify pushes a status to all in-process subscribers of the execution.
func (k *kvStateStore) notify(eid string, st Status) {
	k.subMu.Lock()
	subs := make([]chan Status, 0, len(k.subs[eid]))
	for ch := range k.subs[eid] {
		subs = append(subs, ch)
	}
	k.subMu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- st:
		default: // subscriber is lagging; the poller is the safety net
		}
	}
}

func (k *kvStateStore) WorkflowList(ctx context.Context, f *WorkflowFilter) ([]WorkflowListEntry, error) {
	keys, err := k.listKeys(ctx, kvWorkflowPrefix)
	if err != nil {
		return nil, err
	}

	// The list contains one entry per stored key; dedupe by execution id.
	eids := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		eid := key
		if i := strings.IndexByte(key, '/'); i >= 0 {
			eid = key[:i]
		}
		if _, ok := seen[eid]; ok {
			continue
		}
		seen[eid] = struct{}{}
		eids = append(eids, eid)
	}

	entries := make([]WorkflowListEntry, 0, len(eids))
	for _, eid := range eids {
		st, err := k.WorkflowLoad(ctx, eid)
		if err != nil {
			// Tolerate partially written records.
			continue
		}
		if f != nil {
			if f.WorkflowID != "" && st.WorkflowID != f.WorkflowID {
				continue
			}
			if len(f.Statuses) > 0 {
				found := false
				for _, s := range f.Statuses {
					if st.Status == s {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
		}
		entries = append(entries, WorkflowListEntry{EID: eid, State: st})
	}

	// Newest first.
	sort.Slice(entries, func(i, j int) bool {
		if !entries[i].State.StartedAt.Equal(entries[j].State.StartedAt) {
			return entries[i].State.StartedAt.After(entries[j].State.StartedAt)
		}
		return entries[i].EID < entries[j].EID
	})

	if f != nil {
		if f.Offset > 0 {
			if int(f.Offset) >= len(entries) {
				return []WorkflowListEntry{}, nil
			}
			entries = entries[f.Offset:]
		}
		if f.Limit > 0 && int(f.Limit) < len(entries) {
			entries = entries[:f.Limit]
		}
	}

	return entries, nil
}

func (k *kvStateStore) StepLoad(ctx context.Context, eid, sid string) (*StepState, error) {
	data, err := k.read(ctx, kvStepKey(eid, sid, "status"))
	if err != nil {
		return nil, err
	}
	st := &StepState{
		EID:    eid,
		StepID: sid,
		Status: statusFromString(string(data)),
	}
	base := kvStepKey(eid, sid, "")

	if data, ok, err := k.readOpt(ctx, base+"req"); err != nil {
		return nil, err
	} else if ok {
		st.Req = data
	}

	if data, ok, err := k.readOpt(ctx, base+"rsp"); err != nil {
		return nil, err
	} else if ok {
		st.Rsp = data
	}

	if data, ok, err := k.readOpt(ctx, base+"error"); err != nil {
		return nil, err
	} else if ok {
		st.Error = string(data)
	}

	if data, ok, err := k.readOpt(ctx, base+"attempt"); err != nil {
		return nil, err
	} else if ok {
		n, perr := strconv.Atoi(string(data))
		if perr != nil {
			return nil, perr
		}
		st.Attempt = n
	}

	if data, ok, err := k.readOpt(ctx, base+"started_at"); err != nil {
		return nil, err
	} else if ok {
		st.StartedAt = timeFromData(data)
	}

	if data, ok, err := k.readOpt(ctx, base+"finished_at"); err != nil {
		return nil, err
	} else if ok {
		st.FinishedAt = timeFromData(data)
	}

	if data, ok, err := k.readOpt(ctx, base+"checkpoint"); err != nil {
		return nil, err
	} else if ok {
		st.Checkpoint = data
	}

	return st, nil
}

func (k *kvStateStore) StepSave(ctx context.Context, eid, sid string, st *StepState) error {
	base := kvStepKey(eid, sid, "")
	if err := k.write(ctx, base+"status", []byte(st.Status.String())); err != nil {
		return err
	}
	if len(st.Req) > 0 {
		if err := k.write(ctx, base+"req", st.Req); err != nil {
			return err
		}
	}
	if len(st.Rsp) > 0 {
		if err := k.write(ctx, base+"rsp", st.Rsp); err != nil {
			return err
		}
	}
	if st.Error != "" {
		if err := k.write(ctx, base+"error", []byte(st.Error)); err != nil {
			return err
		}
	}
	if st.Attempt > 0 {
		if err := k.write(ctx, base+"attempt", []byte(strconv.Itoa(st.Attempt))); err != nil {
			return err
		}
	}
	if !st.StartedAt.IsZero() {
		if err := k.write(ctx, base+"started_at", []byte(st.StartedAt.UTC().Format(time.RFC3339))); err != nil {
			return err
		}
	}
	if !st.FinishedAt.IsZero() {
		if err := k.write(ctx, base+"finished_at", []byte(st.FinishedAt.UTC().Format(time.RFC3339))); err != nil {
			return err
		}
	}
	if st.Checkpoint != nil {
		if err := k.write(ctx, base+"checkpoint", st.Checkpoint); err != nil {
			return err
		}
	} else {
		// Clear any stale checkpoint (e.g. the step finished successfully).
		if err := k.delete(ctx, base+"checkpoint"); err != nil {
			return err
		}
	}
	return nil
}

func (k *kvStateStore) StepList(ctx context.Context, eid string) (map[string]*StepState, error) {
	keys, err := k.listKeys(ctx, kvStepPrefix+"/"+eid+"/")
	if err != nil {
		return nil, err
	}

	sids := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		sid := key
		if i := strings.IndexByte(key, '/'); i >= 0 {
			sid = key[:i]
		}
		if _, ok := seen[sid]; ok {
			continue
		}
		seen[sid] = struct{}{}
		sids = append(sids, sid)
	}

	steps := make(map[string]*StepState, len(sids))
	for _, sid := range sids {
		st, err := k.StepLoad(ctx, eid, sid)
		if err != nil {
			// Tolerate partially written records.
			continue
		}
		steps[sid] = st
	}
	return steps, nil
}

func (k *kvStateStore) StepSetCheckpoint(ctx context.Context, eid, sid string, data []byte) error {
	return k.write(ctx, kvStepKey(eid, sid, "checkpoint"), data)
}

// noopStateStore is a StateStore that accepts and discards all data. It is
// used when the flow is configured without a store, so execution state is
// simply not persisted.
type noopStateStore struct{}

func (noopStateStore) WorkflowSave(ctx context.Context, eid string, st *WorkflowState) error {
	return nil
}

func (noopStateStore) WorkflowLoad(ctx context.Context, eid string) (*WorkflowState, error) {
	return &WorkflowState{EID: eid, Status: StatusRunning}, nil
}

func (noopStateStore) WorkflowDelete(ctx context.Context, eid string) error {
	return nil
}

func (noopStateStore) WorkflowStatus(ctx context.Context, eid string) (Status, error) {
	return StatusRunning, nil
}

func (noopStateStore) WorkflowSetStatus(ctx context.Context, eid string, status Status) error {
	return nil
}

func (noopStateStore) WorkflowSetLastStep(ctx context.Context, eid, stepID string) error {
	return nil
}

func (noopStateStore) WorkflowList(ctx context.Context, filter *WorkflowFilter) ([]WorkflowListEntry, error) {
	return nil, nil
}

func (noopStateStore) StepLoad(ctx context.Context, eid, sid string) (*StepState, error) {
	return &StepState{EID: eid, StepID: sid, Status: StatusRunning}, nil
}

func (noopStateStore) StepSave(ctx context.Context, eid, sid string, st *StepState) error {
	return nil
}

func (noopStateStore) StepList(ctx context.Context, eid string) (map[string]*StepState, error) {
	return map[string]*StepState{}, nil
}

func (noopStateStore) StepSetCheckpoint(ctx context.Context, eid, sid string, data []byte) error {
	return nil
}

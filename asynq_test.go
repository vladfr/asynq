package asynq

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/hibiken/asynq/internal/rdb"
)

// This file defines test helper functions used by
// other test files.

func setup(t *testing.T) *redis.Client {
	t.Helper()
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
	// Start each test with a clean slate.
	if err := r.FlushDB().Err(); err != nil {
		panic(err)
	}
	return r
}

var sortTaskOpt = cmp.Transformer("SortMsg", func(in []*Task) []*Task {
	out := append([]*Task(nil), in...) // Copy input to avoid mutating it
	sort.Slice(out, func(i, j int) bool {
		return out[i].Type < out[j].Type
	})
	return out
})

var sortMsgOpt = cmp.Transformer("SortMsg", func(in []*rdb.TaskMessage) []*rdb.TaskMessage {
	out := append([]*rdb.TaskMessage(nil), in...) // Copy input to avoid mutating it
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID.String() < out[j].ID.String()
	})
	return out
})

func randomTask(taskType, qname string, payload map[string]interface{}) *rdb.TaskMessage {
	return &rdb.TaskMessage{
		ID:      uuid.New(),
		Type:    taskType,
		Queue:   qname,
		Retry:   defaultMaxRetry,
		Payload: make(map[string]interface{}),
	}
}

func mustMarshal(t *testing.T, task *rdb.TaskMessage) string {
	t.Helper()
	data, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func mustUnmarshal(t *testing.T, data string) *rdb.TaskMessage {
	t.Helper()
	var task rdb.TaskMessage
	err := json.Unmarshal([]byte(data), &task)
	if err != nil {
		t.Fatal(err)
	}
	return &task
}

func mustMarshalSlice(t *testing.T, tasks []*rdb.TaskMessage) []string {
	t.Helper()
	var data []string
	for _, task := range tasks {
		data = append(data, mustMarshal(t, task))
	}
	return data
}

func mustUnmarshalSlice(t *testing.T, data []string) []*rdb.TaskMessage {
	t.Helper()
	var tasks []*rdb.TaskMessage
	for _, s := range data {
		tasks = append(tasks, mustUnmarshal(t, s))
	}
	return tasks
}

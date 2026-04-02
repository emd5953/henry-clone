// Package agent implements a multi-agent orchestration engine.
// This mirrors Henry's Golang orchestration that chains AI agents
// to process unstructured documents into presentation-ready decks.
// Each agent is independent — if one fails, the pipeline can reroute,
// retry, or flag the failure without taking down the whole job.
package agent

import (
	"context"
	"sync"
	"time"
)

// Agent is a single unit of work in the pipeline.
// Each agent has a name, a function to execute, and dependencies
// on other agents that must complete first.
type Agent struct {
	Name    string
	Fn      func(ctx context.Context, state *PipelineState) error
	DependsOn []string
	Retries int // max retry attempts, 0 = no retry
	Timeout time.Duration
}

// PipelineState is the shared mutable state that agents read from and write to.
// Agents communicate through this shared state rather than direct message passing.
// Each field is written by one agent and read by downstream agents.
type PipelineState struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewPipelineState() *PipelineState {
	return &PipelineState{data: make(map[string]any)}
}

func (s *PipelineState) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *PipelineState) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// GetTyped is a convenience helper to avoid type assertions everywhere.
func GetTyped[T any](s *PipelineState, key string) (T, bool) {
	v, ok := s.Get(key)
	if !ok {
		var zero T
		return zero, false
	}
	typed, ok := v.(T)
	return typed, ok
}

// AgentResult captures what happened when an agent ran.
type AgentResult struct {
	AgentName string
	Status    AgentStatus
	Duration  time.Duration
	Error     error
	Attempts  int
}

type AgentStatus string

const (
	StatusSuccess AgentStatus = "success"
	StatusFailed  AgentStatus = "failed"
	StatusSkipped AgentStatus = "skipped"
)

// Pipeline orchestrates a DAG of agents with dependency resolution,
// parallel execution where possible, and per-agent retry/timeout.
type Pipeline struct {
	agents  []Agent
	results map[string]*AgentResult
	mu      sync.Mutex
}

func NewPipeline(agents ...Agent) *Pipeline {
	return &Pipeline{
		agents:  agents,
		results: make(map[string]*AgentResult),
	}
}

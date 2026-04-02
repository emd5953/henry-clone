package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Run executes the pipeline: resolves dependencies, runs agents in
// parallel where possible, handles retries and timeouts.
// Returns all results — callers decide what to do with failures.
func (p *Pipeline) Run(ctx context.Context, state *PipelineState) []AgentResult {
	// Build lookup and dependency graph
	agentMap := make(map[string]*Agent)
	for i := range p.agents {
		agentMap[p.agents[i].Name] = &p.agents[i]
	}

	// Track completion
	done := make(map[string]chan struct{})
	for _, a := range p.agents {
		done[a.Name] = make(chan struct{})
	}

	var wg sync.WaitGroup

	for i := range p.agents {
		wg.Add(1)
		go func(a Agent) {
			defer wg.Done()
			defer close(done[a.Name])

			// Wait for dependencies
			for _, dep := range a.DependsOn {
				ch, ok := done[dep]
				if !ok {
					p.recordResult(a.Name, StatusFailed, 0, 0, fmt.Errorf("unknown dependency: %s", dep))
					return
				}
				select {
				case <-ch:
					// Check if dependency succeeded
					p.mu.Lock()
					depResult := p.results[dep]
					p.mu.Unlock()
					if depResult != nil && depResult.Status == StatusFailed {
						p.recordResult(a.Name, StatusSkipped, 0, 0, fmt.Errorf("dependency %s failed", dep))
						return
					}
				case <-ctx.Done():
					p.recordResult(a.Name, StatusFailed, 0, 0, ctx.Err())
					return
				}
			}

			// Execute with retries
			maxAttempts := a.Retries + 1
			var lastErr error
			start := time.Now()

			for attempt := 1; attempt <= maxAttempts; attempt++ {
				var execCtx context.Context
				var cancel context.CancelFunc

				if a.Timeout > 0 {
					execCtx, cancel = context.WithTimeout(ctx, a.Timeout)
				} else {
					execCtx, cancel = context.WithCancel(ctx)
				}

				lastErr = a.Fn(execCtx, state)
				cancel()

				if lastErr == nil {
					p.recordResult(a.Name, StatusSuccess, time.Since(start), attempt, nil)
					return
				}

				// Don't retry on context cancellation
				if ctx.Err() != nil {
					break
				}
			}

			p.recordResult(a.Name, StatusFailed, time.Since(start), maxAttempts, lastErr)
		}(p.agents[i])
	}

	wg.Wait()

	// Collect results in order
	results := make([]AgentResult, 0, len(p.agents))
	p.mu.Lock()
	for _, a := range p.agents {
		if r, ok := p.results[a.Name]; ok {
			results = append(results, *r)
		}
	}
	p.mu.Unlock()

	return results
}

func (p *Pipeline) recordResult(name string, status AgentStatus, dur time.Duration, attempts int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.results[name] = &AgentResult{
		AgentName: name,
		Status:    status,
		Duration:  dur,
		Error:     err,
		Attempts:  attempts,
	}
}

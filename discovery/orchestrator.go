package discovery

import (
	"context"
	"sync"
)

// Orchestrator runs multiple discovery providers.
type Orchestrator struct {
	providers []Provider
	mgr       *Manager
}

func NewOrchestrator(mgr *Manager, providers ...Provider) *Orchestrator {
	return &Orchestrator{providers: providers, mgr: mgr}
}

func (o *Orchestrator) Discover(ctx context.Context, opts Options) ([]Result, []error) {
	if len(o.providers) == 0 {
		return []Result{}, nil
	}
	type res struct {
		items []Result
		err   error
	}
	ch := make(chan res, len(o.providers))
	var wg sync.WaitGroup
	for _, p := range o.providers {
		p := p
		wg.Add(1)
		go func() {
			defer wg.Done()
			items, err := p.Discover(ctx, opts)
			ch <- res{items: items, err: err}
		}()
	}
	wg.Wait()
	close(ch)
	all := []Result{}
	var errs []error
	seen := map[string]bool{}
	for r := range ch {
		if r.err != nil {
			errs = append(errs, r.err)
			continue
		}
		for _, it := range r.items {
			if !Validate(it) {
				continue
			}
			key := it.ID
			if key == "" {
				key = it.Address
			}
			if seen[key] {
				continue
			}
			seen[key] = true
			all = append(all, it)
		}
	}
	if o.mgr != nil {
		all = o.mgr.Deduplicate(all)
	}
	if opts.MaxResults > 0 && len(all) > opts.MaxResults {
		all = all[:opts.MaxResults]
	}
	return all, errs
}

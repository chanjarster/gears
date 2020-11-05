package confstore

import (
	"sync"
	"time"
)

var (
	NoopLoadPolicy   LoadPolicy = &noopLoadPolicy{}     // Never load config key-values from Persister to Interface
	SimpleLoadPolicy            = NewSimpleLoadPolicy() // Load config key-values from Persister immediately
)

// Never load config key-values from Persister to Interface
type noopLoadPolicy struct{}

func (n *noopLoadPolicy) DoLoad(s Interface, p Persister) error {
	return nil
}

// Create a default load policy that load config key-values from Persister
// only when last load timestamp is `minInterval` before current time.
//
// This implementation protect underlying Persister from being queried too often
func NewLoadPolicy(minInterval time.Duration) LoadPolicy {
	return &defaultLoadPolicy{
		loadLock:    &sync.Mutex{},
		minInterval: minInterval,
	}
}

type defaultLoadPolicy struct {
	loadLock    *sync.Mutex
	lastLoadTs  time.Time
	minInterval time.Duration
}

func (t *defaultLoadPolicy) DoLoad(s Interface, p Persister) error {
	t.loadLock.Lock()
	defer t.loadLock.Unlock()

	if t.minInterval <= 0 {
		return p.Load(s)
	}

	now := time.Now()
	if now.Add(-t.minInterval).Before(t.lastLoadTs) {
		return nil
	}

	err := p.Load(s)
	if err != nil {
		return err
	}
	t.lastLoadTs = time.Now()
	return nil
}

func NewSimpleLoadPolicy() LoadPolicy {
	return &simpleLoadPolicy{
	}
}

type simpleLoadPolicy struct {
}

func (d *simpleLoadPolicy) DoLoad(s Interface, p Persister) error {
	return p.Load(s)
}

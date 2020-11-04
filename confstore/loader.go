package configstore

import (
	"sync"
	"time"
)

var (
	NoopLoadPolicy LoadPolicy = &noopLoadPolicy{} // Never load config key-values from Persister to Interface
)

// Never load config key-values from Persister to Interface
type noopLoadPolicy struct{}

func (n *noopLoadPolicy) DoLoad(s Interface, p Persister) {}

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

func (t defaultLoadPolicy) DoLoad(s Interface, p Persister) {
	t.loadLock.Lock()
	defer t.loadLock.Unlock()

	if t.minInterval <= 0 {
		p.Load(s)
		return
	}

	now := time.Now()
	if now.Add(-t.minInterval).Before(t.lastLoadTs) {
		return
	}

	p.Load(s)
	t.lastLoadTs = time.Now()
}

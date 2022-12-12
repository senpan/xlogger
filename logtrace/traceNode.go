package logtrace

import "sync"

type TraceNode struct {
	metadata map[string]string
	lock     *sync.RWMutex
}

func NewTraceNode() *TraceNode {
	t := new(TraceNode)
	t.metadata = make(map[string]string, 6)
	t.lock = new(sync.RWMutex)
	return t
}

func (tn *TraceNode) Get(key string) string {
	tn.lock.RLock()
	defer tn.lock.RUnlock()
	return tn.metadata[key]
}

func (tn *TraceNode) Set(key, val string) {
	tn.lock.Lock()
	defer tn.lock.Unlock()
	tn.metadata[key] = val
}

func (tn *TraceNode) ForkMap() map[string]string {
	ret := make(map[string]string, 5)
	tn.lock.RLock()
	defer tn.lock.RUnlock()
	for k, v := range tn.metadata {
		ret[k] = v
	}
	return ret
}

package session

import (
	"sync"
	"time"
	"container/list"
	"fmt"
)

type MemoryProvider struct {
	maxLifeTime int64
	lock        sync.RWMutex
	bytes       int64
	data        map[string]*list.Element
	list        *list.List
	options     MemoryOptions
}

type MemoryOptions struct {
	BytesLimit int64
}

type item struct {
	sid        string
	lastAccess time.Time
	data       []byte
}

func (p *MemoryProvider) Init(maxLifeTime int64, options interface{}) error {
	p.maxLifeTime = maxLifeTime
	p.options = options.(MemoryOptions)
	return nil
}

func (p *MemoryProvider) Exist(sid string) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	_, ok := p.data[sid]
	return ok
}

func (p *MemoryProvider) Read(sid string) (*Session, error) {
	p.lock.RLock()
	elem, ok := p.data[sid]
	p.lock.RUnlock()

	if !ok {
		return NewSession(p, sid, nil)
	}

	p.touch(sid)

	data, err := DecodeGob(elem.Value.(*item).data)
	if err != nil {
		return nil, err
	}

	return NewSession(p, sid, data)
}

func (p *MemoryProvider) touch(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if elem, ok := p.data[sid]; ok {
		elem.Value.(*item).lastAccess = time.Now()
		p.list.MoveToFront(elem)
	}

	return nil
}

func (p *MemoryProvider) Write(sid string, data map[interface{}]interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	encoded, err := EncodeGob(data)
	if err != nil {
		return err
	}

	bytes := p.bytes + int64(len(encoded))
	if bytes > p.options.BytesLimit {
		return fmt.Errorf("session size reached to size limit %d", p.options.BytesLimit)
	}

	p.bytes = bytes
	p.data[sid] = p.list.PushBack(&item{
		sid: sid,
		lastAccess: time.Now(),
		data: encoded,
	})

	return nil
}

func (p *MemoryProvider) Destroy(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if elem, ok := p.data[sid]; ok {
		p.bytes -= int64(len(elem.Value.(*item).data))
		p.list.Remove(elem)
		delete(p.data, sid)
	}

	return nil
}

func (p *MemoryProvider) GC() {
	p.lock.RLock()

	for {
		elem := p.list.Back()
		if elem == nil {
			break
		}

		if elem.Value.(*item).lastAccess.Unix() + p.maxLifeTime < time.Now().Unix() {
			p.lock.RUnlock()
			p.lock.Lock()
			p.bytes -= int64(len(elem.Value.(*item).data))
			p.list.Remove(elem)
			delete(p.data, elem.Value.(*item).sid)
			p.lock.Unlock()
			p.lock.RLock()
		} else {
			break
		}
	}

	p.lock.RUnlock()
}

func init() {
	Register("memory", &MemoryProvider{
		list: list.New(),
		data: make(map[string]*list.Element),
	})
}
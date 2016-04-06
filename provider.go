package session

import (
	"fmt"
)

var adapters = make(map[string]Provider)

type Provider interface {
	Init(maxLifeTime int64, options interface{}) error

	Exist(sid string) bool
	Read(sid string) (*Session, error)
	Write(sid string, data map[interface{}]interface{}) error

	Destroy(sid string) error

	// GC calls GC to clean expired sessions.
	GC()
}

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session.Register(): cannot register nil provider")
	}
	if _, dup := adapters[name]; dup {
		panic(fmt.Errorf("session.Register(): provider '%s' already exists", name))
	}
	adapters[name] = provider
}

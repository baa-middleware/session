package session

import (
	"fmt"
)

var adapters = make(map[string]Provider)

// Provider sessio provider interface
type Provider interface {
	// Init initialize provider
	Init(maxLifeTime int64, options interface{}) error
	// Exist check session id is exist
	Exist(sid string) bool
	// Read read session data from provider
	Read(sid string) (*Session, error)
	// Write write session data to provider
	Write(sid string, data map[interface{}]interface{}) error
	// Destroy destroy session id from provider
	Destroy(sid string) error
	// GC calls GC to clean expired sessions
	GC()
}

// Register register a provider
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session.Register(): cannot register nil provider")
	}
	if _, dup := adapters[name]; dup {
		panic(fmt.Errorf("session.Register(): provider '%s' already exists", name))
	}
	adapters[name] = provider
}

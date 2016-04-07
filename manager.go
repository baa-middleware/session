package session

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/baa.v1"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type Manager struct {
	options  Options
	provider Provider
}

func NewManager(options Options) (*Manager, error) {
	if len(options.Name) == 0 {
		panic("session.NewManager(): name cannot be empty")
	}

	if options.IDLength == 0 {
		options.IDLength = 16
	}

	if options.GCInterval == 0 {
		if options.MaxLifeTime > 0 {
			options.GCInterval = options.MaxLifeTime
		} else {
			options.GCInterval = 1440
		}
	}

	if options.Provider == nil {
		panic("session.NewManager(): provider option cannot be nil")
	}

	provider, ok := adapters[options.Provider.Adapter]
	if !ok {
		panic(fmt.Errorf("session.NewManager(): unknown provider adapter '%s'", options.Provider.Adapter))
	}

	if err := provider.Init(options.MaxLifeTime, options.Provider.Config); err != nil {
		panic(fmt.Errorf("session.NewManager(): init provider '%s' error with config", options.Provider.Adapter))
	}

	if options.Cookie == nil {
		options.Cookie = new(CookieOptions)
	}

	manager := &Manager{
		options:  options,
		provider: provider,
	}

	return manager, nil
}

func (m *Manager) Start(c *baa.Context) (*Session, error) {
	var session *Session

	sid := c.GetCookie(m.options.Name)
	if len(sid) > 0 && m.provider.Exist(sid) {
		return m.provider.Read(sid)
	}

	sid = m.sessionId()
	session, err := m.provider.Read(sid)
	if err != nil {
		return nil, err
	}

	c.SetCookie(
		m.options.Name,
		sid,
		m.options.Cookie.LifeTime,
		m.options.Cookie.Path,
		m.options.Cookie.Domain,
		m.options.Cookie.Secure,
		m.options.Cookie.HttpOnly,
	)

	return session, nil
}

func (m *Manager) GC() {
	m.provider.GC()
}

func (m *Manager) startGC() {
	m.GC()
	time.AfterFunc(time.Duration(m.options.GCInterval)*time.Second, func() {
		m.startGC()
	})
}

func (m *Manager) sessionId() string {
	return randString(m.options.IDLength)
}

func randString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

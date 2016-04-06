package session

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gopkg.in/baa.v1"
	"sync"
)

type Session struct {
	provider Provider
	sid      string
	lock     sync.RWMutex
	data     map[interface{}]interface{}
}

func (s *Session) ID() string {
	return s.sid
}

func (s *Session) Get(key interface{}) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data[key]
}

func (s *Session) Set(key interface{}, val interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = val
	return nil
}

func (s *Session) Delete(key interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.data, key)
	return nil
}

// Destroy remove all data stored in a session
func (s *Session) Destroy() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data = make(map[interface{}]interface{})
	return nil
}

func (s *Session) Close() error {
	return s.provider.Write(s.sid, s.data)
}

func NewSession(provider Provider, sid string, data map[interface{}]interface{}) (*Session, error) {
	if provider == nil {
		return nil, fmt.Errorf("session.New(): provider cannot be nil")
	}

	if len(sid) == 0 {
		return nil, fmt.Errorf("session.New(): invalid session id")
	}

	return &Session{
		provider: provider,
		sid:      sid,
		data:     data,
	}, nil
}

func Middleware(option Options) baa.Handler {
	manager, err := NewManager(option)
	if err != nil {
		panic(err)
	}
	go manager.startGC()

	return func(c *baa.Context) {
		// Start session
		session, err := manager.Start(c)
		if err != nil {
			panic("session.Start(): " + err.Error())
		}

		// SetDI allows reference session instance duration a request
		c.Baa().SetDI("session", session)

		c.Next()

		// Close session
		if err := session.Close(); err != nil {
			panic("session.Close(): " + err.Error())
		}
	}
}

func EncodeGob(object map[interface{}]interface{}) ([]byte, error) {
	// Register key and value type in order to encode without error
	for key, value := range object {
		gob.Register(key)
		gob.Register(value)
	}

	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(object)

	return buf.Bytes(), err
}

func DecodeGob(encoded []byte) (out map[interface{}]interface{}, err error) {
	buf := bytes.NewBuffer(encoded)
	err = gob.NewDecoder(buf).Decode(&out)
	return out, err
}

package session

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
	"sync"

	"gopkg.in/baa.v1"
)

// Session session bag
type Session struct {
	provider Provider
	sid      string
	lock     sync.RWMutex
	data     map[interface{}]interface{}
}

// ID returns session id
func (s *Session) ID() string {
	return s.sid
}

// Get get a value by the key from session
func (s *Session) Get(key interface{}) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data[key]
}

// Set set a value by the key to session
func (s *Session) Set(key interface{}, val interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = val
	return nil
}

// Delete delete a key from session
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

// Close save data in session
func (s *Session) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.provider.Write(s.sid, s.data)
}

// NewSession create a session
func NewSession(provider Provider, sid string, data map[interface{}]interface{}) (*Session, error) {
	if provider == nil {
		return nil, fmt.Errorf("session.New(): provider cannot be nil")
	}

	if len(sid) == 0 {
		return nil, fmt.Errorf("session.New(): invalid session id")
	}

	if data == nil {
		data = make(map[interface{}]interface{})
	}

	return &Session{
		provider: provider,
		sid:      sid,
		data:     data,
	}, nil
}

// Middleware returns a middleware for baa
func Middleware(option Options) baa.HandlerFunc {
	manager, err := NewManager(option)
	if err != nil {
		panic(err)
	}
	go manager.startGC()

	var reStatic = regexp.MustCompile(`\.(jpeg|jpg|png|gif|ico|js|css|txt|zip)$`)

	return func(c *baa.Context) {
		// skip static request
		if reStatic.MatchString(c.Req.URL.Path) {
			c.Next()
			return
		}

		// Start session
		session, err := manager.Start(c)
		if err != nil {
			panic("session.Start(): " + err.Error())
		}

		// allows reference session instance in context
		c.Set("session", session)

		c.Next()

		// Close session
		if err := session.Close(); err != nil {
			panic("session.Close(): " + err.Error())
		}
	}
}

// EncodeGob encode data by gob
func EncodeGob(object map[interface{}]interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(object)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

// DecodeGob decode data by gob
func DecodeGob(encoded []byte) (out map[interface{}]interface{}, err error) {
	buf := bytes.NewBuffer(encoded)
	err = gob.NewDecoder(buf).Decode(&out)
	return out, err
}

package redis

import (
	"time"

	"github.com/baa-middleware/session"
	"gopkg.in/redis.v3"
)

type RedisProvider struct {
	client      *redis.Client
	maxLifeTime int64
	options     Options
}

type Options struct {
	redis.Options
	Prefix string
}

func (p *RedisProvider) Init(maxLifeTime int64, options interface{}) error {
	p.options = options.(Options)

	client := redis.NewClient(&redis.Options{
		Network: p.options.Network,
		Addr: p.options.Addr,
		Dialer: p.options.Dialer,
		Password: p.options.Password,
		DB: p.options.DB,
		DialTimeout: p.options.DialTimeout,
		ReadTimeout: p.options.ReadTimeout,
		WriteTimeout: p.options.WriteTimeout,
		PoolSize: p.options.PoolSize,
		IdleTimeout: p.options.IdleTimeout,
	})
	err := client.Ping().Err()
	if err != nil {
		return err
	}

	p.client = client
	p.maxLifeTime = maxLifeTime

	return nil
}

func (p *RedisProvider) Exist(sid string) bool {
	has, err := p.client.Exists(p.options.Prefix+sid).Result()
	if err != nil {
		return false
	}
	return has
}

func (p *RedisProvider) Read(sid string) (*session.Session, error) {
	psid := p.options.Prefix + sid
	if !p.Exist(sid) {
		if err := p.client.Set(psid, "", time.Second * time.Duration(p.maxLifeTime)).Err(); err != nil {
			return nil, err
		}
	}

	var data map[interface{}]interface{}
	raw, err := p.client.Get(psid).Result()
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		data = make(map[interface{}]interface{})
	} else {
		data, err = session.DecodeGob([]byte(raw))
		if err != nil {
			return nil, err
		}
	}

	return session.NewSession(p, sid, data)
}

func (p *RedisProvider) Write(sid string, data map[interface{}]interface{}) error {
	encoded, err := session.EncodeGob(data)
	if err != nil {
		return err
	}
	return p.client.Set(
		p.options.Prefix + sid,
		string(encoded),
		time.Second * time.Duration(p.maxLifeTime),
	).Err()
}

func (p *RedisProvider) Destroy(sid string) error {
	return p.client.Del(p.options.Prefix + sid).Err()
}

func (_ *RedisProvider) GC() {}

func init() {
	session.Register("redis", &RedisProvider{})
}
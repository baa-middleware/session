# session [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/baa-middleware/session) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/baa-middleware/session/master/LICENSE) [![Build Status](http://img.shields.io/travis/go-baa/cache.svg?style=flat-square)](https://travis-ci.org/baa-middleware/session) [![Coverage Status](http://img.shields.io/coveralls/baa-middleware/session.svg?style=flat-square)](https://coveralls.io/r/baa-middleware/session)

baa middleware for provides the session management.

## Session Provider Adapters

- redis
- memory

## Getting Started

```
package main

import (
	"gopkg.in/baa.v1"
	"github.com/baa-middleware/session"
	"github.com/baa-middleware/session/redis"
)

func main() {
	// new app
	app := baa.New()

	// use session middleware
	redisOptions := redis.Options{}
	redisOptions.Addr = "127.0.0.1:6379"
	redisOptions.Prefix = "Prefix:"

	app.Use(session.Middleware(session.Options{
		Name: "GSESSION",
		Provider: &session.ProviderOptions{
			Adapter: "redis",
			Config: redisOptions,
		},
	}))

	// router
	app.Get("/", func(c *baa.Context) {
		// get the session handler
		session := c.Baa().GetDI("session").(*session.Session)

		// get
		session.Get("key")
		
		// set
		session.Set("key", "value")
		
		// delete
		session.Delete("key")

		// destroy
		session.Destroy()
	})

	// run app
	app.Run(":1323")
}
```

## Configuration

Name        string
	IDLength    int
	Provider    *ProviderOptions
	Cookie      *CookieOptions
	GCInterval  int64
	MaxLifeTime int64

### Global

#### Name `string`

session name

#### IDLength `int`

session id length, default is `16`

#### Provider `*ProviderOptions`

provider options

##### Adapter `string`

provider adapter name, currently support `redis` and `memory`

##### Config `interface{}`
 
provider adapter config, each adapter has its own config

#### Cookie `*CookieOptions`

session cookie options

##### Domain `string`

cookie domain, default is `''`

##### Path `string`

cookie path, default is `/`

##### Secure `bool`

##### LifeTime `int64`

cookie life time, default is `0`, known as session cookie

##### HttpOnly `bool`

#### GCInterval `int64`

garbage collection run interval, used for `memory` adapter only

#### MaxLifeTime `int64`

After this number of seconds, stored data will be seen as 'garbage' and cleaned up by the garbage collection process


## Credits

Get inspirations from [macaron session](https://github.com/go-macaron/session)
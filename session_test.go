package session

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/baa.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app = baa.New()

func TestSessionMiddleware(t *testing.T) {
	Convey("middleware memory", t, func() {
		app.Use(Middleware(Options{
			Name: "GSESSION",
			Provider: &ProviderOptions{
				Adapter: "memory",
				Config: MemoryOptions{
					BytesLimit: 10000,
				},
			},
		}))

		app.Get("/", func(c *baa.Context) {
			// get the session handler
			session := c.Baa().GetDI("session").(*Session)

			err := session.Set("count", 1)
			So(err, ShouldBeNil)

			count := session.Get("count").(int)
			So(count, ShouldEqual, 1)

			err = session.Delete("count")
			So(err, ShouldBeNil)

			So(session.Get("count"), ShouldBeNil)
		})

		w := request("GET", "/")
		So(w.Code, ShouldEqual, http.StatusOK)
	})

	Convey("middleware redis", t, func() {

		redisOptions := RedisOptions{}
		redisOptions.Addr = "127.0.0.1:6379"
		redisOptions.Prefix = "test:"

		app.Use(Middleware(Options{
			Name: "GSESSION",
			Provider: &ProviderOptions{
				Adapter: "redis",
				Config: redisOptions,
			},
			MaxLifeTime: 1440,
		}))

		app.Get("/", func(c *baa.Context) {
			// get the session handler
			session := c.Baa().GetDI("session").(*Session)

			So(session.ID(), ShouldNotBeNil)

			err := session.Set("count", 1)
			So(err, ShouldBeNil)

			count := session.Get("count").(int)
			So(count, ShouldEqual, 1)

			err = session.Delete("count")
			So(err, ShouldBeNil)

			So(session.Get("count"), ShouldBeNil)

			So(session.Destroy(), ShouldBeNil)
		})

		w := request("GET", "/")
		So(w.Code, ShouldEqual, http.StatusOK)
	})

	Convey("new session", t, func() {

		redisOptions := RedisOptions{}
		redisOptions.Addr = "127.0.0.1:6379"
		redisOptions.Prefix = "test:"

		provider := &RedisProvider{}
		provider.Init(10, redisOptions)

		sid := "test"

		data := make(map[interface{}]interface{})

		Convey("provider is nil", func() {
			NewSession(nil, sid, data)
		})

		Convey("sid is empty", func() {
			NewSession(provider, "", data)
		})
	})
}

func request(method, uri string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, uri, nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	return w
}

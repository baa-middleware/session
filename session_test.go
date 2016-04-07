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
	Convey("middleware", t, func() {
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
}

func request(method, uri string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, uri, nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	return w
}

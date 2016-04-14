package session

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSessionManager(t *testing.T) {
	Convey("new manager", t, func() {

		Convey("panic by name", func() {
			So(func() {
				NewManager(Options{
					Name: "",
				})
			}, ShouldPanic)
		})

		Convey("panic by provider", func() {
			So(func() {
				NewManager(Options{
					Name: "SESSION",
				})
			}, ShouldPanic)
		})

		Convey("panic by provider adapter name", func() {
			So(func() {
				NewManager(Options{
					Name:     "SESSION",
					Provider: &ProviderOptions{},
				})
			}, ShouldPanic)
		})

		Convey("panic by provider adapter config", func() {
			So(func() {
				redisOptions := RedisOptions{}
				redisOptions.Addr = "127.0.0.2:1234"

				NewManager(Options{
					Name: "SESSION",
					Provider: &ProviderOptions{
						Adapter: "redis",
						Config:  redisOptions,
					},
				})
			}, ShouldPanic)
		})
	})

	Convey("start", t, func() {
		// TODO
	})

	Convey("gc", t, func() {
		// TODO
	})
}

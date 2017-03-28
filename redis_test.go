package session

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSessionRedisProvider(t *testing.T) {

	redisOptions := RedisOptions{}
	redisOptions.Addr = "127.0.0.1:6379"
	redisOptions.Prefix = "test:"

	provider := &RedisProvider{}
	provider.Init(1440, redisOptions)

	Convey("exist", t, func() {
		So(provider.Exist("foo"), ShouldBeFalse)
	})

	Convey("read exists", t, func() {
		data := make(map[interface{}]interface{})
		data["a"] = "b"
		provider.Write("fooOOOOO", data)

		session, err := provider.Read("fooOOOOO")
		So(err, ShouldBeNil)
		So(session, ShouldNotBeNil)

		provider.Destroy("fooOOOOO")
	})

	Convey("read not exists", t, func() {
		session, err := provider.Read("fooOOOOO")
		So(err, ShouldBeNil)
		So(session, ShouldNotBeNil)
	})

	Convey("write", t, func() {
		Convey("ok", func() {
			data := make(map[interface{}]interface{})
			data["a"] = "b"

			So(provider.Write("foo", data), ShouldBeNil)
		})
	})

	Convey("destroy", t, func() {
		So(provider.Destroy("foo"), ShouldBeNil)
	})
}

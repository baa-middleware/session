package session

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"container/list"
	"time"
)

func TestSessionMemoryProvider(t *testing.T) {

	provider := &MemoryProvider{
		list: list.New(),
		data: make(map[string]*list.Element),
	}
	provider.Init(3, MemoryOptions{
		BytesLimit: 50,
	})

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

		Convey("size limit", func() {
			data := make(map[interface{}]interface{})
			data["a"] = string(randString(1000))
			So(provider.Write("foo", data), ShouldNotBeNil)
		})
	})

	Convey("destroy", t, func() {
		So(provider.Destroy("foo"), ShouldBeNil)
	})

	Convey("gc", t, func() {
		data := make(map[interface{}]interface{})
		data["a"] = "b"
		provider.Write("first", data)

		time.Sleep(time.Second * 4)

		data["a"] = "b"
		provider.Write("second", data)

		provider.GC()
	})
}
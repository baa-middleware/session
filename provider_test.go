package session

import (
	"container/list"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSessionProvider(t *testing.T) {
	Convey("register", t, func() {

		Convey("panic by provider", func() {
			So(func() {
				Register("", nil)
			}, ShouldPanic)
		})

		Convey("panic by name already exists", func() {
			So(func() {
				Register("memory", &MemoryProvider{
					list: list.New(),
					data: make(map[string]*list.Element),
				})

				Register("memory", &MemoryProvider{
					list: list.New(),
					data: make(map[string]*list.Element),
				})
			}, ShouldPanic)
		})
	})
}

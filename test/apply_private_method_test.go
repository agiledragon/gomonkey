package test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2/test/fake"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/agiledragon/gomonkey/v2"
)

func TestApplyPrivateMethod(t *testing.T) {
	Convey("TestApplyPrivateMethod", t, func() {
		Convey("patch private pointer method in the different package", func() {
			f := new(fake.PrivateMethodStruct)
			var s *fake.PrivateMethodStruct
			patches := ApplyPrivateMethod(s, "ok", func(_ *fake.PrivateMethodStruct) bool {
				return false
			})
			defer patches.Reset()
			result := f.Happy()
			So(result, ShouldEqual, "unhappy")
		})

		Convey("patch private value method in the different package", func() {
			s := fake.PrivateMethodStruct{}
			patches := ApplyPrivateMethod(s, "haveEaten", func(_ fake.PrivateMethodStruct) bool {
				return true
			})
			defer patches.Reset()
			result := s.AreYouHungry()
			So(result, ShouldEqual, "I am full")
		})

		Convey("repeat patch same method", func() {
			var s *fake.PrivateMethodStruct
			patches := ApplyPrivateMethod(s, "ok", func(_ *fake.PrivateMethodStruct) bool {
				return false
			})
			result := s.Happy()
			So(result, ShouldEqual, "unhappy")

			patches.ApplyPrivateMethod(s, "ok", func(_ *fake.PrivateMethodStruct) bool {
				return true
			})
			result = s.Happy()
			So(result, ShouldEqual, "happy")

			patches.Reset()
			result = s.Happy()
			So(result, ShouldEqual, "unhappy")
		})
	})

}

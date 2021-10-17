package test

import (
    "github.com/agiledragon/gomonkey/v2/test/fake"
    "reflect"
    "testing"

    . "github.com/agiledragon/gomonkey/v2"
    . "github.com/smartystreets/goconvey/convey"
)

func TestApplyPrivate(t *testing.T) {
	Convey("TestApplyPrivate", t, func() {
        Convey("patch private method in the different package", func() {
            f := new(fake.PrivateMethodStruct)
            var s *fake.PrivateMethodStruct
            patches := ApplyPrivateMethod(reflect.TypeOf(s), "ok", func(_ *fake.PrivateMethodStruct) bool {
                return false
            })
            defer patches.Reset()
            result := f.Happy()
            So(result, ShouldEqual, "unhappy")
        })
    })
}

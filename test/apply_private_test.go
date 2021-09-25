package test

import (
    "reflect"
    "testing"

    . "github.com/agiledragon/gomonkey/v2"
    . "github.com/smartystreets/goconvey/convey"
)

func fakeFunc() bool {
    return true
}

type fakeStruct struct {

}

func (*fakeStruct) fakeMethod() bool {
    return true
}

func TestApplyPrivate(t *testing.T) {
    Convey("TestApplyPrivate", t, func() {
        Convey("func", func() {
            patches := ApplyFunc(fakeFunc, func() bool {
                return false
            })
            defer patches.Reset()
            flag := fakeFunc()
            So(flag, ShouldEqual, false)
        })

        Convey("method", func() {
            f := new(fakeStruct)
            var s *fakeStruct
            patches := ApplyMethod(reflect.TypeOf(s), "fakeMethod", func(_ *fakeStruct) bool {
                return false
            })
            defer patches.Reset()
            flag := f.fakeMethod()
            So(flag, ShouldEqual, false)
        })


    })
}

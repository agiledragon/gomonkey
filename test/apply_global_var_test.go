package test

import (
    . "github.com/henrylee2cn/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
)

var num = 10

func TestApplyGlobalVar(t *testing.T) {
    Convey("TestApplyGlobalVar", t, func() {

        Convey("change", func() {
            patches := ApplyGlobalVar(&num, 150)
            defer patches.Reset()
            So(num, ShouldEqual, 150)
        })

        Convey("recover", func() {
            So(num, ShouldEqual, 10)
        })
    })
}


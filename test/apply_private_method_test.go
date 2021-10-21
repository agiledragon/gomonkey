package test

import (
    "github.com/agiledragon/gomonkey/v2/test/fake"
    "reflect"
    "testing"

    . "github.com/agiledragon/gomonkey/v2"
    . "github.com/smartystreets/goconvey/convey"
)

func TestApplyPrivateMethod(t *testing.T) {
	Convey("TestApplyPrivateMethod", t, func() {
        Convey("patch private pointer method in the different package", func() {
            f := new(fake.PrivateMethodStruct)
            var s *fake.PrivateMethodStruct
            patches := ApplyPrivateMethod(reflect.TypeOf(s), "ok", func(_ *fake.PrivateMethodStruct) bool {
                return false
            })
            defer patches.Reset()
            result := f.Happy()
            So(result, ShouldEqual, "unhappy")
        })

        Convey("patch private value method in the different package", func() {
            s := fake.PrivateMethodStruct{}
            patches := ApplyPrivateMethod(reflect.TypeOf(s), "haveEaten", func(_ fake.PrivateMethodStruct) bool {
                return false
            })
            defer patches.Reset()
            result := s.AreYouHungry()
            So(result, ShouldEqual, "I am hungry")
        })
    })

}


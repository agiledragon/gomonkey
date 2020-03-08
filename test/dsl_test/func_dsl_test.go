package dsl_test

import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/agiledragon/gomonkey/dsl"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
)

func TestPbBuilderFunc(t *testing.T) {
    Convey("TestPbBuilderFunc", t, func() {

        Convey("first dsl", func() {
            patches := NewPatches()
            defer patches.Reset()
            patchBuilder := NewPatchBuilder(patches)

            patchBuilder.
                Func(Belong).
                Stubs().
                With(Eq("zxl"), Any()).
                Will(Return(true)).
                Then(Repeat(Return(false), 2)).
                End()

            flag := Belong("zxl", []string{})
            So(flag, ShouldBeTrue)

            defer func() {
                if p := recover(); p != nil {
                    str, ok := p.(string)
                    So(ok, ShouldBeTrue)
                    So(str, ShouldEqual, "input paras ddd is not matched")
                }
            }()
            Belong("ddd", []string{"abc"})
        })

    })
}

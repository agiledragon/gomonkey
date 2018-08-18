package test

import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "reflect"
    "github.com/agiledragon/gomonkey/test/fake"
)


func TestApplyMethod(t *testing.T) {
    slice := fake.NewSlice()
    var s *fake.Slice
    Convey("TestApplyMethod", t, func() {

        Convey("for succ", func() {
            err := slice.Add(1)
            So(err, ShouldEqual, nil)
            patches := ApplyMethod(reflect.TypeOf(s), "Add", func(_ *fake.Slice, _ int) error {
                return nil
            })
            defer patches.Reset()
            err = slice.Add(1)
            So(err, ShouldEqual, nil)
            err = slice.Remove(1)
            So(err, ShouldEqual, nil)
            So(len(slice), ShouldEqual, 0)
        })

        Convey("for already exist", func() {
            err := slice.Add(2)
            So(err, ShouldEqual, nil)
            patches := ApplyMethod(reflect.TypeOf(s), "Add", func(_ *fake.Slice, _ int) error {
                return fake.ERR_ELEM_EXIST
            })
            defer patches.Reset()
            err = slice.Add(1)
            So(err, ShouldEqual, fake.ERR_ELEM_EXIST)
            err = slice.Remove(2)
            So(err, ShouldEqual, nil)
            So(len(slice), ShouldEqual, 0)
        })

        Convey("two methods", func() {
            err := slice.Add(3)
            So(err, ShouldEqual, nil)
            patches := ApplyMethod(reflect.TypeOf(s), "Add", func(_ *fake.Slice, _ int) error {
                return fake.ERR_ELEM_EXIST
            })
            defer patches.Reset()
            patches.ApplyMethod(reflect.TypeOf(s), "Remove", func(_ *fake.Slice, _ int) error {
                return fake.ERR_ELEM_NT_EXIST
            })
            err = slice.Add(2)
            So(err, ShouldEqual, fake.ERR_ELEM_EXIST)
            err = slice.Remove(1)
            So(err, ShouldEqual, fake.ERR_ELEM_NT_EXIST)
            So(len(slice), ShouldEqual, 1)
            So(slice[0], ShouldEqual, 3)
        })

        Convey("one func and one method", func() {
            err := slice.Add(4)
            So(err, ShouldEqual, nil)
            patches := ApplyFunc(fake.Exec, func(_ string, _ ...string) (string, error) {
                return outputExpect, nil
            })
            defer patches.Reset()
            patches.ApplyMethod(reflect.TypeOf(s), "Remove", func(_ *fake.Slice, _ int) error {
                return fake.ERR_ELEM_NT_EXIST
            })
            output, err := fake.Exec("", "")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, outputExpect)
            err = slice.Remove(1)
            So(err, ShouldEqual, fake.ERR_ELEM_NT_EXIST)
            So(len(slice), ShouldEqual, 2)
            So(slice[0], ShouldEqual, 3)
            So(slice[1], ShouldEqual, 4)
        })
    })
}


package test

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApplyMethod(t *testing.T) {
	slice := fake.NewSlice()
	var s *fake.Slice
	Convey("TestApplyMethod", t, func() {

		Convey("for succ", func() {
			err := slice.Add(1)
			So(err, ShouldEqual, nil)
			patches := ApplyMethod(s, "Add", func(_ *fake.Slice, _ int) error {
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
			patches := ApplyMethod(s, "Add", func(_ *fake.Slice, _ int) error {
				return fake.ErrElemExsit
			})
			defer patches.Reset()
			err = slice.Add(1)
			So(err, ShouldEqual, fake.ErrElemExsit)
			err = slice.Remove(2)
			So(err, ShouldEqual, nil)
			So(len(slice), ShouldEqual, 0)
		})

		Convey("two methods", func() {
			err := slice.Add(3)
			So(err, ShouldEqual, nil)
			defer slice.Remove(3)
			patches := ApplyMethod(s, "Add", func(_ *fake.Slice, _ int) error {
				return fake.ErrElemExsit
			})
			defer patches.Reset()
			patches.ApplyMethod(s, "Remove", func(_ *fake.Slice, _ int) error {
				return fake.ErrElemNotExsit
			})
			err = slice.Add(2)
			So(err, ShouldEqual, fake.ErrElemExsit)
			err = slice.Remove(1)
			So(err, ShouldEqual, fake.ErrElemNotExsit)
			So(len(slice), ShouldEqual, 1)
			So(slice[0], ShouldEqual, 3)
		})

		Convey("one func and one method", func() {
			err := slice.Add(4)
			So(err, ShouldEqual, nil)
			defer slice.Remove(4)
			patches := ApplyFunc(fake.Exec, func(_ string, _ ...string) (string, error) {
				return outputExpect, nil
			})
			defer patches.Reset()
			patches.ApplyMethod(s, "Remove", func(_ *fake.Slice, _ int) error {
				return fake.ErrElemNotExsit
			})
			output, err := fake.Exec("", "")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, outputExpect)
			err = slice.Remove(1)
			So(err, ShouldEqual, fake.ErrElemNotExsit)
			So(len(slice), ShouldEqual, 1)
			So(slice[0], ShouldEqual, 4)
		})
	})
}

package test

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	. "github.com/smartystreets/goconvey/convey"
)

/*
	compare with apply_method_test.go, no need pass receiver
*/

func TestApplyMethodFunc(t *testing.T) {
	slice := fake.NewSlice()
	var s *fake.Slice
	Convey("TestApplyMethodFunc", t, func() {
		Convey("for succ", func() {
			err := slice.Add(1)
			So(err, ShouldEqual, nil)
			patches := ApplyMethodFunc(s, "Add", func(_ int) error {
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
			patches := ApplyMethodFunc(s, "Add", func(_ int) error {
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
			patches := ApplyMethodFunc(s, "Add", func(_ int) error {
				return fake.ErrElemExsit
			})
			defer patches.Reset()
			patches.ApplyMethodFunc(s, "Remove", func(_ int) error {
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
			patches.ApplyMethodFunc(s, "Remove", func(_ int) error {
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

		Convey("for variadic method", func() {
			slice = fake.NewSlice()
			count := slice.Append(1, 2, 3)
			So(count, ShouldEqual, 3)
			patches := ApplyMethodFunc(s, "Append", func(_ ...int) int {
				return 0
			})
			defer patches.Reset()
			count = slice.Append(4, 5, 6)
			So(count, ShouldEqual, 0)
			So(len(slice), ShouldEqual, 3)
		})
	})
}

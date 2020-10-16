package test

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApplyFuncVarSeq(t *testing.T) {
	Convey("TestApplyFuncVarSeq", t, func() {

		Convey("default times is 1", func() {
			info1 := "hello cpp"
			info2 := "hello golang"
			info3 := "hello gomonkey"
			outputs := []OutputCell{
				{Values: Params{[]byte(info1), nil}},
				{Values: Params{[]byte(info2), nil}},
				{Values: Params{[]byte(info3), nil}},
			}
			patches := ApplyFuncVarSeq(&fake.Marshal, outputs)
			defer patches.Reset()
			bytes, err := fake.Marshal("")
			So(err, ShouldEqual, nil)
			So(string(bytes), ShouldEqual, info1)
			bytes, err = fake.Marshal("")
			So(err, ShouldEqual, nil)
			So(string(bytes), ShouldEqual, info2)
			bytes, err = fake.Marshal("")
			So(err, ShouldEqual, nil)
			So(string(bytes), ShouldEqual, info3)
		})

		Convey("retry succ util the third times", func() {
			info1 := "hello cpp"
			outputs := []OutputCell{
				{Values: Params{[]byte(""), fake.ErrActual}, Times: 2},
				{Values: Params{[]byte(info1), nil}},
			}
			patches := ApplyFuncVarSeq(&fake.Marshal, outputs)
			defer patches.Reset()
			bytes, err := fake.Marshal("")
			So(err, ShouldEqual, fake.ErrActual)
			bytes, err = fake.Marshal("")
			So(err, ShouldEqual, fake.ErrActual)
			bytes, err = fake.Marshal("")
			So(err, ShouldEqual, nil)
			So(string(bytes), ShouldEqual, info1)
		})

		Convey("batch operations failed on the third time", func() {
			info1 := "hello gomonkey"
			outputs := []OutputCell{
				{Values: Params{[]byte(info1), nil}, Times: 2},
				{Values: Params{[]byte(""), fake.ErrActual}},
			}
			patches := ApplyFuncVarSeq(&fake.Marshal, outputs)
			defer patches.Reset()
			bytes, err := fake.Marshal("")
			So(err, ShouldEqual, nil)
			So(string(bytes), ShouldEqual, info1)
			bytes, err = fake.Marshal("")
			So(err, ShouldEqual, nil)
			So(string(bytes), ShouldEqual, info1)
			bytes, err = fake.Marshal("")
			So(err, ShouldEqual, fake.ErrActual)
		})

	})
}

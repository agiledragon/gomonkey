package test

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApplyMethodSeq(t *testing.T) {
	e := &fake.Etcd{}
	Convey("TestApplyMethodSeq", t, func() {

		Convey("default times is 1", func() {
			info1 := "hello cpp"
			info2 := "hello golang"
			info3 := "hello gomonkey"
			outputs := []OutputCell{
				{Values: Params{info1, nil}},
				{Values: Params{info2, nil}},
				{Values: Params{info3, nil}},
			}
			patches := ApplyMethodSeq(e, "Retrieve", outputs)
			defer patches.Reset()
			output, err := e.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info1)
			output, err = e.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info2)
			output, err = e.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info3)
		})

		Convey("retry succ util the third times", func() {
			info1 := "hello cpp"
			outputs := []OutputCell{
				{Values: Params{"", fake.ErrActual}, Times: 2},
				{Values: Params{info1, nil}},
			}
			patches := ApplyMethodSeq(e, "Retrieve", outputs)
			defer patches.Reset()
			output, err := e.Retrieve("")
			So(err, ShouldEqual, fake.ErrActual)
			output, err = e.Retrieve("")
			So(err, ShouldEqual, fake.ErrActual)
			output, err = e.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info1)
		})

		Convey("batch operations failed on the third time", func() {
			info1 := "hello gomonkey"
			outputs := []OutputCell{
				{Values: Params{info1, nil}, Times: 2},
				{Values: Params{"", fake.ErrActual}},
			}
			patches := ApplyMethodSeq(e, "Retrieve", outputs)
			defer patches.Reset()
			output, err := e.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info1)
			output, err = e.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info1)
			output, err = e.Retrieve("")
			So(err, ShouldEqual, fake.ErrActual)
		})

	})
}

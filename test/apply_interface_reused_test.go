package test

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApplyInterfaceReused(t *testing.T) {
	e := &fake.Etcd{}

	Convey("TestApplyInterfaceReused", t, func() {
		patches := ApplyFunc(fake.NewDb, func(_ string) fake.Db {
			return e
		})
		defer patches.Reset()
		db := fake.NewDb("mysql")

		Convey("TestApplyInterface", func() {
			info := "hello interface"
			patches.ApplyMethod(e, "Retrieve",
				func(_ *fake.Etcd, _ string) (string, error) {
					return info, nil
				})
			output, err := db.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info)
		})

		Convey("TestApplyInterfaceSeq", func() {
			info1 := "hello cpp"
			info2 := "hello golang"
			info3 := "hello gomonkey"
			outputs := []OutputCell{
				{Values: Params{info1, nil}},
				{Values: Params{info2, nil}},
				{Values: Params{info3, nil}},
			}
			patches.ApplyMethodSeq(e, "Retrieve", outputs)
			output, err := db.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info1)
			output, err = db.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info2)
			output, err = db.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info3)
		})

		Convey("the arg type can be interface", func() {
			info := "hello interface"
			patches.ApplyMethod(e, "Retrieve",
				func(_ fake.Db, _ string) (string, error) {
					return info, nil
				})
			output, err := db.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info)
		})
	})
}

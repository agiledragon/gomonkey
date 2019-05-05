package test

import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
	"reflect"
	"github.com/agiledragon/gomonkey/test/fake"
)

func TestApplyInterfaceReused(t *testing.T) {

    Convey("TestApplyInterfaceReused", t, func() {
		e := &fake.Etcd{}
        Convey("TestApplyInterface", func() {
			patches := ApplyFunc(fake.NewDb, func(_ string) fake.Db {
				return e
			})
			defer patches.Reset()
			patchForInterface := "support a patch for a interface"
			patches.ApplyMethod(reflect.TypeOf(e), "Retrieve", 
				func(_ *fake.Etcd, _ string) (string, error) {
					return patchForInterface, nil
            })
            db := fake.NewDb("mysql")
			output, err := db.Retrieve("")
            So(err, ShouldEqual, nil)
            So(patchForInterface, ShouldEqual, output)
		})
		
		Convey("TestApplyInterfaceSeq", func() {
			e := &fake.Etcd{}
			patches := ApplyFunc(fake.NewDb, func(_ string) fake.Db {
				return e
			})
			defer patches.Reset()
			info1 := "hello cpp"
            info2 := "hello golang"
            info3 := "hello gomonkey"
            outputs := []OutputCell{
                {Values: Params{info1, nil}},
                {Values: Params{info2, nil}},
                {Values: Params{info3, nil}},
            }
			patches.ApplyMethodSeq(reflect.TypeOf(e), "Retrieve", outputs)
			db := fake.NewDb("mysql")
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
    })
}
package test

import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "github.com/agiledragon/gomonkey/test/fake"
    "reflect"
)

func TestApplyMethodSeq(t *testing.T) {
    e := &fake.Etcd{}
    Convey("TestApplyMethodSeq", t, func() {

        Convey("default times is 1", func() {
            info1 := "hello cpp"
            info2 := "hello golang"
            info3 := "hello gomonkey"
            outputs := []Output{
                {Values: Values{info1, nil}},
                {Values: Values{info2, nil}},
                {Values: Values{info3, nil}},
            }
            patches := ApplyMethodSeq(reflect.TypeOf(e), "Retrieve", outputs)
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
            outputs := []Output{
                {Values: Values{"", fake.ErrActual}, Times: 2},
                {Values: Values{info1, nil}},
            }
            patches := ApplyMethodSeq(reflect.TypeOf(e), "Retrieve", outputs)
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
            outputs := []Output{
                {Values: Values{info1, nil}, Times: 2},
                {Values: Values{"", fake.ErrActual}},
            }
            patches := ApplyMethodSeq(reflect.TypeOf(e), "Retrieve", outputs)
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



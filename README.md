# gomonkey

gomonkey is a library to make monkey patching in unit tests easy.

## Features
+ support functions
+ support member methods
+ support a specified sequence values for a function
+ support a specified sequence values for a method

## Notes
+ gomonkey fails to patch a function if inlining is enabled, please running your tests with inlining disabled by adding the command line argument that is `-gcflags=-l`.
+ gomonkey support amd64 temporarily.

## Installation
------------
	$ go get github.com/agiledragon/gomonkey
	
## Using gomonkey

The following just make some tests as typical examples.
**Please refer to the test cases, very complete and detailed.**

### ApplyFunc

```go
import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "github.com/agiledragon/gomonkey/test/fake"
    "encoding/json"
)

var (
    outputExpect = "xxx-vethName100-yyy"
)

func TestApplyFunc(t *testing.T) {
    Convey("TestApplyFunc", t, func() {

        Convey("one func for succ", func() {
            patches := ApplyFunc(fake.Exec, func(_ string, _ ...string) (string, error) {
                    return outputExpect, nil
                })
            defer patches.Reset()
            output, err := fake.Exec("", "")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, outputExpect)
        })

        Convey("one func for fail", func() {
            patches := ApplyFunc(fake.Exec, func(_ string, _ ...string) (string, error) {
                return "", fake.ErrActual
            })
            defer patches.Reset()
            output, err := fake.Exec("", "")
            So(err, ShouldEqual, fake.ErrActual)
            So(output, ShouldEqual, "")
        })

        Convey("two funcs", func() {
            patches := ApplyFunc(fake.Exec, func(_ string, _ ...string) (string, error) {
                return outputExpect, nil
            })
            defer patches.Reset()
            patches.ApplyFunc(fake.Belong, func(_ string, _ []string) bool {
                return true
            })
            output, err := fake.Exec("", "")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, outputExpect)
            flag := fake.Belong("", nil)
            So(flag, ShouldBeTrue)
        })

        Convey("input and output param", func() {
            patches := ApplyFunc(json.Unmarshal, func(_ []byte, v interface{}) error {
                p := v.(*map[int]int)
                *p = make(map[int]int)
                (*p)[1] = 2
                (*p)[2] = 4
                return nil
            })
            defer patches.Reset()
            var m map[int]int
            err := json.Unmarshal(nil, &m)
            So(err, ShouldEqual, nil)
            So(m[1], ShouldEqual, 2)
            So(m[2], ShouldEqual, 4)
        })
    })
}


```

### ApplyFunc

```go
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
            defer slice.Remove(3)
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
            defer slice.Remove(4)
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
            So(len(slice), ShouldEqual, 1)
            So(slice[0], ShouldEqual, 4)
        })
    })
}

```

### ApplyFunc

//TO DO

### ApplyFunc

//TO DO
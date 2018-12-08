# gomonkey

gomonkey is a library to make monkey patching in unit tests easy.

## Features

+ support a patch for a function
+ support a patch for a member method
+ support a patch for a global variable
+ support a patch for a function variable
+ support patches of a specified sequence for a function
+ support patches of a specified sequence for a member method
+ support patches of a specified sequence for a function variable

## Notes
+ gomonkey fails to patch a function or a member method if inlining is enabled, please running your tests with inlining disabled by adding the command line argument that is `-gcflags=-l`(below go1.10) or `-gcflags=all=-l`(go1.10 and above).
+ gomonkey should work on any amd64 system.
+ A panic may happen when a goroutine is patching a function or a member method that is visited by another goroutine at the same time. That is to say, gomonkey is not threadsafe.
+ go1.6 version of the reflection mechanism supports the query of private member methods, but go1.7 and above does not support it. However, all versions of the reflection mechanism support the query of private functions, so gomonkey will trigger a `panic` for only patching a private member method when go1.7 and above is used.

## Installation
```go
$ go get github.com/agiledragon/gomonkey
```
## Using gomonkey

The following just make some tests as idioms.
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

### ApplyMethod

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

## ApplyGlobalVar

```go
import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
)

var num = 10

func TestApplyGlobalVar(t *testing.T) {
    Convey("TestApplyGlobalVar", t, func() {

        Convey("change", func() {
            patches := ApplyGlobalVar(&num, 150)
            defer patches.Reset()
            So(num, ShouldEqual, 150)
        })

        Convey("recover", func() {
            So(num, ShouldEqual, 10)
        })
    })
}
```

## ApplyFuncVar

```go
import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "github.com/agiledragon/gomonkey/test/fake"
)


func TestApplyFuncVar(t *testing.T) {
    Convey("TestApplyFuncVar", t, func() {

        Convey("for succ", func() {
            str := "hello"
            patches := ApplyFuncVar(&fake.Marshal, func (_ interface{}) ([]byte, error) {
                return []byte(str), nil
            })
            defer patches.Reset()
            bytes, err := fake.Marshal(nil)
            So(err, ShouldEqual, nil)
            So(string(bytes), ShouldEqual, str)
        })

        Convey("for fail", func() {
            patches := ApplyFuncVar(&fake.Marshal, func (_ interface{}) ([]byte, error) {
                return nil, fake.ErrActual
            })
            defer patches.Reset()
            _, err := fake.Marshal(nil)
            So(err, ShouldEqual, fake.ErrActual)
        })
    })
}

```


### ApplyFuncSeq

```go
import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "github.com/agiledragon/gomonkey/test/fake"
)

func TestApplyFuncSeq(t *testing.T) {
    Convey("TestApplyFuncSeq", t, func() {

        Convey("default times is 1", func() {
            info1 := "hello cpp"
            info2 := "hello golang"
            info3 := "hello gomonkey"
            outputs := []OutputCell{
                {Values: Params{info1, nil}},
                {Values: Params{info2, nil}},
                {Values: Params{info3, nil}},
            }
            patches := ApplyFuncSeq(fake.ReadLeaf, outputs)
            defer patches.Reset()
            output, err := fake.ReadLeaf("")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, info1)
            output, err = fake.ReadLeaf("")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, info2)
            output, err = fake.ReadLeaf("")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, info3)
        })

        Convey("retry succ util the third times", func() {
            info1 := "hello cpp"
            outputs := []OutputCell{
                {Values: Params{"", fake.ErrActual}, Times: 2},
                {Values: Params{info1, nil}},
            }
            patches := ApplyFuncSeq(fake.ReadLeaf, outputs)
            defer patches.Reset()
            output, err := fake.ReadLeaf("")
            So(err, ShouldEqual, fake.ErrActual)
            output, err = fake.ReadLeaf("")
            So(err, ShouldEqual, fake.ErrActual)
            output, err = fake.ReadLeaf("")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, info1)
        })

        Convey("batch operations failed on the third time", func() {
            info1 := "hello gomonkey"
            outputs := []OutputCell{
                {Values: Params{info1, nil}, Times: 2},
                {Values: Params{"", fake.ErrActual}},
            }
            patches := ApplyFuncSeq(fake.ReadLeaf, outputs)
            defer patches.Reset()
            output, err := fake.ReadLeaf("")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, info1)
            output, err = fake.ReadLeaf("")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, info1)
            output, err = fake.ReadLeaf("")
            So(err, ShouldEqual, fake.ErrActual)
        })

    })
}


```

### ApplyMethodSeq

```go
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
            outputs := []OutputCell{
                {Values: Params{info1, nil}},
                {Values: Params{info2, nil}},
                {Values: Params{info3, nil}},
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
            outputs := []OutputCell{
                {Values: Params{"", fake.ErrActual}, Times: 2},
                {Values: Params{info1, nil}},
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
            outputs := []OutputCell{
                {Values: Params{info1, nil}, Times: 2},
                {Values: Params{"", fake.ErrActual}},
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


```

## ApplyFuncVarSeq

```go
import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "github.com/agiledragon/gomonkey/test/fake"
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

```

## NewPatches

```go
import (
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "github.com/agiledragon/gomonkey/test/fake"
    "encoding/json"
)

func TestPatchPair(t *testing.T) {
   
    Convey("TestPatchPair", t, func() {
        
        Convey("TestPatchPair", func() {
            patchPairs := [][2]interface{} {
                {
                    fake.Exec,
                    func(_ string, _ ...string) (string, error) {
                        return outputExpect, nil
                    },
                },
                {
                    json.Unmarshal,
                    func(_ []byte, v interface{}) error {
                        p := v.(*map[int]int)
                        *p = make(map[int]int)
                        (*p)[1] = 2
                        (*p)[2] = 4
                        return nil
                    },
                },

            }
            patches := NewPatches()
            defer patches.Reset()
            for _, pair := range patchPairs {
                patches.ApplyFunc(pair[0], pair[1])
            }

            output, err := fake.Exec("", "")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, outputExpect)

            var m map[int]int
            err = json.Unmarshal(nil, &m)
            So(err, ShouldEqual, nil)
            So(m[1], ShouldEqual, 2)
            So(m[2], ShouldEqual, 4)
        })
    })
}

```



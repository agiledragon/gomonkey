package test

import (
   // "fmt"
    "errors"
    . "github.com/agiledragon/gomonkey"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "reflect"
)

var (
    ERR_ELEM_EXIST = errors.New("elem already exist")
    ERR_ELEM_NT_EXIST = errors.New("elem not exist")
)

type Slice []int

func NewSlice() Slice {
    return make(Slice, 0)
}

func (this* Slice) Add(elem int) error {
    for _, v := range *this {
        if v == elem {
           // fmt.Printf("Slice: Add elem: %v already exist\n", elem)
            return ERR_ELEM_EXIST
        }
    }
    *this = append(*this, elem)
    //fmt.Printf("Slice: Add elem: %v succ\n", elem)
    return nil
}

func (this* Slice) Remove(elem int) error {
    found := false
    for i, v := range *this {
        if v == elem {
            if i == len(*this) - 1 {
                *this = (*this)[:i]

            } else {
                *this = append((*this)[:i], (*this)[i+1:]...)
            }
            found = true
            break
        }
    }
    if !found {
        //fmt.Printf("Slice: Remove elem: %v not exist\n", elem)
        return ERR_ELEM_NT_EXIST
    }
    //fmt.Printf("Slice: Remove elem: %v succ\n", elem)
    return nil
}

func TestApplyMethod(t *testing.T) {
    slice := NewSlice()
    var s *Slice
    Convey("TestApplyMethod", t, func() {

        Convey("for succ", func() {
            err := slice.Add(1)
            So(err, ShouldEqual, nil)
            patches := ApplyMethod(reflect.TypeOf(s), "Add", func(_ *Slice, _ int) error {
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
            patches := ApplyMethod(reflect.TypeOf(s), "Add", func(_ *Slice, _ int) error {
                return ERR_ELEM_EXIST
            })
            defer patches.Reset()
            err = slice.Add(1)
            So(err, ShouldEqual, ERR_ELEM_EXIST)
            err = slice.Remove(2)
            So(err, ShouldEqual, nil)
            So(len(slice), ShouldEqual, 0)
        })

        Convey("two methods", func() {
            err := slice.Add(3)
            So(err, ShouldEqual, nil)
            patches := ApplyMethod(reflect.TypeOf(s), "Add", func(_ *Slice, _ int) error {
                return ERR_ELEM_EXIST
            })
            defer patches.Reset()
            patches.ApplyMethod(reflect.TypeOf(s), "Remove", func(_ *Slice, _ int) error {
                return ERR_ELEM_NT_EXIST
            })
            err = slice.Add(2)
            So(err, ShouldEqual, ERR_ELEM_EXIST)
            err = slice.Remove(1)
            So(err, ShouldEqual, ERR_ELEM_NT_EXIST)
            So(len(slice), ShouldEqual, 1)
            So(slice[0], ShouldEqual, 3)
        })

        Convey("one func and one method", func() {
            err := slice.Add(4)
            So(err, ShouldEqual, nil)
            patches := ApplyFunc(Exec, func(_ string, _ ...string) (string, error) {
                return outputExpect, nil
            })
            defer patches.Reset()
            patches.ApplyMethod(reflect.TypeOf(s), "Remove", func(_ *Slice, _ int) error {
                return ERR_ELEM_NT_EXIST
            })
            output, err := Exec("", "")
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, outputExpect)
            err = slice.Remove(1)
            So(err, ShouldEqual, ERR_ELEM_NT_EXIST)
            So(len(slice), ShouldEqual, 2)
            So(slice[0], ShouldEqual, 3)
            So(slice[1], ShouldEqual, 4)
        })
    })
}


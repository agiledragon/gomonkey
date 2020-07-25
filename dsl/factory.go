package dsl

import . "github.com/henrylee2cn/gomonkey"

func Any() Constraint {
    return &AnyConstraint{}
}

func Eq(x interface{}) Constraint {
    return &EqConstraint{x: x}
}

func Return(x ...interface{}) Behavior {
    r := &ReturnBehavior{rets: make([]Params, 0), params: make(Params, 0)}
    r.params = append(r.params, x...)
    return r
}

func Repeat(behavior Behavior, times int) Behavior {
    return &RepeatBehavior{rets: make([]Params, 0), behavior: behavior, times: times}
}



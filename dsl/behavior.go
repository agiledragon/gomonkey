package dsl

import . "github.com/henrylee2cn/gomonkey"

type Behavior interface {
    Apply() []Params
}

type ReturnBehavior struct {
    rets   []Params
    params Params
}

func (this *ReturnBehavior) Apply() []Params {
    this.rets = append(this.rets, this.params)
    return this.rets
}

type RepeatBehavior struct {
    rets     []Params
    behavior Behavior
    times    int
}

func (this *RepeatBehavior) Apply() []Params {
    for i := 0; i < this.times; i++ {
        this.rets = append(this.rets, this.behavior.Apply()[0])
    }
    return this.rets
}

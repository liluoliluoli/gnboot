package exp_util

import "fmt"

type TernaryExpression[T any] struct {
	b bool
}

func (t TernaryExpression[T]) Then(r T) TernaryExpression[T] {
	if t.b {
		panic(r)
	}
	return t
}

func (t TernaryExpression[T]) Else(r T) T {
	return r
}
func Ternary[T any](f func() T) (r T) {
	defer func() {
		if e := recover(); e != nil {
			r = e.(T)
		}
	}()
	r = f()
	return
}

func If[T any](b bool) TernaryExpression[T] {
	return TernaryExpression[T]{b: b}
}

func main() {
	x := Ternary(func() float64 { return If[float64](3 > 2).Then(3).Else(2) })
	fmt.Println("r", x)
}

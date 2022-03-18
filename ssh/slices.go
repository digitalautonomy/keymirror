package ssh

import "golang.org/x/exp/constraints"

func transform[T, U any](l []T, f func(T) U) []U {
	output := []U{}
	for _, l := range l {
		output = append(output, f(l))
	}
	return output
}

func isEmptySlice[T any](v []T) bool {
	return len(v) == 0
}

func isEmpty[T comparable](v T) bool {
	var zero T

	// IDEA does not handle this correctly yet. But Golang 1.18 does
	return zero == v
}

func isGreaterThan[T constraints.Ordered](cutoff T) predicate[T] {
	return func(v T) bool {
		return v > cutoff
	}
}

func isEqualTo[T comparable](expected T) predicate[T] {
	return func(v T) bool {
		return expected == v
	}
}

func isNil[T any](v *T) bool {
	return v == nil
}

type func1[T, U any] func(T) U
type predicate[T any] func1[T, bool]

func not[T any](f predicate[T]) predicate[T] {
	return func(v T) bool {
		return !f(v)
	}
}

func filter[T any](v []T, f predicate[T]) []T {
	res := []T{}

	for _, t := range v {
		if f(t) {
			res = append(res, t)
		}
	}
	return res
}

func ignoringErrors[T, R any](f func(T) (R, error)) func(T) R {
	return func(s T) R {
		b, _ := f(s)
		return b
	}
}

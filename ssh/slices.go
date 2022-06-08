package ssh

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func appendingWith[T, U any](f func1[T, U]) func([]U, T) []U {
	return func(previous []U, current T) []U {
		return append(previous, f(current))
	}
}

func appendingWhen[T any](f predicate[T]) func([]T, T) []T {
	return func(previous []T, current T) []T {
		if f(current) {
			return append(previous, current)
		}
		return previous
	}
}

func transform[T, U any](l []T, f func1[T, U]) []U {
	return foldLeft(l, []U{}, appendingWith(f))
}

func filter[T any](v []T, f predicate[T]) []T {
	return foldLeft(v, []T{}, appendingWhen(f))
}

func foldLeft[T any, R any](l []T, start R, f func(R, T) R) R {
	result := start

	for _, current := range l {
		result = f(result, current)
	}

	return result
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

func loggingErrors[T, R any](l logrus.FieldLogger, message string, f func(T) (R, error)) func(T) R {
	return func(s T) R {
		b, e := f(s)
		if e != nil {
			l.WithError(e).Error(message)
		}
		return b
	}
}

func ignoringErrors[T, R any](f func(T) (R, error)) func(T) R {
	return func(s T) R {
		b, _ := f(s)
		return b
	}
}

func existsIn[T comparable](l []T) predicate[T] {
	return func(e T) bool {
		return slices.Contains(l, e)
	}
}

func foreachValue[K comparable, V any](values map[K]V, f func(V)) {
	for _, v := range values {
		f(v)
	}
}

func foreach[T any](values []T, f func(T)) {
	for _, v := range values {
		f(v)
	}
}

func concat[T any](s ...[]T) []T {
	if len(s) == 0 {
		return nil
	}

	return append(s[0], concat(s[1:]...)...)
}

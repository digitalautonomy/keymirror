package collections

import "golang.org/x/exp/maps"

// Set represents an unordered collection of objects.
// Sets can be compared for equality, and the elements manipulated
// in the expected ways
type Set[T comparable] interface {
	Add(T)
	AddAll(...T)
	Intersection(Set[T]) Set[T]
	Union(Set[T]) Set[T]
	Difference(Set[T]) Set[T]
	SubSetOf(Set[T]) bool
	Equal(Set[T]) bool
	Contains(T) bool
	Size() int
	Empty() bool
	ToSlice() []T
}

type simpleSet[T comparable] struct {
	values map[T]bool
}

// SetFrom will create a set containing the given values
func SetFrom[T comparable](values ...T) Set[T] {
	s := NewSet[T]()
	s.AddAll(values...)
	return s
}

// NewSet creates a new empty set
func NewSet[T comparable]() Set[T] {
	return &simpleSet[T]{
		values: map[T]bool{},
	}
}

// Add implements the Set interface for simpleSet
func (s *simpleSet[T]) Add(v T)  {
	s.values[v] = true
}

// Size implements the Set interface for simpleSet
func (s *simpleSet[T]) Size() int {
	return len(s.values)
}

// Empty implements the Set interface for simpleSet
func (s *simpleSet[T]) Empty() bool {
	return s.Size() == 0
}

// AddAll implements the Set interface for simpleSet
func (s *simpleSet[T]) AddAll(values ...T)  {
	for _, v := range values {
		s.Add(v)
	}
}

// Intersection implements the Set interface for simpleSet
func (s *simpleSet[T]) Intersection(other Set[T]) Set[T] {
	newSet := NewSet[T]()

	for k := range s.values {
		if other.Contains(k) {
			newSet.Add(k)
		}
	}
	
	return newSet
}

// Difference implements the Set interface for simpleSet
func (s *simpleSet[T]) Difference(other Set[T]) Set[T] {
	newSet := NewSet[T]()

	for k := range s.values {
		if !other.Contains(k) {
			newSet.Add(k)
		}
	}
	
	return newSet
}

// SubSetOf implements the Set interface for simpleSet
func (s *simpleSet[T]) SubSetOf(other Set[T]) bool {
	return s.Difference(other).Empty()
}

// Union implements the Set interface for simpleSet
func (s *simpleSet[T]) Union(other Set[T]) Set[T] {
	newSet := SetFrom(s.ToSlice()...)
	newSet.AddAll(other.ToSlice()...)
	return newSet
}

// Contains implements the Set interface for simpleSet
func (s *simpleSet[T]) Contains(v T) bool {
	return s.values[v]
}

// ToSlice implements the Set interface for simpleSet
func (s *simpleSet[T]) ToSlice() []T {
	return maps.Keys(s.values)
}

// Equal implements the Set interface for simpleSet
func (s *simpleSet[T]) Equal(other Set[T]) bool {
	return s.Difference(other).Empty() && other.Difference(s).Empty()
}

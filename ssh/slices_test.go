package ssh

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type genericsSuite struct {
	suite.Suite
}

func TestGenericsSuite(t *testing.T) {
	suite.Run(t, new(genericsSuite))
}

func (s *genericsSuite) Test_isEmpty() {
	s.True(isEmpty(false))
	s.False(isEmpty(true))

	s.True(isEmpty(0))
	s.False(isEmpty(1))

	s.True(isEmpty(""))
	s.False(isEmpty("fo"))
}

func (s *genericsSuite) Test_isEmptySlice() {
	s.True(isEmptySlice([]string{}))
	s.False(isEmptySlice([]string{""}))
}

func (s *genericsSuite) Test_isNil() {
	var t1 *string
	var t2 *int
	var t3 *string = new(string)
	var t4 *int = new(int)

	s.True(isNil(t1))
	s.True(isNil(t2))
	s.False(isNil(t3))
	s.False(isNil(t4))
}

func (s *genericsSuite) Test_not() {
	f1 := func(s string) bool {
		return s == "foo" || s == "bar"
	}

	s.True(f1("foo"))
	s.True(f1("bar"))
	s.False(f1("quux"))

	inv := not(f1)
	s.False(inv("foo"))
	s.False(inv("bar"))
	s.True(inv("quux"))

	f2 := func(s int) bool {
		return s == 42
	}
	inv2 := not(f2)
	s.False(inv2(42))
}

func (s *genericsSuite) Test_filter() {
	l := []string{"a", "", "b"}
	s.Equal([]string{""}, filter(l, isEmpty[string]))
	s.Equal([]string{"a", "b"}, filter(l, not(isEmpty[string])))
	s.Equal([]int{43, 43}, filter([]int{41, 42, 43, 44, 43}, isEqualTo(43)))
	s.Equal([]int{43, 44, 43}, filter([]int{41, 42, 43, 44, 43}, isGreaterThan(42)))
}

func (s *genericsSuite) Test_foldLeft_sum() {
	plus := func(l, r int) int {
		return l + r
	}

	input := []int{2, 3, 5, 9, 13}
	s.Equal(32, foldLeft(input, 0, plus))
}

func (s *genericsSuite) Test_foldLeft_product() {
	mult := func(l, r int) int {
		return l * r
	}

	input := []int{2, 3, 5, 9}
	s.Equal(270, foldLeft(input, 1, mult))
}

func (s *genericsSuite) Test_existsIn_returnsAPredicateCheckingForExistance() {
	p1 := existsIn([]string{"foo", "bar"})
	p2 := existsIn([]int{1, 2, 4, 8})

	s.False(p1("hello"))
	s.True(p1("foo"))
	s.True(p1("bar"))

	s.False(p2(0))
	s.True(p2(1))
	s.True(p2(2))
	s.False(p2(3))
	s.True(p2(4))
}

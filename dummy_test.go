package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type dummySuite struct {
	suite.Suite
}

func TestDummySuite(t *testing.T) {
	suite.Run(t, new(dummySuite))
}

func (s *dummySuite) Test_Dummy() {
	s.Equal(1, 1)
}

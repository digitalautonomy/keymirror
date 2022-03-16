package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type mainSuite struct {
	suite.Suite
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(mainSuite))
}

func (s *mainSuite) Test_AlwaysSucceeds() {
	main()
	s.Equal(1, 1)
}

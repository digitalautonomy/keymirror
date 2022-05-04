package i18n

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type localizeSuite struct {
	suite.Suite
}

func TestLocalizeSuite(t *testing.T) {
	suite.Run(t, new(localizeSuite))
}

func (s *localizeSuite) Test_Local_ReturnsTheStringPassedAsParameter() {
	m := "this string should be returned without being modified..."
	s.Equal(Local(m), m)

	m = "...same here!"
	s.Equal(Local(m), m)
}

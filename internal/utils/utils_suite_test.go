package utils_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UtilsTestSuite))
}

package data_models

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type ChroniclesWindowSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestChroniclesWindow(t *testing.T) {
	suite.Run(t, new(ChroniclesWindowSuite))
}

func (suite *ChroniclesWindowSuite) TestEOFErrorIs() {
	suite.Nil(nil)
	err := &ScanHttpErrorRetry{errors.New("Some error")}
	otherErr := &ScanHttpErrorRetry{}
	suite.True(errors.Is(err, otherErr))
}

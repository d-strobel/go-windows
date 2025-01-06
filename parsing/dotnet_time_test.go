package parsing

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all Time parsing functions
type DotnetTimeUnitTestSuite struct {
	suite.Suite
	// Fixtures
	dotNetDatetime          string
	expectedDatetime        time.Time
	jsonData                string
	expectedUnmarshaledJSON TestUnmarshalObject
}

// Test struct for unmarshal tests
type TestUnmarshalObject struct {
	Name    string     `json:"Name"`
	Created DotnetTime `json:"Created"`
	Updated DotnetTime `json:"Updated"`
}

func (suite *DotnetTimeUnitTestSuite) SetupSuite() {
	// Fixtures
	suite.dotNetDatetime = `"\/Date(1701379505092)\/"`
	suite.expectedDatetime = time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)
	suite.jsonData = `{
    "Name":  "tester",
    "Created":  "\/Date(1701379505092)\/",
    "Updated":  "\/Date(1701379505092)\/"
}
	`
	suite.expectedUnmarshaledJSON = TestUnmarshalObject{
		Name:    "tester",
		Created: DotnetTime{Time: suite.expectedDatetime},
		Updated: DotnetTime{Time: suite.expectedDatetime},
	}
}

func TestDotnetTimeUnitTestSuite(t *testing.T) {
	suite.Run(t, &DotnetTimeUnitTestSuite{})
}

func (suite *DotnetTimeUnitTestSuite) TestUnmarshalJSON() {
	suite.T().Parallel()

	suite.Run("should unmarshal the dotnet timestring to DotnetTime object", func() {
		winTime := DotnetTime{}
		expectedDotnetTime := DotnetTime{Time: suite.expectedDatetime}
		err := winTime.UnmarshalJSON([]byte(suite.dotNetDatetime))
		suite.NoError(err)
		suite.Equal(expectedDotnetTime, winTime)
	})

	suite.Run("should return error with invalid dotnet timestring", func() {
		winTime := DotnetTime{}
		err := winTime.UnmarshalJSON([]byte("2023-20-10"))
		suite.ErrorContains(err, "parsing.UnmarshalJSON(DotnetTime): input string is not a dotnet JSON datetime")
	})

	suite.Run("should unmarshal the whole json object correctly", func() {
		actualResult := TestUnmarshalObject{}
		err := json.Unmarshal([]byte(suite.jsonData), &actualResult)
		suite.NoError(err)
		suite.Equal(suite.expectedUnmarshaledJSON, actualResult)
	})
}

func (suite *DotnetTimeUnitTestSuite) TestTimeMethods() {
	suite.Run("should be able to format dotnet time to RFC3389", func() {
		dotnetTime := DotnetTime{Time: time.Date(2023, time.November, 30, 21, 25, 5, 0, time.UTC)}
		actualResult := dotnetTime.Format(time.RFC3339)
		expectedResult := "2023-11-30T21:25:05Z"
		suite.Equal(expectedResult, actualResult)
	})
}

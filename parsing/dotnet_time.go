package parsing

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// DotnetTime is a custom time type that embeds the time.Time type. It is designed to handle
// the unmarshalling of dotnet JSON datetime strings in the format "\"/Date(timestamp)/\""
// when used as a field in a struct that is being unmarshalled from JSON.
type DotnetTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface for the DotnetTime type.
// It parses a JSON-encoded dotnet JSON datetime string and converts it into a DotnetTime object.
// The input byte slice is expected to represent a dotnet JSON datetime string in the format "\"/Date(timestamp)/\"".
func (t *DotnetTime) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package
	if string(b) == "null" || string(b) == `""` {
		return nil
	}

	// Check for valid dotnet timestring
	re, err := regexp.Compile(`^"\\/Date\(\d+\)\\/"$`)
	if err != nil {
		return fmt.Errorf("parsing.UnmarshalJSON(DotnetTime): %s", err)
	}

	if !re.Match(b) {
		return fmt.Errorf("parsing.UnmarshalJSON(DotnetTime): input string is not a dotnet JSON datetime: %s", string(b))
	}

	// Extract timestamp
	re, err = regexp.Compile(`\d+`)
	if err != nil {
		return fmt.Errorf("parsing.UnmarshalJSON(DotnetTime): %s", err)
	}
	timestamp := re.Find(b)

	// Convert to seconds
	i, err := strconv.Atoi(string(timestamp))
	if err != nil {
		return fmt.Errorf("parsing.UnmarshalJSON(DotnetTime): %s", err)
	}
	seconds := int64(i / 1000)

	// Unmarshal unix time into DotnetTime object
	unixTime := time.Unix(seconds, 0).UTC()
	*t = DotnetTime{unixTime}

	return nil
}

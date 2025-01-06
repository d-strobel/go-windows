package parsing

import (
	"encoding/json"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCimIpAddressUnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("ValidIPv4Address", func(t *testing.T) {
		// Valid integer representation of 10.100.91.51
		jsonData := `861627402`
		var ip CimIpAddress

		err := json.Unmarshal([]byte(jsonData), &ip)
		require.NoError(t, err)

		expectedIP := CimIpAddress{
			netip.MustParseAddr("10.100.91.51"),
		}
		assert.Equal(t, expectedIP, ip)
	})

	t.Run("InvalidJSONFormat", func(t *testing.T) {
		// Invalid JSON (string instead of integer)
		jsonData := `"notAnInteger"`
		var ip CimIpAddress

		err := json.Unmarshal([]byte(jsonData), &ip)
		assert.Error(t, err)
	})

	t.Run("ZeroAddress", func(t *testing.T) {
		// Test edge case with 0
		jsonData := `0`
		var ip CimIpAddress

		err := json.Unmarshal([]byte(jsonData), &ip)
		require.NoError(t, err)

		expectedIP := CimIpAddress{
			netip.MustParseAddr("0.0.0.0"),
		}
		assert.Equal(t, expectedIP, ip)
	})

	t.Run("MaxIPv4Address", func(t *testing.T) {
		// Test edge case with maximum IPv4 address (255.255.255.255)
		jsonData := `4294967295` // Equivalent to 255.255.255.255
		var ip CimIpAddress

		err := json.Unmarshal([]byte(jsonData), &ip)
		require.NoError(t, err)

		expectedIP := CimIpAddress{
			netip.MustParseAddr("255.255.255.255"),
		}
		assert.Equal(t, expectedIP, ip)
	})
}

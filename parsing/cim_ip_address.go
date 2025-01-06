package parsing

import (
	"encoding/json"
	"fmt"
	"net/netip"
)

// CimIpAddress represents an IP Address value.
type CimIpAddress struct {
	netip.Addr
}

// UnmarshalJSON implements the json.Unmarshaler interface for the CimIpAddress type.
func (a *CimIpAddress) UnmarshalJSON(b []byte) error {
	// Parse the integer from the JSON input
	var addr uint32
	if err := json.Unmarshal(b, &addr); err != nil {
		return fmt.Errorf("parsing.UnmarshalJSON(CimIpAddress): failed to parse IP address from JSON: %w", err)
	}

	// Convert integer to IPv4 address
	ip := [4]byte{
		byte(addr),
		byte(addr >> 8),
		byte(addr >> 16),
		byte(addr >> 24),
	}

	// Use netip to create the address
	parsedAddr := netip.AddrFrom4(ip)

	// Assign the parsed address to the pointer receiver
	*a = CimIpAddress{
		parsedAddr,
	}

	return nil
}

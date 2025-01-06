package dhcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"github.com/d-strobel/gowindows/winerror"
)

// ExclusionRangeV4 represents an IPv4 DHCP exclusion range.
type ExclusionRangeV4 struct {
	ScopeId    addressString `json:"ScopeId"`
	StartRange addressString `json:"StartRange"`
	EndRange   addressString `json:"EndRange"`
}

// ExclusionRangeV4ReadParams represents parameters for the ipv4 exclusion range read function.
type ExclusionRangeV4ReadParams struct {
	// Specifies the end IP address of the range that is excluded.
	EndRange netip.Addr

	// Specifies the scope ID, in IPv4 address format, from which the excluded IP address range is returned.
	ScopeId netip.Addr

	// Specifies the starting IP address of the range that is excluded.
	StartRange netip.Addr
}

// pwshCommand returns the PowerShell command to read an IPV4 DHCP exclusion range.
func (params ExclusionRangeV4ReadParams) pwshCommand() string {
	return fmt.Sprintf(
		"Get-DhcpServerv4ExclusionRange -ScopeId '%s' | Where-Object {$_.StartRange.IPAddressToString -eq '%s' -and $_.EndRange.IPAddressToString -eq '%s'} | ConvertTo-Json -Compress",
		params.ScopeId,
		params.StartRange,
		params.EndRange,
	)
}

// ExclusionRangeV4Read gets a DHCP exclusion range. It returns a ExclusionRangeV4 object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) ExclusionRangeV4Read(ctx context.Context, params ExclusionRangeV4ReadParams) (ExclusionRangeV4, error) {
	var s ExclusionRangeV4

	// Assert needed parameters
	if !params.ScopeId.Is4() || !params.StartRange.Is4() || !params.EndRange.Is4() {
		return s, errors.New("windows.dhcp.ExclusionRangeV4Read: scope parameters 'ScopeId', 'StartRange' and 'EndRange' must be a valid IPv4 address")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &s); err != nil {
		return s, winerror.Errorf(cmd, "windows.dhcp.ExclusionRangeV4Read: %s", err)
	}

	// If the output of the command is empty, return an error.
	if !s.ScopeId.Address.Is4() {
		return s, winerror.Errorf(cmd, "windows.dhcp.ExclusionRangeV4Read: exclusion range not found")
	}

	return s, nil
}

// ExclusionRangeV4CreateParams represents parameters for the IPv4 exclusion range create function.
type ExclusionRangeV4CreateParams struct {
	// Specifies the ending IP address of the range in the subnet from which IP addresses should
	// be leased by the DHCP server service.
	EndRange netip.Addr

	// Specifies the identifier (ID) of the IPv4 scope from which the IP addresses are excluded.
	ScopeId netip.Addr

	// Specifies the starting IP address of the range in the subnet from which IP addresses should be leased
	// by the DHCP server service.
	StartRange netip.Addr
}

// pwshCommand returns the PowerShell command to create an IPv4 exclusion range.
func (params ExclusionRangeV4CreateParams) pwshCommand() string {
	return fmt.Sprintf(
		"Add-DhcpServerv4ExclusionRange -PassThru -Confirm:$false -ScopeId '%s' -StartRange '%s' -EndRange '%s' | ConvertTo-Json -Compress",
		params.ScopeId,
		params.StartRange,
		params.EndRange,
	)
}

// ExclusionRangeV4Create creates a new IPv4 exclusion range. It returns a ExclusionRangeV4 object.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) ExclusionRangeV4Create(ctx context.Context, params ExclusionRangeV4CreateParams) (ExclusionRangeV4, error) {
	var s ExclusionRangeV4

	// Assert needed parameters
	if !params.ScopeId.Is4() || !params.StartRange.Is4() || !params.EndRange.Is4() {
		return s, errors.New("windows.dhcp.ExclusionRangeV4Create: exclusion range parameter 'ScopeId', 'StartRange' and 'EndRange' must be a valid IPv4 address")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &s); err != nil {
		return s, winerror.Errorf(cmd, "windows.dhcp.ExclusionRangeV4Create: %s", err)
	}

	return s, nil
}

// ExclusionRangeV4DeleteParams represents parameters for the IPv4 exclusion range delete function.
type ExclusionRangeV4DeleteParams struct {
	// Specifies the ending IPv4 address of the excluded IP range which is to be deleted.
	EndRange netip.Addr

	// Specifies the scope identifier (ID), in IPv4 address format, from which the exclusion range is to be deleted.
	ScopeId netip.Addr

	// Specifies the starting IPv4 address of the excluded IP range which is to be deleted.
	StartRange netip.Addr
}

// pwshCommand returns the PowerShell command to delete a DHCP scope.
func (params ExclusionRangeV4DeleteParams) pwshCommand() string {
	return fmt.Sprintf(
		"Remove-DhcpServerv4ExclusionRange -Confirm:$false -ScopeId '%s' -StartRange '%s' -EndRange '%s'",
		params.ScopeId,
		params.StartRange,
		params.EndRange,
	)
}

// ExclusionRangeV4Delete removes an IPv4 exclusion range.
// It returns a *winerror.WinError if the windows client returns an error.
func (c *Client) ExclusionRangeV4Delete(ctx context.Context, params ExclusionRangeV4DeleteParams) error {
	var s ExclusionRangeV4

	// Assert needed parameters
	if !params.ScopeId.Is4() || !params.StartRange.Is4() || !params.EndRange.Is4() {
		return errors.New("windows.dhcp.ExclusionRangeV4Delete: exclusion range parameter 'ScopeId', 'StartRange' and 'EndRange' must be a valid IPv4 address")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &s); err != nil {
		return winerror.Errorf(cmd, "windows.dhcp.ExclusionRangeV4Delete: %s", err)
	}

	return nil
}

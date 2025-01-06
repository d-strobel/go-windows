package dhcp

import (
	"context"
	"net/netip"

	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
)

// Fixtures
const (
	// exclusion range json output
	exclusionRangeV4Json = `{"ScopeId":{"Address":698560,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.0"},"SubnetMask":{"Address":16777215,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"255.255.255.0"},"StartRange":{"Address":84584640,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.5"},"EndRange":{"Address":168470720,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.10"},"ActivatePolicies":true,"Delay":0,"Description":"Test description","LeaseDuration":{"Ticks":6912000000000,"Days":8,"Hours":0,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":8,"TotalHours":192,"TotalMilliseconds":691200000,"TotalMinutes":11520,"TotalSeconds":691200},"MaxBootpClients":4294967295,"Name":"test","NapEnable":false,"NapProfile":"","State":"Active","SuperscopeName":"","Type":"Dhcp","PSComputerName":null}`
)

var (
	expectedExclusionRangeV4 = ExclusionRangeV4{
		ScopeId: addressString{
			Address: netip.MustParseAddr("192.168.10.0"),
		},
		StartRange: addressString{
			Address: netip.MustParseAddr("192.168.10.5"),
		},
		EndRange: addressString{
			netip.MustParseAddr("192.168.10.10"),
		},
	}
)

// Test ExclusionRangeV4 related methods.
func (suite *DhcpServerUnitTestSuite) TestExclusionRangeV4ReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ExclusionRangeV4ReadParams
			expectedCmd     string
		}{
			{
				"assert correct command ExclusionRangeV4",
				ExclusionRangeV4ReadParams{ScopeId: netip.MustParseAddr("192.168.10.0"), StartRange: netip.MustParseAddr("192.168.10.5"), EndRange: netip.MustParseAddr("192.168.10.10")},
				"Get-DhcpServerv4ExclusionRange -ScopeId '192.168.10.0' | Where-Object {$_.StartRange.IPAddressToString -eq '192.168.10.5' -and $_.EndRange.IPAddressToString -eq '192.168.10.10'} | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestExclusionRangeV4Read() {
	suite.Run("should return the correct ExclusionRangeV4", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-DhcpServerv4ExclusionRange -ScopeId '192.168.10.0' | Where-Object {$_.StartRange.IPAddressToString -eq '192.168.10.5' -and $_.EndRange.IPAddressToString -eq '192.168.10.10'} | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: exclusionRangeV4Json}, nil)
		actualExclustionRangeV4, err := c.ExclusionRangeV4Read(ctx, ExclusionRangeV4ReadParams{
			ScopeId:    netip.MustParseAddr("192.168.10.0"),
			StartRange: netip.MustParseAddr("192.168.10.5"),
			EndRange:   netip.MustParseAddr("192.168.10.10"),
		})
		suite.NoError(err)
		suite.Equal(expectedExclusionRangeV4, actualExclustionRangeV4)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ExclusionRangeV4ReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ExclusionRangeV4ReadParams{},
				"windows.dhcp.ExclusionRangeV4Read: scope parameters 'ScopeId', 'StartRange' and 'EndRange' must be a valid IPv4 address",
			},
			{
				"assert error with ipv6",
				ExclusionRangeV4ReadParams{ScopeId: netip.MustParseAddr("fe80:0010::"), StartRange: netip.MustParseAddr("192.168.10.5"), EndRange: netip.MustParseAddr("192.168.10.10")},
				"windows.dhcp.ExclusionRangeV4Read: scope parameters 'ScopeId', 'StartRange' and 'EndRange' must be a valid IPv4 address",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnection(suite.T())
			c := &Client{
				Connection:      mockConn,
				decodeCliXmlErr: func(s string) (string, error) { return "", nil },
			}
			_, err := c.ExclusionRangeV4Read(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test ExclusionRangeV4Create related methods.
func (suite *DhcpServerUnitTestSuite) TestExclusionRangeV4CreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ExclusionRangeV4CreateParams
			expectedCmd     string
		}{
			{
				"assert correct command with necessary parameters",
				ExclusionRangeV4CreateParams{
					ScopeId:    netip.MustParseAddr("192.168.10.0"),
					StartRange: netip.MustParseAddr("192.168.10.5"),
					EndRange:   netip.MustParseAddr("192.168.10.10"),
				},
				"Add-DhcpServerv4ExclusionRange -PassThru -Confirm:$false -ScopeId '192.168.10.0' -StartRange '192.168.10.5' -EndRange '192.168.10.10' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestExclusionRangeV4Create() {
	suite.Run("should return the correct ExclusionRangeV4", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Add-DhcpServerv4ExclusionRange -PassThru -Confirm:$false -ScopeId '192.168.10.0' -StartRange '192.168.10.5' -EndRange '192.168.10.10' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: exclusionRangeV4Json}, nil)
		actualExclusionRangeV4, err := c.ExclusionRangeV4Create(ctx, ExclusionRangeV4CreateParams{
			ScopeId:    netip.MustParseAddr("192.168.10.0"),
			StartRange: netip.MustParseAddr("192.168.10.5"),
			EndRange:   netip.MustParseAddr("192.168.10.10"),
		})
		suite.NoError(err)
		suite.Equal(expectedExclusionRangeV4, actualExclusionRangeV4)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ExclusionRangeV4CreateParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ExclusionRangeV4CreateParams{},
				"windows.dhcp.ExclusionRangeV4Create: exclusion range parameter 'ScopeId', 'StartRange' and 'EndRange' must be a valid IPv4 address",
			},
			{
				"assert error with ipv6",
				ExclusionRangeV4CreateParams{ScopeId: netip.MustParseAddr("fe80:0010::"), StartRange: netip.MustParseAddr("192.168.10.5"), EndRange: netip.MustParseAddr("192.168.10.10")},
				"windows.dhcp.ExclusionRangeV4Create: exclusion range parameter 'ScopeId', 'StartRange' and 'EndRange' must be a valid IPv4 address",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnection(suite.T())
			c := &Client{
				Connection:      mockConn,
				decodeCliXmlErr: func(s string) (string, error) { return "", nil },
			}
			_, err := c.ExclusionRangeV4Create(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test ExclusionRangeV4Delete related methods.
func (suite *DhcpServerUnitTestSuite) TestExclusionRangeV4DeletePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ExclusionRangeV4DeleteParams
			expectedCmd     string
		}{
			{
				"assert correct command with necessary parameters",
				ExclusionRangeV4DeleteParams{
					ScopeId:    netip.MustParseAddr("192.168.10.0"),
					StartRange: netip.MustParseAddr("192.168.10.5"),
					EndRange:   netip.MustParseAddr("192.168.10.10"),
				},
				"Remove-DhcpServerv4ExclusionRange -Confirm:$false -ScopeId '192.168.10.0' -StartRange '192.168.10.5' -EndRange '192.168.10.10'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestExclusionRangeV4Delete() {
	suite.Run("should not error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Remove-DhcpServerv4ExclusionRange -Confirm:$false -ScopeId '192.168.10.0' -StartRange '192.168.10.5' -EndRange '192.168.10.10'").
			Return(connection.CmdResult{}, nil)
		err := c.ExclusionRangeV4Delete(ctx, ExclusionRangeV4DeleteParams{
			ScopeId:    netip.MustParseAddr("192.168.10.0"),
			StartRange: netip.MustParseAddr("192.168.10.5"),
			EndRange:   netip.MustParseAddr("192.168.10.10"),
		})
		suite.NoError(err)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ExclusionRangeV4DeleteParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ExclusionRangeV4DeleteParams{},
				"windows.dhcp.ExclusionRangeV4Delete: exclusion range parameter 'ScopeId', 'StartRange' and 'EndRange' must be a valid IPv4 address",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mockConn := mockConnection.NewMockConnection(suite.T())
			c := &Client{
				Connection:      mockConn,
				decodeCliXmlErr: func(s string) (string, error) { return "", nil },
			}
			err := c.ExclusionRangeV4Delete(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

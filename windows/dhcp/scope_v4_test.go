package dhcp

import (
	"context"
	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"net/netip"
	"time"

	"github.com/d-strobel/gowindows/parsing"
)

// Fixtures
const (
	scopeV4Json = `{"ScopeId":{"Address":698560,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.0"},"SubnetMask":{"Address":16777215,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"255.255.255.0"},"StartRange":{"Address":84584640,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.5"},"EndRange":{"Address":168470720,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.10.10"},"ActivatePolicies":true,"Delay":0,"Description":"Test description","LeaseDuration":{"Ticks":6912000000000,"Days":8,"Hours":0,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":8,"TotalHours":192,"TotalMilliseconds":691200000,"TotalMinutes":11520,"TotalSeconds":691200},"MaxBootpClients":4294967295,"Name":"test","NapEnable":false,"NapProfile":"","State":"Active","SuperscopeName":"","Type":"Dhcp","PSComputerName":null}`
)

var (
	expectedScopeV4 = ScopeV4{
		Name:        "test",
		Description: "Test description",
		ScopeId: addressString{
			Address: netip.MustParseAddr("192.168.10.0"),
		},
		StartRange: addressString{
			Address: netip.MustParseAddr("192.168.10.5"),
		},
		EndRange: addressString{
			netip.MustParseAddr("192.168.10.10"),
		},
		SubnetMask: addressString{
			netip.MustParseAddr("255.255.255.0"),
		},
		State:            "Active",
		MaxBootpClients:  4294967295,
		ActivatePolicies: true,
		NapEnable:        false,
		NapProfile:       "",
		Delay:            0,
		LeaseDuration: parsing.CimTimeDuration{
			Duration: time.Hour * 24 * 8,
		},
	}
)

// Test ScopeV4Read related methods.
func (suite *DhcpServerUnitTestSuite) TestScopeV4ReadPwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4ReadParams
			expectedCmd     string
		}{
			{
				"assert correct command ScopeV4 read by ScopeId",
				ScopeV4ReadParams{ScopeId: netip.MustParseAddr("192.168.10.0")},
				"Get-DhcpServerv4Scope -ScopeId '192.168.10.0' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestScopeV4Read() {
	suite.T().Parallel()

	suite.Run("should return the correct ScopeV4 (Read)", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-DhcpServerv4Scope -ScopeId '192.168.10.0' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: scopeV4Json}, nil)
		actualScopeV4, err := c.ScopeV4Read(ctx, ScopeV4ReadParams{
			ScopeId: netip.MustParseAddr("192.168.10.0"),
		})
		suite.NoError(err)
		suite.Equal(expectedScopeV4, actualScopeV4)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4ReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ScopeV4ReadParams{},
				"windows.dhcp.ScopeV4Read: scope parameter 'ScopeId' must be an IPv4 network address",
			},
			{
				"assert error with empty parameters",
				ScopeV4ReadParams{ScopeId: netip.MustParseAddr("fe80:0010::")},
				"windows.dhcp.ScopeV4Read: scope parameter 'ScopeId' must be an IPv4 network address",
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
			_, err := c.ScopeV4Read(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test ScopeV4Create related methods.
func (suite *DhcpServerUnitTestSuite) TestScopeV4CreatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4CreateParams
			expectedCmd     string
		}{
			{
				"assert correct command with neccessary parameters",
				ScopeV4CreateParams{
					Name:       "test",
					StartRange: netip.MustParseAddr("192.168.10.5"),
					EndRange:   netip.MustParseAddr("192.168.10.10"),
					SubnetMask: netip.MustParseAddr("255.255.255.0"),
				},
				"Add-DhcpServerv4Scope -PassThru -Confirm:$false -Name 'test' -StartRange '192.168.10.5' -EndRange '192.168.10.10' -SubnetMask '255.255.255.0' -State 'InActive' | ConvertTo-Json -Compress",
			},
			{
				"assert correct command with additional parameters",
				ScopeV4CreateParams{
					Name:             "test",
					StartRange:       netip.MustParseAddr("192.168.10.5"),
					EndRange:         netip.MustParseAddr("192.168.10.10"),
					SubnetMask:       netip.MustParseAddr("255.255.255.0"),
					Description:      "Test description",
					Enabled:          true,
					MaxBootpClients:  10000,
					ActivatePolicies: true,
					NapEnable:        true,
					NapProfile:       "testNap",
					Delay:            10,
					LeaseDuration:    time.Duration(27*time.Hour + 30*time.Minute + 10*time.Second),
					Type:             "Dhcp",
					Superscope:       "testSuperscope",
				},
				"Add-DhcpServerv4Scope -PassThru -Confirm:$false -Name 'test' -StartRange '192.168.10.5' -EndRange '192.168.10.10' -SubnetMask '255.255.255.0' -Description 'Test description' -State 'Active' -MaxBootpClients 10000 -ActivatePolicies -NapEnable -NapProfile 'testNap' -Delay 10 -LeaseDuration $(New-TimeSpan -Days 1 -Hours 3 -Minutes 30 -Seconds 10) -Type 'Dhcp' -SuperscopeName 'testSuperscope' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestScopeV4Create() {
	suite.Run("should return the correct ScopeV4 (Create)", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Add-DhcpServerv4Scope -PassThru -Confirm:$false -Name 'test' -StartRange '192.168.10.5' -EndRange '192.168.10.10' -SubnetMask '255.255.255.0' -State 'InActive' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: scopeV4Json}, nil)
		actualScopeV4, err := c.ScopeV4Create(ctx, ScopeV4CreateParams{
			Name:       "test",
			StartRange: netip.MustParseAddr("192.168.10.5"),
			EndRange:   netip.MustParseAddr("192.168.10.10"),
			SubnetMask: netip.MustParseAddr("255.255.255.0"),
		})
		suite.NoError(err)
		suite.Equal(expectedScopeV4, actualScopeV4)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4CreateParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ScopeV4CreateParams{},
				"windows.dhcp.ScopeV4Create: scope parameter 'Name' must be set",
			},
			{
				"assert error without SubnetMask needed fields",
				ScopeV4CreateParams{Name: "test", StartRange: netip.MustParseAddr("192.168.10.5"), EndRange: netip.MustParseAddr("192.168.10.10")},
				"windows.dhcp.ScopeV4Create: scope parameter 'StartRange', 'EndRange' and 'SubnetMask' must be a valid IPv4 address",
			},
			{
				"assert error with IPv6 fields",
				ScopeV4CreateParams{Name: "test", StartRange: netip.MustParseAddr("fe80::0010"), EndRange: netip.MustParseAddr("fe80::0020"), SubnetMask: netip.MustParseAddr("fe80::0005")},
				"windows.dhcp.ScopeV4Create: scope parameter 'StartRange', 'EndRange' and 'SubnetMask' must be a valid IPv4 address",
			},
			{
				"assert error with invalid Type",
				ScopeV4CreateParams{Name: "test", StartRange: netip.MustParseAddr("192.168.10.5"), EndRange: netip.MustParseAddr("192.168.10.10"), SubnetMask: netip.MustParseAddr("255.255.255.0"), Type: "test"},
				"windows.dhcp.ScopeV4Create: scope parameter 'Type' must be one of the following values: 'Dhcp', 'Bootp', 'Both'",
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
			_, err := c.ScopeV4Create(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test ScopeV4Update related methods.
func (suite *DhcpServerUnitTestSuite) TestScopeV4UpdatePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4UpdateParams
			expectedCmd     string
		}{
			{
				"assert correct command with neccessary parameters",
				ScopeV4UpdateParams{
					ScopeId: netip.MustParseAddr("192.168.10.0"),
				},
				"Set-DhcpServerv4Scope -PassThru -Confirm:$false -ScopeId '192.168.10.0' -State 'InActive' | ConvertTo-Json -Compress",
			},
			{
				"assert correct command with additional parameters",
				ScopeV4UpdateParams{
					ScopeId:          netip.MustParseAddr("192.168.10.0"),
					Name:             "test",
					StartRange:       netip.MustParseAddr("192.168.10.5"),
					EndRange:         netip.MustParseAddr("192.168.10.10"),
					Description:      "Test description",
					Enabled:          true,
					MaxBootpClients:  10000,
					ActivatePolicies: true,
					NapEnable:        true,
					NapProfile:       "testNap",
					Delay:            10,
					LeaseDuration:    time.Duration(27*time.Hour + 30*time.Minute + 10*time.Second),
					Type:             "Dhcp",
					Superscope:       "testSuperscope",
				},
				"Set-DhcpServerv4Scope -PassThru -Confirm:$false -ScopeId '192.168.10.0' -Name 'test' -StartRange '192.168.10.5' -EndRange '192.168.10.10' -Description 'Test description' -State 'Active' -MaxBootpClients 10000 -ActivatePolicies -NapEnable -NapProfile 'testNap' -Delay 10 -LeaseDuration $(New-TimeSpan -Days 1 -Hours 3 -Minutes 30 -Seconds 10) -Type 'Dhcp' -SuperscopeName 'testSuperscope' | ConvertTo-Json -Compress",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestScopeV4Update() {
	suite.Run("should return the correct ScopeV4 (Update)", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Set-DhcpServerv4Scope -PassThru -Confirm:$false -ScopeId '192.168.10.0' -Name 'test' -StartRange '192.168.10.5' -EndRange '192.168.10.10' -State 'InActive' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: scopeV4Json}, nil)
		actualScopeV4, err := c.ScopeV4Update(ctx, ScopeV4UpdateParams{
			ScopeId:    netip.MustParseAddr("192.168.10.0"),
			Name:       "test",
			StartRange: netip.MustParseAddr("192.168.10.5"),
			EndRange:   netip.MustParseAddr("192.168.10.10"),
		})
		suite.NoError(err)
		suite.Equal(expectedScopeV4, actualScopeV4)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4UpdateParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ScopeV4UpdateParams{},
				"windows.dhcp.ScopeV4Update: scope parameter 'ScopeId' must be a valid IPv4 address",
			},
			{
				"assert error with StartRange only",
				ScopeV4UpdateParams{ScopeId: netip.MustParseAddr("192.168.10.0"), StartRange: netip.MustParseAddr("192.168.10.10")},
				"windows.dhcp.ScopeV4Update: scope parameter 'StartRange' and 'EndRange' must be set together",
			},
			{
				"assert error with EndRange only",
				ScopeV4UpdateParams{ScopeId: netip.MustParseAddr("192.168.10.0"), EndRange: netip.MustParseAddr("192.168.10.10")},
				"windows.dhcp.ScopeV4Update: scope parameter 'StartRange' and 'EndRange' must be set together",
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
			_, err := c.ScopeV4Update(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

// Test ScopeV4Delete related methods.
func (suite *DhcpServerUnitTestSuite) TestScopeV4DeletePwshCommand() {
	suite.Run("should return the correct command", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4DeleteParams
			expectedCmd     string
		}{
			{
				"assert correct command with neccessary parameters",
				ScopeV4DeleteParams{
					ScopeId: netip.MustParseAddr("192.168.10.0"),
				},
				"Remove-DhcpServerv4Scope -Confirm:$false -ScopeId '192.168.10.0'",
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			actualCmd := tc.inputParameters.pwshCommand()
			suite.Equal(tc.expectedCmd, actualCmd)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestScopeV4Delete() {
	suite.Run("should not error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Remove-DhcpServerv4Scope -Confirm:$false -ScopeId '192.168.10.0'").
			Return(connection.CmdResult{}, nil)
		err := c.ScopeV4Delete(ctx, ScopeV4DeleteParams{
			ScopeId: netip.MustParseAddr("192.168.10.0"),
		})
		suite.NoError(err)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters ScopeV4DeleteParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				ScopeV4DeleteParams{},
				"windows.dhcp.ScopeV4Delete: scope parameter 'ScopeId' must be a valid IPv4 address",
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
			err := c.ScopeV4Delete(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

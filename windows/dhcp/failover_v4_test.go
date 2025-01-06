package dhcp

import (
	"context"
	"encoding/json"
	"net/netip"
	"time"

	"github.com/d-strobel/gowindows/connection"
	mockConnection "github.com/d-strobel/gowindows/connection/mocks"
	"github.com/d-strobel/gowindows/parsing"
)

// Fixtures
const (
	// failover json output
	failoverV4Json = `{"ScopeId":{"value":[{"Address":698560,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false},{"Address":1353920,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false}],"Count":2},"PrimaryServerIP":{"Address":111111,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.5.1"},"SecondaryServerIP":{"Address":222222,"AddressFamily":2,"ScopeId":null,"IsIPv6Multicast":false,"IsIPv6LinkLocal":false,"IsIPv6SiteLocal":false,"IsIPv6Teredo":false,"IsIPv4MappedToIPv6":false,"IPAddressToString":"192.168.5.2"},"AutoStateTransition":false,"EnableAuth":true,"LoadBalancePercent":50,"MaxClientLeadTime":{"Ticks":36000000000,"Days":0,"Hours":1,"Milliseconds":0,"Minutes":0,"Seconds":0,"TotalDays":0.041666666666666664,"TotalHours":1,"TotalMilliseconds":3600000,"TotalMinutes":60,"TotalSeconds":3600},"Mode":"LoadBalance","Name":"dhcp-master\u003c--\u003edhcp-node","PartnerServer":"DHCP-N-P-01.test.local","PrimaryServerName":"DHCP-M-P-01","ReservePercent":null,"SecondaryServerName":"DHCP-N-P-01.test.local","ServerRole":null,"ServerType":"PrimaryServer","State":"Normal","StateSwitchInterval":null,"PSComputerName":null}`
)

var (
	expectedFailoverV4 = FailoverV4{
		ScopeId: scopeIdVal{
			Value: []addressBytes{
				{
					Address: parsing.CimIpAddress{Addr: netip.MustParseAddr("192.168.10.0")},
				},
				{
					Address: parsing.CimIpAddress{Addr: netip.MustParseAddr("192.168.20.0")},
				},
			},
		},
		Name: "dhcp-master<-->dhcp-node",
		PrimaryServerIp: addressString{
			netip.MustParseAddr("192.168.5.1"),
		},
		PrimaryServerName: "DHCP-M-P-01",
		SecondaryServerIp: addressString{
			netip.MustParseAddr("192.168.5.2"),
		},
		SecondaryServerName: "DHCP-N-P-01.test.local",
		AutoStateTransition: false,
		EnableAuth:          true,
		LoadBalancePercent:  50,
		MaxClientLeadTime: parsing.CimTimeDuration{
			Duration: time.Hour,
		},
		Mode:           "LoadBalance",
		ReservePercent: 0,
		ServerRole:     "",
		ServerType:     "PrimaryServer",
		State:          "Normal",
		StateSwitchInterval: parsing.CimTimeDuration{
			Duration: 0,
		},
	}
)

// Test the unmarshalJSON functionality.
func (suite *DhcpServerUnitTestSuite) TestFailoverV4UnmarshalJSON() {
	var f FailoverV4
	err := json.Unmarshal([]byte(failoverV4Json), &f)
	suite.NoError(err)
	suite.Equal(expectedFailoverV4, f)
}

func (suite *DhcpServerUnitTestSuite) TestFailoverV4ReadPwshCommand() {
	inputParameters := FailoverV4ReadParams{
		Name: "test-failover",
	}
	expectedCmd := "Get-DhcpServerv4Failover -Name 'test-failover' | ConvertTo-Json -Compress"
	actualCmd := inputParameters.pwshCommand()
	suite.Equal(expectedCmd, actualCmd)
}

func (suite *DhcpServerUnitTestSuite) TestFailoverV4Read() {
	suite.Run("should return the correct FailoverV4 (Read)", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Get-DhcpServerv4Failover -Name 'dhcp-master<-->dhcp-node' | ConvertTo-Json -Compress").
			Return(connection.CmdResult{StdOut: failoverV4Json}, nil)
		actualFailoverV4, err := c.FailoverV4Read(ctx, FailoverV4ReadParams{
			Name: "dhcp-master<-->dhcp-node",
		})
		suite.NoError(err)
		suite.Equal(expectedFailoverV4, actualFailoverV4)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters FailoverV4ReadParams
			expectedErr     string
		}{
			{
				"assert error with empty parameters",
				FailoverV4ReadParams{},
				"windows.dhcp.FailoverV4Read: failover parameter 'Name' must be set",
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
			_, err := c.FailoverV4Read(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

func (suite *DhcpServerUnitTestSuite) TestFailoverV4CreatePwshCommand() {
	tcs := []struct {
		description     string
		inputParameters FailoverV4CreateParams
		expectedCmd     string
	}{
		{
			"Test with neccessary parameters (PartnerServerName)",
			FailoverV4CreateParams{
				Name: "test-failover",
				ScopeIds: []netip.Addr{
					netip.MustParseAddr("192.168.10.0"),
				},
				PartnerServerName: "dhcp-n-p-01.test.local",
			},
			"Add-DhcpServerv4Failover -PassThru -Confirm:$false -Name 'test-failover' -PartnerServer 'dhcp-n-p-01.test.local' -ScopeId @('192.168.10.0')",
		},
		{
			"Test with neccessary parameters (PartnerServerIp)",
			FailoverV4CreateParams{
				Name: "test-failover",
				ScopeIds: []netip.Addr{
					netip.MustParseAddr("192.168.10.0"),
				},
				PartnerServerIp: netip.MustParseAddr("192.168.5.2"),
			},
			"Add-DhcpServerv4Failover -PassThru -Confirm:$false -Name 'test-failover' -PartnerServer '192.168.5.2' -ScopeId @('192.168.10.0')",
		},
		{
			"Test with additional parameters",
			FailoverV4CreateParams{
				Name: "test-failover",
				ScopeIds: []netip.Addr{
					netip.MustParseAddr("192.168.10.0"),
					netip.MustParseAddr("192.168.20.0"),
				},
				PartnerServerIp:     netip.MustParseAddr("192.168.5.2"),
				LoadBalancePercent:  60,
				MaxClientLeadTime:   time.Hour*20 + time.Minute*20,
				ReservePercent:      20,
				ServerRole:          "PrimaryServer",
				SharedSecret:        "testsharedsecret",
				StateSwitchInterval: time.Hour*4 + time.Minute*30 + time.Second*10,
			},
			"Add-DhcpServerv4Failover -PassThru -Confirm:$false -Name 'test-failover' -PartnerServer '192.168.5.2' -ScopeId @('192.168.10.0','192.168.20.0') -LoadBalancePercent 60 -MaxClientLeadTime $(New-TimeSpan -Days 0 -Hours 20 -Minutes 20 -Seconds 0) -ReservePercent 20 -ServerRole 'PrimaryServer' -SharedSecret 'testsharedsecret' -StateSwitchInterval $(New-TimeSpan -Days 0 -Hours 4 -Minutes 30 -Seconds 10)",
		},
	}

	for _, tc := range tcs {
		suite.T().Logf("test case: %s", tc.description)
		actualCmd := tc.inputParameters.pwshCommand()
		suite.Equal(tc.expectedCmd, actualCmd)
	}
}

func (suite *DhcpServerUnitTestSuite) TestFailoverV4Create() {
	suite.Run("should return the correct FailoverV4 (Create)", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mockConnection.NewMockConnection(suite.T())
		c := &Client{
			Connection:      mockConn,
			decodeCliXmlErr: func(s string) (string, error) { return "", nil },
		}
		mockConn.EXPECT().
			RunWithPowershell(ctx, "Add-DhcpServerv4Failover -PassThru -Confirm:$false -Name 'dhcp-master<-->dhcp-node' -PartnerServer '192.168.5.2' -ScopeId @('192.168.10.0','192.168.20.0')").
			Return(connection.CmdResult{StdOut: failoverV4Json}, nil)
		actualFailoverV4, err := c.FailoverV4Create(ctx, FailoverV4CreateParams{
			Name:            "dhcp-master<-->dhcp-node",
			PartnerServerIp: netip.MustParseAddr("192.168.5.2"),
			ScopeIds: []netip.Addr{
				netip.MustParseAddr("192.168.10.0"),
				netip.MustParseAddr("192.168.20.0"),
			},
		})
		suite.NoError(err)
		suite.Equal(expectedFailoverV4, actualFailoverV4)
	})

	suite.Run("should return specific errors", func() {
		tcs := []struct {
			description     string
			inputParameters FailoverV4CreateParams
			expectedErr     string
		}{
			{
				"assert with empty parameters",
				FailoverV4CreateParams{},
				"windows.dhcp.FailoverV4Create: failover parameters 'Name', 'ScopeIds' and one of 'PartnerServerName', 'PartnerServerIp' must be set",
			},
			{
				"assert with without PartnerServerName",
				FailoverV4CreateParams{
					Name: "test",
					ScopeIds: []netip.Addr{
						netip.MustParseAddr("192.168.10.0"),
					},
				},
				"windows.dhcp.FailoverV4Create: failover parameters 'Name', 'ScopeIds' and one of 'PartnerServerName', 'PartnerServerIp' must be set",
			},
			{
				"assert with without ScopeIds",
				FailoverV4CreateParams{
					Name:              "test",
					PartnerServerName: "test",
				},
				"windows.dhcp.FailoverV4Create: failover parameters 'Name', 'ScopeIds' and one of 'PartnerServerName', 'PartnerServerIp' must be set",
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
			_, err := c.FailoverV4Create(ctx, tc.inputParameters)
			suite.EqualError(err, tc.expectedErr)
		}
	})
}

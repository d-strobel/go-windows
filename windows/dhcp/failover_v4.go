package dhcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/d-strobel/gowindows/parsing"
	"github.com/d-strobel/gowindows/winerror"
)

// FailoverV4 represents an IPv4 failover object.
type FailoverV4 struct {
	Name                string                  `json:"Name"`
	ScopeId             scopeIdVal              `json:"ScopeId"`
	PrimaryServerIp     addressString           `json:"PrimaryServerIP"`
	PrimaryServerName   string                  `json:"PrimaryServerName"`
	SecondaryServerIp   addressString           `json:"SecondaryServerIP"`
	SecondaryServerName string                  `json:"SecondaryServerName"`
	AutoStateTransition bool                    `json:"AutoStateTransition"`
	EnableAuth          bool                    `json:"EnableAuth"`
	LoadBalancePercent  uint32                  `json:"LoadBalancePercent"`
	MaxClientLeadTime   parsing.CimTimeDuration `json:"MaxClientLeadTime"`
	Mode                string                  `json:"Mode"`
	ReservePercent      uint32                  `json:"ReservePercent"`
	ServerRole          string                  `json:"ServerRole"`
	ServerType          string                  `json:"ServerType"`
	State               string                  `json:"State"`
	StateSwitchInterval parsing.CimTimeDuration `json:"StateSwitchInterval"`
}

// FailoverV4ReadParams represents the parameters for the FailoverV4Read function.
type FailoverV4ReadParams struct {
	// Specifies the name of a failover relationship for which the properties are returned.
	Name string
}

// pwshCommand returns the PowerShell command to create an IPv4 Failover.
func (params FailoverV4ReadParams) pwshCommand() string {
	return fmt.Sprintf(
		"Get-DhcpServerv4Failover -Name '%s' | ConvertTo-Json -Compress",
		params.Name,
	)
}

// FailoverV4Read returns a FailoverV4 object.
func (c *Client) FailoverV4Read(ctx context.Context, params FailoverV4ReadParams) (FailoverV4, error) {
	var f FailoverV4

	// Assert needed parameters.
	if params.Name == "" {
		return f, errors.New("windows.dhcp.FailoverV4Read: failover parameter 'Name' must be set")
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &f); err != nil {
		return f, winerror.Errorf(cmd, "windows.dhcp.FailoverV4Read: %s", err)
	}

	return f, nil
}

// FailoverV4CreateParams represents the parameters for the FailoverV4Create function.
type FailoverV4CreateParams struct {
	// Specifies the percentage of DHCP client requests which should be served by the
	// local DHCP server service or the DHCP server service that runs on the computer
	// specified in the ComputerName parameter. The remaining requests would be served
	// by the partner DHCP server service.
	//
	// The default value is 50%.
	LoadBalancePercent uint32

	// Specifies the maximum client lead time for the failover relationship.
	//
	// The default value is 1 hour.
	MaxClientLeadTime time.Duration

	// Specifies the name of the failover relationship to create.
	Name string

	// Specifies the host name, of the partner DHCP server service
	// with which the failover relationship is created.
	PartnerServerName string

	// Specifies the IPv4 address of the partner DHCP server service
	// with which the failover relationship is created.
	PartnerServerIp netip.Addr

	// Specifies the percentage of free IPv4 addresses in the IPv4 address pool of
	// the scope which should be reserved on the standby DHCP server service.
	// In the case of a failover, the IPv4 address from this reserved pool on the
	// standby DHCP server service will be leased to new DHCP clients.
	//
	//The default value is 5.
	ReservePercent uint32

	// Specifies the scope identifiers, in IPv4 address format,
	// which are to be added to the failover relationship.
	ScopeIds []netip.Addr

	// Specifies the role of the local DHCP server service in the hot standby mode.
	//
	// The acceptable values for this parameter are: Active or Standby.
	//
	// The default value is Active for the local DHCP server service,
	// such as the partner DHCP server service that is specified will be a standby DHCP server service.
	ServerRole string

	// Specifies the shared secret to be used for message digest authentication.
	// If not specified, the message digest authentication is turned off.
	SharedSecret string

	// Specifies the time interval for which the DHCP server service operates
	// in the COMMUNICATION INTERRUPTED state before transitioning to the PARTNER DOWN state.
	StateSwitchInterval time.Duration
}

// pwshCommand returns the PowerShell command to create an IPv4 Failover.
func (params FailoverV4CreateParams) pwshCommand() string {
	var scopeList []string

	// Create base command.
	cmd := []string{
		fmt.Sprintf("Add-DhcpServerv4Failover -PassThru -Confirm:$false -Name '%s'",
			params.Name,
		),
	}

	// Add additional required parameters.
	if params.PartnerServerName != "" {
		cmd = append(cmd, fmt.Sprintf("-PartnerServer '%s'", params.PartnerServerName))
	} else if params.PartnerServerIp.Is4() {
		cmd = append(cmd, fmt.Sprintf("-PartnerServer '%s'", params.PartnerServerIp))
	}

	// Add scopeIds with single quotes and join them with commas.
	for _, scope := range params.ScopeIds {
		scopeList = append(scopeList, fmt.Sprintf("'%s'", scope))
	}
	cmd = append(cmd, fmt.Sprintf("-ScopeId @(%s)", strings.Join(scopeList, ",")))

	// Add optional parameters.
	if params.LoadBalancePercent != 0 {
		cmd = append(cmd, fmt.Sprintf("-LoadBalancePercent %d", params.LoadBalancePercent))
	}

	if params.MaxClientLeadTime != 0 {
		cmd = append(cmd, fmt.Sprintf("-MaxClientLeadTime %s", parsing.PwshTimespanString(params.MaxClientLeadTime)))
	}

	if params.ReservePercent != 0 {
		cmd = append(cmd, fmt.Sprintf("-ReservePercent %d", params.ReservePercent))
	}

	if params.ServerRole != "" {
		cmd = append(cmd, fmt.Sprintf("-ServerRole '%s'", params.ServerRole))
	}

	if params.SharedSecret != "" {
		cmd = append(cmd, fmt.Sprintf("-SharedSecret '%s'", params.SharedSecret))
	}

	if params.StateSwitchInterval != 0 {
		cmd = append(cmd, fmt.Sprintf("-StateSwitchInterval %s", parsing.PwshTimespanString(params.StateSwitchInterval)))
	}

	// Return the full command.
	return strings.Join(cmd, " ")
}

// FailoverV4Create creates a new IPv4 failover and returns a FailoverV4 object.
func (c *Client) FailoverV4Create(ctx context.Context, params FailoverV4CreateParams) (FailoverV4, error) {
	var f FailoverV4

	// Assert needed parameters.
	if params.Name == "" || len(params.ScopeIds) == 0 || (params.PartnerServerName == "" && !params.PartnerServerIp.Is4()) {
		return f, errors.New(
			"windows.dhcp.FailoverV4Create: failover parameters 'Name', 'ScopeIds' and one of 'PartnerServerName', 'PartnerServerIp' must be set",
		)
	}

	// Run command
	cmd := params.pwshCommand()
	if err := run(ctx, c, cmd, &f); err != nil {
		return f, winerror.Errorf(cmd, "windows.dhcp.FailoverV4Create: %s", err)
	}

	return f, nil
}

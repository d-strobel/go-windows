// Package dns provides a Go library for handling Windows DNS Server.
// The functions are related to the Powershell dns server cmdlets provided by Windows.
//
// https://learn.microsoft.com/en-us/powershell/module/dnsserver/?view=windowsserver2022-ps
package dns

import (
	"context"
	"encoding/json"
	"time"

	"errors"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parsing"
)

// dns is a type constraint for the run function, ensuring it works with specific types.
type dns interface {
	Zone | []Zone | recordObject | []recordObject
}

// Default Windows DNS TTL.
// https://learn.microsoft.com/en-us/windows/win32/ad/configuration-of-ttl-limits?source=recommendations
var defaultTimeToLive time.Duration = time.Second * 86400

// recordObject contains the unmarshaled json of the powershell record object.
type recordObject struct {
	DistinguishedName string                  `json:"DistinguishedName"`
	Name              string                  `json:"HostName"`
	RecordData        recordRecordData        `json:"RecordData"`
	RecordType        string                  `json:"RecordType"`
	Timestamp         parsing.DotnetTime      `json:"Timestamp"`
	Type              int8                    `json:"Type"`
	TimeToLive        parsing.CimTimeDuration `json:"TimeToLive"`
}
type recordRecordData struct {
	CimInstanceProperties parsing.CimClassKeyVal `json:"CimInstanceProperties"`
}

// Client represents a client for handling DNS server functions.
type Client struct {
	// Connection represents a connection.Connection object.
	Connection connection.Connection

	// decodeCliXmlErr represents a function that decodes a CLIXML error and returns aa  human readable string.
	decodeCliXmlErr func(string) (string, error)
}

// NewClient returns a new instance of the Client.
func NewClient(conn connection.Connection) *Client {
	return NewClientWithParser(conn, parsing.DecodeCliXmlErr)
}

// NewClientWithParser returns a new instance of the Client.
// It requires a connection and parsing as input parameters.
func NewClientWithParser(conn connection.Connection, parsing func(string) (string, error)) *Client {
	return &Client{Connection: conn, decodeCliXmlErr: parsing}
}

// run runs a PowerShell command against a Windows system, handles the command results,
// and unmarshals the output into a local object type.
func run[T dns](ctx context.Context, c *Client, cmd string, d *T) error {
	// Run the command
	result, err := c.Connection.RunWithPowershell(ctx, cmd)
	if err != nil {
		return err
	}

	// Handle stderr
	if result.StdErr != "" {
		stderr, err := c.decodeCliXmlErr(result.StdErr)
		if err != nil {
			return err
		}

		return errors.New(stderr)
	}

	if result.StdOut == "" {
		return nil
	}

	// Unmarshal stdout
	if err = json.Unmarshal([]byte(result.StdOut), &d); err != nil {
		return err
	}

	return nil
}

// Package accounts provides a Go library for handling local Windows accounts.
// The functions are related to the Powershell local accounts cmdlets provided by Windows.
//
// https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.localaccounts/?view=powershell-5.1
package accounts

import (
	"context"
	"encoding/json"

	"errors"
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parsing"
)

// accounts is a type constraint for the run function, ensuring it works with specific types.
type accounts interface {
	Group | []Group | User | []User | GroupMember | []GroupMember
}

// Client represents a client for handling local Windows functions.
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

// SID represents the Security Identifier (SID) of a security principal.
// The Value field contains the actual SID value.
type SID struct {
	Value string `json:"Value"`
}

// run runs a PowerShell command against a Windows system, handles the command results,
// and unmarshals the output into a local object type.
func run[T accounts](ctx context.Context, c *Client, cmd string, l *T) error {
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
	if err = json.Unmarshal([]byte(result.StdOut), &l); err != nil {
		return err
	}

	return nil
}

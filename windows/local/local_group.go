package local

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/d-strobel/gowindows/parser"
	"github.com/d-strobel/gowindows/winerror"
)

type Group struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	SID         SID    `json:"SID"`
}

type GroupParams struct {
	Name        string
	Description string
	SID         string
}

// GroupRead gets a group by a SID or Name and returns a Group object.
func (c *Client) GroupRead(ctx context.Context, params GroupParams) (Group, error) {

	// Declare Group object
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return g, winerror.Errorf(winerror.ConfigError, "GroupRead: group parameter 'Name' or 'SID' must be set")
	}

	// Base command
	cmds := []string{"Get-LocalGroup"}

	// Add parameters
	// Prefer SID over Name
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmd := strings.Join(cmds, " ")

	// Optional parameters
	opts := &parser.PwshOpts{
		JSONOutput: true,
	}

	// Powershell command object
	pwshCmd := parser.NewPwshCommand([]string{cmd}, opts)

	// Run the comand
	result, err := c.Connection.Run(ctx, pwshCmd)
	if err != nil {
		return g, err
	}

	// Handle stderr
	if result.StdErr != "" {
		errXML, err := parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return g, err
		}

		return g, winerror.Errorf(winerror.WindowsError, "GroupRead:\n%s", errXML)
	}

	// Unmarshal result
	err = json.Unmarshal([]byte(result.StdOut), &g)
	if err != nil {
		return g, err
	}

	return g, nil
}

// GroupList returns all groups.
func (c *Client) GroupList(ctx context.Context) ([]Group, error) {

	// Declare slice of Group object
	var g []Group

	// Command
	cmd := "Get-LocalGroup"

	// Optional parameters
	opts := &parser.PwshOpts{
		JSONOutput: true,
	}

	// Powershell command object
	pwshCmd := parser.NewPwshCommand([]string{cmd}, opts)

	// Run the comand
	result, err := c.Connection.Run(ctx, pwshCmd)
	if err != nil {
		return nil, err
	}

	// Handle stderr
	if result.StdErr != "" {
		errXML, err := parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return nil, err
		}

		return nil, winerror.Errorf(winerror.WindowsError, "GroupList:\n%s", errXML)
	}

	// Unmarshal result
	err = json.Unmarshal([]byte(result.StdOut), &g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// GroupCreate creates a new group and returns the Group object.
func (c *Client) GroupCreate(ctx context.Context, params GroupParams) (Group, error) {

	// Declare Group object
	var g Group

	// Assert needed parameters
	if params.Name == "" {
		return g, winerror.Errorf(winerror.ConfigError, "GroupCreate: group parameter 'Name' must be set")
	}

	// Base command
	cmds := []string{"New-LocalGroup"}

	// Add parameters
	cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))

	if params.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Description '%s'", params.Description))
	}

	cmd := strings.Join(cmds, " ")

	// Optional parameters
	opts := &parser.PwshOpts{
		JSONOutput: true,
	}

	// Powershell command object
	pwshCmd := parser.NewPwshCommand([]string{cmd}, opts)

	// Run the comand
	result, err := c.Connection.Run(ctx, pwshCmd)
	if err != nil {
		return g, err
	}

	// Handle stderr
	if result.StdErr != "" {
		errXML, err := parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return g, err
		}

		return g, winerror.Errorf(winerror.WindowsError, "GroupCreate:\n%s", errXML)
	}

	// Unmarshal result
	err = json.Unmarshal([]byte(result.StdOut), &g)
	if err != nil {
		return g, err
	}

	return g, nil
}

// GroupUpdate updates a group and returns the Group object.
// Currently only the description parameter can be changed.
func (c *Client) GroupUpdate(ctx context.Context, params GroupParams) (Group, error) {

	// Declare Group object
	var g Group

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return g, winerror.Errorf(winerror.ConfigError, "GroupUpdate: group parameter 'Name' or 'SID' must be set")
	}

	if params.Description == "" {
		return g, winerror.Errorf(winerror.ConfigError, "GroupUpdate: group parameter 'Description' must be set")
	}

	// Base command
	cmds := []string{"Set-LocalGroup"}

	// Add parameters
	// Prefer SID over Name to identifiy group
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmds = append(cmds, fmt.Sprintf("-Description '%s'", params.Description))

	cmd := strings.Join(cmds, " ")

	// Optional parameters
	opts := &parser.PwshOpts{
		JSONOutput: false,
	}

	// Powershell command object
	pwshCmd := parser.NewPwshCommand([]string{cmd}, opts)

	// Run the comand
	result, err := c.Connection.Run(ctx, pwshCmd)
	if err != nil {
		return g, err
	}

	// Handle stderr
	if result.StdErr != "" {
		errXML, err := parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return g, err
		}

		return g, winerror.Errorf(winerror.WindowsError, "GroupUpdate:\n%s", errXML)
	}

	// Read out group to return the new group object
	group, err := c.GroupRead(ctx, params)
	if err != nil {
		return g, err
	}

	return group, nil
}

// GroupDelete removes a group by a SID or Name.
func (c *Client) GroupDelete(ctx context.Context, params GroupParams) error {

	// Assert needed parameters
	if params.Name == "" && params.SID == "" {
		return winerror.Errorf(winerror.ConfigError, "GroupDelete: group parameter 'Name' or 'SID' must be set")
	}

	// Base command
	cmds := []string{"Remove-LocalGroup"}

	// Add parameters
	// Prefer SID over Name to identifiy group
	if params.SID != "" {
		cmds = append(cmds, fmt.Sprintf("-SID %s", params.SID))
	} else if params.Name != "" {
		cmds = append(cmds, fmt.Sprintf("-Name '%s'", params.Name))
	}

	cmd := strings.Join(cmds, " ")

	// Optional parameters
	opts := &parser.PwshOpts{
		JSONOutput: false,
	}

	// Powershell command object
	pwshCmd := parser.NewPwshCommand([]string{cmd}, opts)

	// Run the comand
	result, err := c.Connection.Run(ctx, pwshCmd)
	if err != nil {
		return err
	}

	// Handle stderr
	if result.StdErr != "" {
		errXML, err := parser.DecodeCLIXML(result.StdErr)
		if err != nil {
			return err
		}

		return winerror.Errorf(winerror.WindowsError, "GroupDelete:\n%s", errXML)
	}

	return nil
}

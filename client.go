package gowindows

import (
	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
	"github.com/d-strobel/gowindows/windows/local"
)

type Client struct {
	Connection *connection.Connection
	parser     *parser.Parser
	Local      *local.LocalClient
}

// NewClient returns a client object that contains the connection and the windows package.
// Use this client to run the functions inside the windows subpackages.
func NewClient(conf *connection.Config) (*Client, error) {

	var err error

	// Init new client
	c := &Client{}

	// Store new connection to the client
	c.Connection, err = connection.NewConnection(conf)
	if err != nil {
		return nil, err
	}

	// Store parser to the client
	c.parser = parser.NewParser()

	// Build the client with the subpackages
	c.Local = local.NewLocalClient(c.Connection, c.parser)

	return c, nil
}

// Close closes any open connection.
// For now only ssh connections will be terminated.
// To avoid surprises in the future, this should always be called in a defer statement.
func (c *Client) Close() error {
	return c.Connection.Close()
}

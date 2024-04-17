package connection_test

import (
	"testing"

	"github.com/d-strobel/gowindows/connection"
	"github.com/stretchr/testify/suite"
)

const (
	sshHost     = "127.0.0.1"
	sshPort     = 1222
	sshUsername = "vagrant"
	sshPassword = "vagrant"

	// ED25519 key
	sshKeyPathED25519 = "./fixtures/ed25519"
	sshKeyED25519     = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCNKqLkDaaa4KGp+xaT0X94XVxGiwG6RHsymEc9/m39hwAAAJjpeDkr6Xg5
KwAAAAtzc2gtZWQyNTUxOQAAACCNKqLkDaaa4KGp+xaT0X94XVxGiwG6RHsymEc9/m39hw
AAAEAMT15+Ut2N+m9HW9wXgIeVR+qKeoT3UlVCxxnPsnoA5o0qouQNpprgoan7FpPRf3hd
XEaLAbpEezKYRz3+bf2HAAAAD2RzdHJvYmVsQE5CMDc4NAECAwQFBg==
-----END OPENSSH PRIVATE KEY-----`

	// RSA key
	sshKeyPathRSA = "./fixtures/id_rsa"
	sshKeyRSA     = `-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`
)

// Init acceptance test suite for SSH
type ConnectionSSHAccTestSuite struct {
	suite.Suite
}

func TestConnectionSSHAccTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, &ConnectionSSHAccTestSuite{})
}

func (suite *ConnectionSSHAccTestSuite) TestNewConnectionWithSSH() {
	suite.T().Parallel()

	suite.Run("should establish a connection via password", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:     sshHost,
			SSHPort:     sshPort,
			SSHUsername: sshUsername,
			SSHPassword: sshPassword,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey path with ed25519", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:           sshHost,
			SSHPort:           sshPort,
			SSHUsername:       sshUsername,
			SSHPrivateKeyPath: sshKeyPathED25519,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey path with rsa", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:           sshHost,
			SSHPort:           sshPort,
			SSHUsername:       sshUsername,
			SSHPrivateKeyPath: sshKeyPathRSA,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey with ed25519", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:       sshHost,
			SSHPort:       sshPort,
			SSHUsername:   sshUsername,
			SSHPrivateKey: sshKeyED25519,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey with rsa", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:       sshHost,
			SSHPort:       sshPort,
			SSHUsername:   sshUsername,
			SSHPrivateKey: sshKeyRSA,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})
}
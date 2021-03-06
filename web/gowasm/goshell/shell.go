package goshell

import (
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

const maxSSHRetries = 10
const maxSSHDelay = 50 * time.Microsecond

type SecureShell struct {
	client *ssh.Client
	output io.Writer
}

func NewSecureShell(output io.Writer, host, username, password string, port ...int) (*SecureShell, error) {
	sshPort := 22
	if port != nil {
		sshPort = port[0]
	}

	keyboardInteractiveChallenge := func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil
		}
		return []string{password}, nil
	}

	// Retry few times if ssh connection fails
	for i := 0; i < maxSSHRetries; i++ {
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, sshPort), &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
				ssh.KeyboardInteractive(keyboardInteractiveChallenge),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			time.Sleep(maxSSHDelay)
			errMessage := fmt.Sprintf("Failed to dial host: %+v\n", err)
			output.Write([]byte(errMessage))
			fmt.Println(errMessage)
			continue
		}
		s, err := client.NewSession()
		if err != nil {
			client.Close()
			time.Sleep(maxSSHDelay)
			continue
		}
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		// Request pseudo terminal
		if err := s.RequestPty("xterm", 40, 80, modes); err != nil {
			return nil, fmt.Errorf("failed to get pseudo-terminal: %v", err)
		}
		output.Write([]byte(fmt.Sprintf("---Connected to %s successfully.---\n", host)))
		return &SecureShell{client: client, output: output}, nil
	}
	output.Write([]byte(fmt.Sprintf("---Connected to %s. unsuccessfully---\n", host)))
	return nil, fmt.Errorf("retry times was exceeded 10")
}

func (s *SecureShell) ExecuteCommand(cmd string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	w := io.MultiWriter(s.output)
	session.Stdout = w
	session.Stderr = w
	if err := session.Start(cmd); err != nil {
		return err
	}
	if err := session.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *SecureShell) Output(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	b, err := session.Output(cmd)
	return string(b), err
}

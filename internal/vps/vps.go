package vps

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// VPSConfig struct
type VPSConfig struct {
	ServerName    string `json:"server_name"`
	ServerUser    string `json:"server_user"`
	ServerPasswd  string `json:"server_passwd"`
	ServerIP      string `json:"server_ip"`
	ServerPort    int    `json:"server_port"`
	ADSLAccount   string `json:"adsl_account"`
	ADSLPasswd    string `json:"adsl_passwd"`
	VPSOpenTime   string `json:"vps_opentime"`
	VPSExpireTime string `json:"vps_expiretime"`
}

// Vps struct
type Vps struct {
	conf *ssh.ClientConfig
	Addr string

	stdout  *bytes.Buffer
	session *ssh.Session
	client  *ssh.Client
}

// OpenVPS func
func OpenVPS(conf *VPSConfig) *Vps {
	v := &Vps{
		conf: &ssh.ClientConfig{
			User:            conf.ServerUser,
			Auth:            []ssh.AuthMethod{ssh.Password(conf.ServerPasswd)},
			Timeout:         30 * time.Second,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		Addr:   fmt.Sprintf("%s:%d", conf.ServerIP, conf.ServerPort),
		stdout: bytes.NewBuffer([]byte("")),
	}
	return v
}

// Run func
func (v *Vps) Run(cmd string) (string, error) {
	defer func() {
		if v.client != nil {
			v.client.Close()
		}
	}()
	session, err := v.dial()
	defer func() {
		if session != nil {
			session.Close()
		}
	}()
	if err != nil {
		return "", err
	}
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(out), fmt.Errorf("%v, %s", err, string(out))
	}
	return string(out), nil
}

func (v *Vps) dial() (*ssh.Session, error) {
	client, err := ssh.Dial("tcp", v.Addr, v.conf)
	v.client = client
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	v.session = session
	if err != nil {
		return nil, err
	}
	return session, nil
}

// TerminalRun fn
func (v *Vps) TerminalRun() error {
	session, err := v.dial()
	if err != nil {
		return err
	}
	defer session.Close()
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		return err
	}
	return session.Run("/bin/bash")
}

// Close func
func (v *Vps) Close() {
	if v.session != nil {
		v.session.Close()
	}
	if v.client != nil {
		v.client.Close()
	}
}

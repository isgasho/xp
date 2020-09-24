// https://www.cnblogs.com/chenqionghe/p/8267326.html

package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	privateKey = fmt.Sprintf("%s/.ssh/id_rsa", home)
}

var privateKey string

type Cli struct {
	IP         string      //IP地址
	Username   string      //用户名
	Password   string      //密码
	Port       int         //端口号
	client     *ssh.Client //ssh客户端
	LastResult string      //最近一次Run的结果
}

//创建命令行对象
//@param ip IP地址
//@param username 用户名
//@param password 密码
//@param port 端口号,默认22
func New(ip string, username string, password string, port ...int) *Cli {
	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}
	return cli
}

//执行shell
//@param shell shell脚本命令
func (c Cli) Run(shell string) (string, error) {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return "", err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err
}

//连接
func (c *Cli) connect() error {
	var config ssh.ClientConfig
	if c.Password != "" {
		config = ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
		}
	} else {
		key, err := ioutil.ReadFile(privateKey)
		if err != nil {
			fmt.Printf("Unable to read private key: %v\n", err)
			return err
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			fmt.Printf("Unable to parse private key: %v\n", err)
			return err
		}

		config = ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
		}
	}

	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}

//执行带交互的命令
func (c *Cli) RunTerminal(shell string, stdout, stderr io.Writer) error {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}

	session.Run(shell)
	return nil
}

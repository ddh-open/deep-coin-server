package hostShell

import (
	"bytes"
	"devops-http/app/module/base"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"time"
)

type Context struct {
	UUID uuid.UUID
	Ip   string
	User string
	Port int
	Logs chan string

	SSHBuffer  *SShBuffer
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
	Start      bool
	ExecStart  chan struct{}
	ExecEnd    chan struct{}
	Cancel     chan struct{}
}

type SShBuffer struct {
	outBuf   *bytes.Buffer
	stdinBuf io.WriteCloser
}

func NewContext(ip string, port int, user string) *Context {
	var stdinBuf io.WriteCloser
	return &Context{
		Ip:        ip,
		Port:      port,
		User:      user,
		Logs:      make(chan string, 20),
		ExecStart: make(chan struct{}, 1),
		ExecEnd:   make(chan struct{}, 1),
		Cancel:    make(chan struct{}, 1),
		SSHBuffer: &SShBuffer{
			bytes.NewBuffer(make([]byte, 0)),
			stdinBuf,
		},
	}
}

func (c *Context) InitTerminalWithPassword(password string) error {
	if c.Start {
		base.Logger.Error("session is start terminal")
	}
	err := c.InitSession(password)
	if err != nil {
		return err
	}
	session := c.SSHSession
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("xterm", 400, 300, modes); err != nil {
		base.Logger.Error("get pty error", zap.Error(err))
		return err
	}
	stdinBuf, err := session.StdinPipe()
	if err != nil {
		base.Logger.Error("get stdin pipe error", zap.Error(err))
		return err
	}
	c.SSHBuffer.stdinBuf = stdinBuf
	session.Stdout = c.SSHBuffer.outBuf

	err = session.Shell()
	if err != nil {
		base.Logger.Error("shell session error :", zap.Error(err))
		return err
	}
	c.ExecStart <- struct{}{}
	go c.listenMessages(true)
	<-c.ExecEnd
	go session.Wait()
	return err
}

func (c *Context) InitSession(password string) error {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 3,
		User:            c.User,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", c.Ip, c.Port), config)
	if err != nil {
		base.Logger.Error("dial SSH error :", zap.Error(err))
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		base.Logger.Error("get session error :", zap.Error(err))
		return err
	}
	c.SSHSession = session
	c.User = config.User
	return nil
}

func (c *Context) InitClient() error {
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("dou.190824")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", "45.136.184.165:22", config)
	if err != nil {
		base.Logger.Error("dial SSH error :", zap.Error(err))
		return err
	}
	c.SSHClient = client
	<-c.ExecEnd
	c.User = config.User
	return nil
}

func (c *Context) listenMessages(wait bool) error {
	buf := make([]byte, 8192)
	var t int
	terminator := []byte{'$', ' '}
	if c.User == "root" {
		terminator = []byte{'#', ' '}
	}
out:
	for {
		time.Sleep(200 * time.Millisecond)
		select {
		case <-c.Cancel:
			// 开始的可以进行下一个任务了
			<-c.ExecStart
			// 执行的结束的反馈
			if wait {
				c.ExecEnd <- struct{}{}
			}
			break out
		default:
			n, err := c.SSHBuffer.outBuf.Read(buf)
			if err != nil && err != io.EOF {
				base.Logger.Error("read out buffer err:%v", zap.Error(err))
				break out
			}
			if n > 0 {
				c.Logs <- string(buf[:n])
				t = bytes.LastIndex(buf, terminator)
				if t > 0 {
					// 开始的可以进行下一个任务了
					<-c.ExecStart
					// 执行的结束的反馈
					if wait {
						c.ExecEnd <- struct{}{}
					}
					break out
				}
			}
		}
	}
	return nil
}

func (c *Context) EnableSudo() error {
	err := c.SendCmd("su")
	err = c.SendCmd("root")
	c.User = "root"
	return err
}

func (c *Context) SendCmd(cmd string) error {
	c.ExecStart <- struct{}{}
	c.Logs <- "  " + cmd + "\n"
	_, err := c.SSHBuffer.stdinBuf.Write([]byte(fmt.Sprintf("%v\n", cmd)))
	if err != nil {
		<-c.ExecStart
		return err
	}
	go c.listenMessages(true)
	<-c.ExecEnd
	return err
}

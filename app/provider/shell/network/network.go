package network

import (
	"bytes"
	"context"
	"devops-http/app/module/base"
	"encoding/binary"
	"errors"
	"go.uber.org/zap"
	"hash/crc32"
	"io"
	"math"
	"net"
	"sync"
	"time"
)

var errTooLong = errors.New("too long")
var errChecksum = errors.New("invalid checksum")
var errTimeout = errors.New("timeout")

// Conn network connection
type Conn struct {
	c        net.Conn
	lockRead sync.Mutex
	sizeRead [6]byte
	chWrite  chan []byte
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewConn create connection
func NewConn(c net.Conn) *Conn {
	ctx, cancel := context.WithCancel(context.Background())
	conn := &Conn{
		c:       c,
		chWrite: make(chan []byte, 1024),
		ctx:     ctx,
		cancel:  cancel,
	}
	go conn.loopWrite()
	return conn
}

// Close connection
func (c *Conn) Close() {
	c.c.Close()
	c.cancel()
}

func (c *Conn) read(timeout time.Duration) (uint32, uint16, []byte, error) {
	c.lockRead.Lock()
	defer c.lockRead.Unlock()
	c.c.SetReadDeadline(time.Now().Add(timeout))
	_, err := io.ReadFull(c.c, c.sizeRead[:])
	if err != nil {
		return 0, 0, nil, err
	}
	size := binary.BigEndian.Uint16(c.sizeRead[:])
	enc := binary.BigEndian.Uint32(c.sizeRead[2:])
	buf := make([]byte, size)
	_, err = io.ReadFull(c.c, buf)
	if err != nil {
		return 0, 0, nil, err
	}
	return enc, size, buf, nil
}

// ReadMessage read message with timeout
func (c *Conn) ReadMessage(timeout time.Duration) ([]byte, uint16, error) {
	enc, size, buf, err := c.read(timeout)
	if err != nil {
		return nil, 0, err
	}
	if crc32.ChecksumIEEE(buf) != enc {
		return nil, 0, errChecksum
	}
	return buf, size, nil
}

// WriteMessage write message with timeout
func (c *Conn) WriteMessage(data []byte, timeout time.Duration) error {
	if len(data) > math.MaxUint16 {
		return errTooLong
	}
	buf := make([]byte, len(data)+len(c.sizeRead))
	binary.BigEndian.PutUint16(buf, uint16(len(data)))
	binary.BigEndian.PutUint32(buf[2:], crc32.ChecksumIEEE(data))
	copy(buf[len(c.sizeRead):], data)
	select {
	case c.chWrite <- buf:
		return nil
	case <-time.After(timeout):
		return errTimeout
	}
}

// RemoteAddr get connection remote address
func (c *Conn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

// LocalAddr get connection local address
func (c *Conn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

func (c *Conn) loopWrite() {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return
		case data := <-c.chWrite:
			_, err := io.Copy(c.c, bytes.NewReader(data))
			if err != nil {
				base.Logger.Error("write data: %v", zap.Error(err))
				return
			}
		}
	}
}

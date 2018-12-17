package viaproxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

// Wrap takes a net.Conn and returns a pointer to Conn that knows how to
// properly identify the remote address if it comes via a proxy that
// supports the Proxy Protocol.
func Wrap(cn net.Conn) (*Conn, error) {
	c := &Conn{cn: cn, r: bufio.NewReader(cn)}
	if err := c.parse(); err != nil {
		return nil, err
	}
	return c, nil
}

// Conn is an implementation of net.Conn interface for TCP connections that come
// from a proxy that users the Proxy Protocol to communicate with the upstream
// servers.
type Conn struct {
	cn     net.Conn
	r      *bufio.Reader
	proxy  net.Addr
	remote net.Addr
}

// ProxyAddr returns the proxy remote network address.
func (c *Conn) ProxyAddr() net.Addr { return c.proxy }

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	if c.remote != nil {
		return c.remote
	}
	return c.cn.RemoteAddr()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr { return c.cn.LocalAddr() }

// Read reads data from the connection.
func (c *Conn) Read(b []byte) (int, error) { return c.r.Read(b) }

// Close closes the connection.
func (c *Conn) Close() error { return c.cn.Close() }

// SetDeadline implements the Conn SetDeadline method.
func (c *Conn) SetDeadline(t time.Time) error { return c.cn.SetDeadline(t) }

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *Conn) SetReadDeadline(t time.Time) error { return c.cn.SetReadDeadline(t) }

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c *Conn) SetWriteDeadline(t time.Time) error { return c.cn.SetWriteDeadline(t) }

// Write implements the Conn Write method.
func (c *Conn) Write(b []byte) (int, error) { return c.cn.Write(b) }

var (
	unknown = []byte("UNKNOWN\r\n")
	proxy   = []byte("PROXY ")
)

func (c *Conn) parse() error {
	buf, err := c.r.Peek(len(proxy))
	if err != nil {
		if err == bufio.ErrBufferFull || err == io.EOF {
			return nil
		}

		return fmt.Errorf("parsing proxy protocol header on loop: %q", err)
	}

	if !bytes.Equal(buf, proxy) {
		return nil
	}

	err = c.init()
	return err
}

func (c *Conn) init() error {
	// PROXY
	buf := make([]byte, 6)
	n, err := c.r.Read(buf)
	if err != nil {
		return err
	}
	if !bytes.Equal(buf, []byte("PROXY ")) {
		return fmt.Errorf("invalid proxy protocol header prefix: %q", buf[:n])
	}

	buf, err = c.r.Peek(len(unknown))
	if err != nil {
		return fmt.Errorf("parsing proxy protocol header: %q", err)
	}
	if bytes.Equal(buf, unknown) {
		_, err = c.r.Discard(len(unknown))
		return err
	}

	// TCP4 || TCP6
	buf = make([]byte, 5)
	// This line cannot return error as the buffer of the *bufio.Reader already contains at least five characters from the call to Peek above.
	c.r.Read(buf)
	if !bytes.Equal([]byte("TCP4 "), buf) && !bytes.Equal([]byte("TCP6 "), buf) {
		return fmt.Errorf("unrecognized protocol: %q", buf)
	}

	// CLIENT IP
	clientIP, err := c.readIP()
	if err != nil {
		return fmt.Errorf("cannot parse client IP: %q", err)
	}

	// PROXY IP
	proxyIP, err := c.readIP()
	if err != nil {
		return fmt.Errorf("cannot parse proxy IP: %q", err)
	}

	// CLIENT PORT
	clientPort, err := c.readPort(' ')
	if err != nil {
		return fmt.Errorf("cannot parse client port: %q", err)
	}

	// PROXY PORT
	proxyPort, err := c.readPort('\r')
	if err != nil {
		return fmt.Errorf("cannot parse proxy port: %q", err)
	}

	// Trailing
	b, err := c.r.ReadByte()
	if err != nil || b != '\n' {
		return fmt.Errorf("invalid trailing: %q", err)
	}

	c.remote = &net.TCPAddr{IP: clientIP, Port: clientPort}
	c.proxy = &net.TCPAddr{IP: proxyIP, Port: proxyPort}

	return nil
}

func (c *Conn) readIP() (net.IP, error) {
	p, err := c.r.ReadString(' ')
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(p[:len(p)-1])
	if ip == nil {
		return nil, fmt.Errorf("cannot parse IP %q", p)
	}

	return ip, nil
}

func (c *Conn) readPort(delim byte) (int, error) {
	p, err := c.r.ReadString(delim)
	if err != nil {
		return 0, err
	}

	port, err := strconv.Atoi(p[:len(p)-1])
	if err != nil {
		return 0, err
	}

	return port, nil
}

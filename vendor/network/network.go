package network

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	Server = iota
	Client
)

type Connection struct {
	Type int

	Conn *net.UDPConn

	Handlers map[string]Handler

	Clients []*net.UDPAddr
}

type Handler func(*Request) interface{}

//udpAddr return resolved addr for udp connection
func udpAddr(addr string) *net.UDPAddr {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln("failed addr", err)
	}
	return raddr
}

func NewHandlers(h map[string]Handler) *Connection {
	c := &Connection{
		Handlers: h,
	}
	return c
}

func (c *Connection) SetHandler(name string, h Handler) {
	c.Handlers[name] = h
}

//Server - start server listener
func (c *Connection) Server(addr string) error {
	conn, err := net.ListenUDP("udp", udpAddr(addr))
	if err != nil {
		return err
	}

	c.Type = Server
	c.Conn = conn

	go c.carrier()

	return nil
}

//Client - initialize connection
func (c *Connection) Client(addr string) error {
	conn, err := net.DialUDP("udp", udpAddr(":0"), udpAddr(addr))
	if err != nil {
		return err
	}

	c.Type = Client
	c.Conn = conn

	go c.carrier()

	return nil
}

func (c *Connection) Close() {
	c.Conn.Close()
	c.Clients = nil
	c.Handlers = nil
}

func (c *Connection) DeleteClient(addr *net.UDPAddr) {
	for i, a := range c.Clients {
		if a == addr {
			c.DeleteClientN(i)
		}
	}
}

func (c *Connection) DeleteClientN(i int) {
	c.Clients[i] = c.Clients[len(c.Clients)-1]
	c.Clients[len(c.Clients)-1] = nil
	c.Clients = c.Clients[:len(c.Clients)-1]
}

func (c *Connection) AddClient(addr *net.UDPAddr) {
	for _, a := range c.Clients {
		if a.String() == addr.String() {
			return
		}
	}

	c.Clients = append(c.Clients, addr)
}

//carrier manage incoming messages
func (c *Connection) carrier() {
	for {
		var b = make([]byte, 8192)
		i, addr, err := c.Conn.ReadFromUDP(b)
		if err != nil {
			c.DeleteClient(addr)
			log.Println(c.Type, "failed read udp package, error: ", err)
			continue
		}

		b = b[:i]
		if len(b) == 0 {
			continue
		}

		c.AddClient(addr)
		// c.Clients = append(c.Clients, addr)

		m, err := c.decodeMessage(b)
		if err != nil {
			log.Println(c.Type, "failed decode message", err)
			continue
		}

		// log.Println(addr)

		req := &Request{
			Conn:       c.Conn,
			RemoteAddr: addr,
			Message:    m,
		}

		responseMsg, err := c.callHandler(req)
		if err != nil {
			log.Println(c.Type, "failed send response:", err)
			continue
		}
		if responseMsg != nil {
			_, err := c.Conn.WriteToUDP(responseMsg, addr)
			if err != nil {
				log.Println("failed write to udp channel:", err)
			}
		}

		// if data != nil {
		// 	if err := req.SendResponse(data); err != nil {
		// 		log.Println("failed send response:", err)
		// 		continue
		// 	}
		// }

		// log.Println("get new message, call function:", m.ResponseFunc)
		// m.remoteAddr = c.Conn.RemoteAddr()

		// if len(m.RequestFunc) > 0 {
		// 	if f, ok := c.Handlers[m.RequestFunc]; ok {
		// 		data := m.Response(f(m.Data))
		// 		if data != nil {
		// 			c.Conn.WriteToUDP(data, addr)
		// 		}
		// 	} else {
		// 		log.Printf("called function `%s` is not found\n")
		// 	}
		// }
		// c.Handler(m)
	}
	log.Println("wtf? exit?")
}

func (c *Connection) decodeMessage(b []byte) (m *Message, err error) {
	buf := bytes.NewBuffer(b)
	err = gob.NewDecoder(buf).Decode(&m)
	return
}

func (c *Connection) callHandler(req *Request) ([]byte, error) {
	if handler, ok := c.Handlers[req.RequestFunc]; ok {
		data := handler(req)
		if data != nil {
			return req.NewResponse(data)
		}
		return nil, nil
	}

	log.Println(c.Handlers)

	return nil, fmt.Errorf("handler `%s` not found", req.RequestFunc)
}

type Request struct {
	Conn       *net.UDPConn
	RemoteAddr *net.UDPAddr
	*Message
}

type Message struct {
	Sendtime     time.Time
	RequestFunc  string
	ResponseFunc string

	Data interface{}
}

func NewMessage(reqfunc, resfunc string, data interface{}) *Message {
	m := &Message{
		Sendtime:     time.Now(),
		RequestFunc:  reqfunc,
		ResponseFunc: resfunc,
		Data:         data,
	}

	return m
}

func (m *Message) encode() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(buf).Encode(m)
	return append(buf.Bytes(), '\n'), err
}

func (c *Connection) Broadcast(reqfunc, resfunc string, data interface{}) error {

	bMsg, err := NewMessage(reqfunc, resfunc, data).encode()
	if err != nil {
		return err
	}

	log.Println("send broadcast to ", c.Clients)
	for i, client := range c.Clients {
		log.Println("send to", client.String())
		_, err := c.Conn.WriteToUDP(bMsg, client)
		if err != nil {
			c.DeleteClientN(i)
			log.Printf("failed send broadcast message to `%s:%s`, reason: %s\n", client, err)
		}
	}

	return nil
}

//SendMessage to server from clients
func (c *Connection) SendMessage(reqfunc, resfunc string, data interface{}) error {
	if c.Conn == nil {
		return errors.New("no connection")
	}

	bMsg, err := NewMessage(reqfunc, resfunc, data).encode()
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(bMsg)
	return err
}

func (req *Request) NewResponse(data interface{}) ([]byte, error) {
	if req.ResponseFunc == "" {
		return nil, errors.New("response function is empty")
	}

	message := NewMessage(req.ResponseFunc, "", data)

	return message.encode()
}

// func (req *Request) SendResponse(data interface{}) error {
// 	if req.ResponseFunc == "" {
// 		return errors.New("response function is empty")
// 	}
// 	if req.Conn == nil {
// 		return errors.New("no connection")
// 	}

// 	bMsg, err := NewMessage(req.ResponseFunc, "", data).encode()
// 	if err != nil {
// 		return err
// 	}
// 	// newReq := &Request{
// 	// 	Sendtime:    time.Now(),
// 	// 	RequestFunc: req.ResponseFunc,
// 	// 	Data:        data,
// 	// }

// 	// // Connection.Conn.WriteToUDP([]byte{}, m.remoteAddr)
// 	// log.Println("send response to:", req.ResponseFunc)
// 	// var buf bytes.Buffer
// 	// err := gob.NewEncoder(&buf).Encode(newReq)
// 	// if err != nil {
// 	// 	return err
// 	// 	// log.Println("failed send response message, error:", err)
// 	// }

// 	_, err = req.Conn.Write(bMsg)
// 	return err
// 	// return buf.Bytes()
// 	// return nil
// }

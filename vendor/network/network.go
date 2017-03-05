package network

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	server = iota
	client
)

const disconnect = 0x04

//Connection type
type Connection struct {
	Type int

	Conn *net.UDPConn

	// HandlersServer reflect.Value
	// HandlersClient reflect.Value

	Handlers map[string]Handler

	Clients map[string]*net.UDPAddr
}

//Handler type of function
type Handler func(*Request) interface{}

//NewHandlers create connection width handlers from map
// func NewHandlers(h map[string]Handler) *Connection {
// 	c := &Connection{
// 		Handlers: h,
// 	}
// 	return c
// }

//NewConnection register handlers to connection from type, example: type Handlers struct{} and then func (Handlers) Funcname...
func NewConnection(h interface{}) *Connection {
	var handlers = make(map[string]Handler)

	// return &Connection{
	// 	HandlersServer: reflect.ValueOf(server),
	// 	HandlersClient: reflect.ValueOf(client),
	// }

	v := reflect.ValueOf(h)
	t := reflect.TypeOf(h)
	for i := 0; i < t.NumMethod(); i++ {

		f, ok := v.Method(i).Interface().(func(*Request) interface{})
		if ok {
			handlers[t.Method(i).Name] = f
		} else {
			log.Fatalf("handler %s is not suitable", t.Method(i).Name)
		}

	}

	return &Connection{
		Handlers: handlers,
		Clients:  make(map[string]*net.UDPAddr),
	}
}

// //SetHandler to exist connection by name
// func (c *Connection) SetHandler(name string, h Handler) {
// 	c.Handlers[name] = h
// }

//Server - start server listener
func (c *Connection) Server(addr string) error {
	conn, err := net.ListenUDP("udp", udpAddr(addr))
	if err != nil {
		return err
	}

	c.Type = server
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

	c.Type = client
	c.Conn = conn

	go c.carrier()

	return nil
}

func (c *Connection) Close() {
	c.Conn.Write([]byte{disconnect})
	// c.Conn.WriteToUDP([]byte{disconnect, '\n'}, client)
	// }
	c.Conn.Close()
	c.Clients = nil
	c = nil
}

func (c *Connection) DeleteClient(addr *net.UDPAddr) {
	// log.Println("DeleteClient", addr, c.Clients)
	delete(c.Clients, addr.String())
	// log.Println(c.Clients)
	// for i, a := range c.Clients {
	// 	if a == addr {
	// 		c.DeleteClientN(i)
	// 	}
	// }
}

// func (c *Connection) DeleteClientN(i int) {
// 	c.Clients[i] = c.Clients[len(c.Clients)-1]
// 	c.Clients[len(c.Clients)-1] = nil
// 	c.Clients = c.Clients[:len(c.Clients)-1]
// }

func (c *Connection) AddClient(addr *net.UDPAddr) {
	c.Clients[addr.String()] = addr
	// for _, a := range c.Clients {
	// 	if a.String() == addr.String() {
	// 		return
	// 	}
	// }

	// c.Clients = append(c.Clients, addr)
}

//carrier manage incoming messages
func (c *Connection) carrier() {
	for {
		if c == nil {
			break
		}

		var b = make([]byte, 8192)
		i, addr, err := c.Conn.ReadFromUDP(b)
		if err != nil {

			// c.DeleteClient(addr)
			// if c.Type == client {
			// 	c.Close()
			// }
			log.Println(c.Type, "failed read udp package, error: ", err)
			c.Close()
			break
		}

		b = b[:i]
		if len(b) == 0 {
			continue
		}

		if i == 1 && b[0] == disconnect {
			log.Println("DISCONNECT", addr)

			c.DeleteClient(addr)

			if c.Type == client {
				c.Close()
			}
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
	// log.Println(c.Type, "wtf? exit?")
}

func (c *Connection) decodeMessage(b []byte) (m *Message, err error) {
	buf := bytes.NewBuffer(b)
	err = gob.NewDecoder(buf).Decode(&m)
	return
}

func (c *Connection) callHandler(req *Request) ([]byte, error) {
	f, ok := c.Handlers[req.RequestFunc]
	if !ok {
		log.Println("request unknown function", req.RequestFunc)
		return nil, nil
	}

	data := f(req)
	if data != nil {
		return req.NewResponse(data)
	}
	return nil, nil

	// if handler, ok := c.Handlers[req.RequestFunc]; ok {
	// 	data := handler(req)
	// 	if data != nil {
	// 		return req.NewResponse(data)
	// 	}
	// 	return nil, nil
	// }

	// log.Println(c.Handlers)

	// return nil, fmt.Errorf("handler `%s` not found", req.RequestFunc)
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

func newMessage(reqfunc, resfunc string, data interface{}) *Message {

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

func (c *Connection) Broadcast(reqfunc, resfunc interface{}, data interface{}) error {

	reqname, resname := getFuncsName(reqfunc, resfunc)

	bMsg, err := newMessage(reqname, resname, data).encode()
	if err != nil {
		return err
	}

	// log.Println("send broadcast to ", c.Clients)
	for _, client := range c.Clients {
		// log.Println("send to", client.String())
		_, err := c.Conn.WriteToUDP(bMsg, client)
		if err != nil {
			c.DeleteClient(client)
			log.Printf("failed send broadcast message to `%s:%s`, reason: %s\n", client, err)
		}
	}

	return nil
}

func getFuncsName(reqfunc, resfunc interface{}) (string, string) {
	reqName := getFuncName(reqfunc)

	var resName string
	if resfunc != nil {
		resName = getFuncName(resfunc)
	}

	return reqName, resName
}

func getFuncName(f interface{}) string {
	if _, ok := f.(string); ok {
		log.Fatalln("incoming function is string not Handler")
	}
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name = strings.TrimPrefix(filepath.Ext(name), ".")
	return strings.TrimSuffix(name, "-fm")
}

//SendMessage to server from clients
func (c *Connection) SendMessage(reqfunc, resfunc interface{}, data interface{}) error {
	if c.Conn == nil {
		return errors.New("no connection")
	}

	reqname, resname := getFuncsName(reqfunc, resfunc)

	bMsg, err := newMessage(reqname, resname, data).encode()
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

	message := newMessage(req.ResponseFunc, "", data)

	return message.encode()
}

//udpAddr return resolved addr for udp connection
func udpAddr(addr string) *net.UDPAddr {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln("failed addr", err)
	}
	return raddr
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

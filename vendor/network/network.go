package network

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"
)

var (
	Connection = &connection{
		Routes: make(map[string]func(interface{}) interface{}),
	}

	err error
)

type connection struct {
	Conn       *net.UDPConn
	remoteAddr *net.UDPAddr

	Routes map[string]func(interface{}) interface{}

	enc *gob.Encoder
	dec *gob.Decoder
}

//udpAddr return resolved addr for udp connection
func udpAddr(addr string) *net.UDPAddr {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln("failed addr", err)
	}
	return raddr
}

//Server - start server listener
func Server(addr string) error {
	Connection.Conn, err = net.ListenUDP("udp", udpAddr(addr))
	if err != nil {
		return err
	}

	// Connection.gob()

	go Connection.carrier()

	return nil
}

//Client - initialize connection
func Client(addr string) error {
	Connection.Conn, err = net.DialUDP("udp", udpAddr(":0"), udpAddr(addr))
	if err != nil {
		return err
	}

	// Connection.remoteAddr = udpAddr(Connection.Conn.LocalAddr().String())

	// log.Println(Connection.Conn.RemoteAddr().String())
	// log.Println(Connection.Conn.LocalAddr().String())

	// Connection.Conn = conn
	// Connection.gob()
	// Connection.Conn.WriteTo([]byte("asd"), Connection.Conn.RemoteAddr())

	go Connection.carrier()

	return nil
}

// func (c *connection) gob() {
// 	c.enc = gob.NewEncoder(c.Conn)
// 	// c.dec = gob.NewDecoder(c.buf)
// }

//carrier manage incoming messages
func (c *connection) carrier() {
	for {
		var b = make([]byte, 2048)
		i, addr, err := c.Conn.ReadFromUDP(b)
		if err != nil {
			log.Println("failed read udp package, error: ", err)
			continue
		}
		b = b[:i]
		// log.Println(i, addr, err)
		// log.Printf(string(b))

		buf := bytes.NewBuffer(b)
		var m Message
		if err := gob.NewDecoder(buf).Decode(&m); err != nil {
			log.Println("failed decode message, error:", err)
			continue
		}
		log.Println("get new message, call function:", m.ResponseFunc)
		// m.remoteAddr = c.Conn.RemoteAddr()

		if f, ok := c.Routes[m.RequestFunc]; ok {
			data := m.Response(f(m.Data))
			c.Conn.WriteToUDP(data, addr)
		} else {
			log.Printf("called function `%s` is not found\n")
		}
		// c.Handler(m)
	}
	log.Println("wtf? exit?")
}

func AddRoute(funcname string, f func(interface{}) interface{}) {
	Connection.Routes[funcname] = f
}

//Message structure, Time contains timestamp when message was sent
type Message struct {
	Time         time.Time
	RequestFunc  string
	ResponseFunc string

	Data interface{}
}

//SendMessage to server or client
func SendMessage(reqfunc, resfunc string, data interface{}) error {
	m := Message{
		Time:         time.Now(),
		RequestFunc:  reqfunc,
		ResponseFunc: resfunc,
		Data:         data,
	}

	buf := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(buf).Encode(m)
	if err != nil {
		return err
	}

	_, err = Connection.Conn.Write(buf.Bytes())
	// _, err = Connection.Conn.WriteToUDP(buf.Bytes(), udpAddr(Connection.Conn.RemoteAddr().String()))
	// _, err = Connection.Conn.WriteTo(buf.Bytes(), Connection.Conn.RemoteAddr())
	// _, err = Connection.Conn.WriteToUDP(buf.Bytes(), Connection.remoteAddr)
	return err

	// log.Println("send message to:", reqfunc, "with answer to:", resfunc)
	// return gob.New
}

func (m Message) Response(data interface{}) []byte {
	if data == nil || m.ResponseFunc == "" {
		return nil
	}

	newM := Message{
		Time:        time.Now(),
		RequestFunc: m.ResponseFunc,
		Data:        data,
	}

	// Connection.Conn.WriteToUDP([]byte{}, m.remoteAddr)
	log.Println("send response to:", m.ResponseFunc)
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(newM)
	if err != nil {
		log.Println("failed send response message, error:", err)
	}

	return buf.Bytes()
}

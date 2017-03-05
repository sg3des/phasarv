package network

import (
	"encoding/gob"
	"fmt"
	"log"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type TServer struct{}
type TClient struct{}

func (TServer) HanName0(req *Request) interface{} {
	s, ok := req.Data.(Data)
	if !ok {
		return nil
	}
	fmt.Println("server get message:", s.Msg)
	return Data{"hello i`m server"}
}
func (TServer) HanName1(req *Request) interface{} {
	return nil
}
func (TClient) HanClient0(req *Request) interface{} {
	s, ok := req.Data.(Data)
	if !ok {
		return nil
	}
	fmt.Println("client get message:", s.Msg)
	return nil
}

type Data struct {
	Msg string
}

func TestNewConnection(t *testing.T) {
	gob.Register(Data{})

	var hs TServer
	s := NewConnection(hs)
	if err := s.Server("127.0.0.1:9690"); err != nil {
		t.Error(err)
	}
	defer s.Close()

	time.Sleep(1 * time.Second)

	var hc TClient
	c := NewConnection(hc)
	if err := c.Client("127.0.0.1:9690"); err != nil {
		t.Error(err)
	}
	defer c.Close()

	err := c.SendMessage(hs.HanName0, hc.HanClient0, Data{"hello i`m client"})
	if err != nil {
		t.Error(err)
	}

	c2 := NewConnection(hc)
	if err := c2.Client("127.0.0.1:9690"); err != nil {
		t.Error(err)
	}
	defer c2.Close()

	time.Sleep(1 * time.Second)
	c.Close()
	time.Sleep(1 * time.Second)
	err = c2.SendMessage(hs.HanName0, hc.HanClient0, Data{"hello i`m client-2"})
	if err != nil {
		t.Error(err)
	}
	c2.Close()

	time.Sleep(1 * time.Second)
	s.Close()
}

// func TestServerClient(t *testing.T) {
// 	gob.Register(Data{})

// 	s := NewHandlers(map[string]Handler{"server": serverHandler})
// 	if err := s.Server(":9690"); err != nil {
// 		t.Error(err)
// 	}
// 	defer s.Close()

// 	c := NewHandlers(map[string]Handler{"client": clientHalder})
// 	if err := c.Client("127.0.0.1:9690"); err != nil {
// 		t.Error(err)
// 	}
// 	defer c.Close()

// 	c2 := NewHandlers(map[string]Handler{"client": clientHalder})
// 	if err := c2.Client("127.0.0.1:9690"); err != nil {
// 		t.Error(err)
// 	}
// 	defer c2.Close()

// 	c.SendMessage("server", "client", getData("client1"))
// 	c2.SendMessage("server", "client", getData("client2"))
// 	// time.Sleep(1 * time.Second)
// 	// s.Broadcast("client", "", "i`m server")

// 	time.Sleep(1 * time.Second)
// 	return
// 	log.Println(s, c)
// }

// type Data struct {
// 	Name string
// 	S    *SubData
// 	s    SubData

// 	Func func()
// }

// type SubData struct {
// 	SubName string
// }

// func getData(name string) *Data {
// 	d := &Data{
// 		Name: "TestName: " + name,
// 		S: &SubData{
// 			SubName: "TestSubName: " + name,
// 		},
// 		s: SubData{
// 			SubName: "Unexported",
// 		},
// 	}

// 	return d
// }

// func serverHandler(req *Request) interface{} {
// 	log.Printf("	S: %++v", req.Data)
// 	// return nil
// 	return "hello from server"
// }

// func clientHalder(req *Request) interface{} {
// 	log.Printf("	C: %++v", req.Data)
// 	return nil
// 	// return "response from client"
// }

// func BenchmarkServerClient(b *testing.B) {
// 	c := NewHandlers(map[string]Handler{"pong": pong})
// 	if err := c.Client("127.0.0.1:9692"); err != nil {
// 		b.Error(err)
// 	}
// 	// defer c.Close()

// 	for i := 0; i < b.N; i++ {
// 		if err := c.SendMessage("ping", "pong", i); err != nil {
// 			// b.Log(i)
// 			b.Error(err)
// 		}
// 	}
// }

// func pong(req *Request) interface{} {
// 	return nil
// }

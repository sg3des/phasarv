package network

import (
	"log"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestServerClient(t *testing.T) {
	s := NewHandlers(map[string]Handler{"server": serverHandler})
	if err := s.Server(":9690"); err != nil {
		t.Error(err)
	}
	defer s.Close()

	c := NewHandlers(map[string]Handler{"client": clientHalder})
	if err := c.Client("127.0.0.1:9690"); err != nil {
		t.Error(err)
	}
	defer c.Close()

	c2 := NewHandlers(map[string]Handler{"client": clientHalder})
	if err := c2.Client("127.0.0.1:9690"); err != nil {
		t.Error(err)
	}
	defer c2.Close()

	c.SendMessage("server", "client", "hello from client")
	c2.SendMessage("server", "client", "hello i`m client2")
	// time.Sleep(1 * time.Second)
	// s.Broadcast("client", "", "i`m server")

	time.Sleep(1 * time.Second)
	return
	log.Println(s, c)
}

func serverHandler(req *Request) interface{} {
	log.Printf("	S: %++v", req.Data)
	// return nil
	return "hello from server"
}

func clientHalder(req *Request) interface{} {
	log.Printf("	C: %++v", req.Data)
	return nil
	// return "response from client"
}

func BenchmarkServerClient(b *testing.B) {
	c := NewHandlers(map[string]Handler{"pong": pong})
	if err := c.Client("127.0.0.1:9692"); err != nil {
		b.Error(err)
	}
	// defer c.Close()

	for i := 0; i < b.N; i++ {
		if err := c.SendMessage("ping", "pong", i); err != nil {
			// b.Log(i)
			b.Error(err)
		}
	}
}

func pong(req *Request) interface{} {
	return nil
}

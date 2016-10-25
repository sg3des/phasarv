package main

import (
	"log"
	"network"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func BenchmarkServer(b *testing.B) {
	if err := network.Client("127.0.0.1" + addr); err != nil {
		b.Error("failed connect to server, error:", err)
	}

	data := []byte("ho-ho-ho")
	for i := 0; i < b.N; i++ {
		network.SendMessage("funcname", "asd", data)
	}
}

package main

import (
	"game"
	"game/db"
	controllers "game/network-controllers"
	"game/rooms"
	"strings"
	// "game/rooms"
	"scenes"

	"log"
	"math/rand"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/sg3des/argum"

	"engine"
)

var args struct {
	Room string `help:"show specify room"`
	Mode string `argum:"pos,|client" help:"current mode"`
}

func init() {
	log.SetFlags(log.Lshortfile)
	argum.MustParse(&args)
}

func main() {
	game.NeedRender = true
	engine.Client(local, keyCallback)
}

func local() {
	scenes.InitEnvironment()

	if args.Room != "" {
		rooms.Init()

		switch strings.ToLower(args.Room) {
		case "index":
			rooms.Index()
			return
		case "hangar":
			rooms.Hangar(randomName())
			return
		default:
			log.Fatalln("unknown room %s", args.Room)
		}
	}

	if args.Mode == "client" {
		go networkPlay()
	} else {
		localPlay()
	}
}

func networkPlay() {
	controllers.Connect("127.0.0.1:9696")
	controllers.SendAuthorize(randomName())
}

func localPlay() {
	scenes.Load("scene00")

	player := db.GetPlayer(randomName())
	player.Ship.LeftWeapon = db.GetWeapon("laser0")
	game.CreateLocalPlayer(player)
	// for i := 0; i < 10; i++ {
	// 	game.CreateEnemy()
	// }
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomName() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		// w.SetShouldClose(true)
		engine.Close()
	}
	// if key == glfw.KeyF1 && action == glfw.Press {
	// 	glfw.SwapInterval(1)
	// }
	if key == glfw.KeySpace && action == glfw.Press {
		for _, p := range game.Players {
			p.Ship.Destroy()
		}
	}
}

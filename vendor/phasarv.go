package main

import (
	"engine"
	"game"
	"game/db"
	"log"
	"math/rand"
	"scenes"
	"time"

	"github.com/sg3des/argum"
)

var args struct {
	TestMode string `argum:"pos" help:"run a specify section for testing"`
}

func init() {
	log.SetFlags(log.Lshortfile)
	argum.MustParse(&args)
}

func main() {
	game.NeedRender = true
	db.SetInitialValues()

	engine.Start(gamemode)
}

func gamemode() {
	scenes.InitEnvironment()

	switch args.TestMode {
	case "battle":
		localBattle()
	case "network":
		networkBattle()
	case "hangar":
		u, _ := db.LookupUser(randomName(), "")
		game.NewHangar(u)
	default:
		game.Start()
	}
}

//
// TEST MODES
//

func networkBattle() {
	// controllers.Connect("127.0.0.1:9696")
	// controllers.SendAuthorize(randomName())
}

func localBattle() {
	user, _ := db.LookupUser(randomName(), "")

	user.Ship.LeftWeapon = db.GetWeapon("laser0")
	user.Ship.RightWeapon = db.GetWeapon("rocket0")

	game.SinglePlay = true
	battle := game.NewBattle(game.NewBattleConfig())
	battle.SetLocalPlayer(user)

	log.Println(battle)

	// battle.CreateEnemy()
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

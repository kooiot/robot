package main

import (
	"math/rand"
	"time"

	"github.com/kooiot/robot/cmd/robot-pong/sub"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	sub.Execute()
}

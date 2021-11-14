package main

import (
	"math/rand"
	"time"

	"github.com/kooiot/robot/cmd/server/sub"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	sub.Execute()
}

package main

import (
	"math/rand"
	"time"

	"github.com/kooiot/robot/cmd/client/sub"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	sub.Execute()
}

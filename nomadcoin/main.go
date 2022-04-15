package main

import (
	"github.com/m00p1ng/learn-go/nomadcoin/cli"
	"github.com/m00p1ng/learn-go/nomadcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}

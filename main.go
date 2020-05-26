package main

import (
	"github.com/jpillora/overseer"
	"log"
	"os"
	"test.local/cmd"
)

func main() {
	overseer.SanityCheck() // 平滑重启
	if err := cmd.Execute(); err != nil {
		log.Fatalf("star err: %v", err)
		os.Exit(1)
	}
}

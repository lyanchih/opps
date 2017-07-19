package main

import (
	_ "gitlab.adlinktech.com/lyan.hung/opps/engine/all"
	"log"
)

func main() {
	cmd, err := newRootCommand().ExecuteC()
	if err != nil {
		log.Fatal(err)
	}

	waitDone(cmd)
}

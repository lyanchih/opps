package main

import (
	_ "gitlab.adlinktech.com/lyan.hung/opps/engine/all"
	"log"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

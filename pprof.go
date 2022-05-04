package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func pprof() {
	if os.Getenv("__ENV__") != "DEV" {
		return
	}

	log.Println(http.ListenAndServe(":6060", nil))
}

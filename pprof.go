package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func pprof() {
	log.Println(http.ListenAndServe(":6060", nil))
}

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	models "zoubun/internal"
)

var endpoint = "http://localhost:3000"

func main() {
	motd := flag.NewFlagSet("motd", flag.ExitOnError)

	incCmd := flag.NewFlagSet("increment", flag.ExitOnError)
	countCmd := flag.NewFlagSet("count", flag.ExitOnError)

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case motd.Name():
		index()
	case incCmd.Name():
		increment()
	case countCmd.Name():
		count()
	default:
		usage()
	}
}

func usage() {
	log.Print("USAGE: cli SUBCOMMAND")
	log.Print("motd displays the message of the day")
	log.Print("increment contributes towards the counting")
	log.Print("count displays the current count")
}

func index() {
	resp, err := http.Get(endpoint + "/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var jsonOutput models.Motd

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		panic(err)
	}

	log.Printf("Message of the day: %v", jsonOutput.Message)
}

func count() {
	resp, err := http.Get(endpoint + "/count")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var jsonOutput models.Counter

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		panic(err)
	}

	log.Printf("Count %v", jsonOutput.Count)
}

func increment() {
	jsonData, err := json.Marshal("{}")
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(endpoint+"/increment", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var jsonOutput models.Counter

	err = json.NewDecoder(resp.Body).Decode(&jsonOutput)
	if err != nil {
		panic(err)
	}

	log.Printf("Incremented to %v", jsonOutput.Count)
}

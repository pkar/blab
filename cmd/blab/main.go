package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pkar/blab"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	host := flag.String("host", "localhost", "The host interface to listen on")
	port := flag.Int("port", 7777, "The port the server accepts connections for clients")
	logDir := flag.String("logs", "", "The directory to write chat logs for each room to. If empty no logs will be written")
	flag.Parse()

	conf := &blab.Config{
		Host:   *host,
		Port:   *port,
		LogDir: *logDir,
	}

	if err := blab.Run(conf); err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}

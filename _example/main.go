package main

import (
	"log"
	"os"
	"strconv"

	"github.com/soopsio/nkill"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("Kills all processes listening on the given TCP ports.\nusage: nkill port")
	}

	// if os.Getpid() != 0 {
	// 	log.Println("WARNING: You are not running this script as superuser.")
	// }

	for _, port := range os.Args[1:] {
		p, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			log.Printf("%s is not a valid port number\n", port)
			continue
		}
		nkill.KillPort(p)
	}

}

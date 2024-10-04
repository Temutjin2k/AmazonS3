package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"triple-s/internal/handler"
	"triple-s/internal/utils"
)

func main() {
	portFlag := flag.Int("port", 8080, "Port number")
	dirFlag := flag.String("dir", "", "Path to the directory")
	helpFlap := flag.Bool("help", false, "Help flag")
	flag.Parse()

	if *helpFlap || !utils.ValidateFlags(*portFlag, *dirFlag) {
		utils.PrintHelp()
		os.Exit(0)
	}

	address := fmt.Sprintf(":%d", *portFlag)
	url := "http://localhost" + address + "/" + *dirFlag

	fmt.Printf("Starting server on port: %v\nURL: %v\n", *portFlag, url)

	http.HandleFunc("/", handler.Handler)
	log.Fatal(http.ListenAndServe(address, nil))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"triple-s/config"
	"triple-s/internal/handler"
	"triple-s/utils"
)

func main() {
	portFlag := flag.Int("port", 8080, "Port number")
	dirFlag := flag.String("dir", "data", "Path to the directory")
	helpFlap := flag.Bool("help", false, "Help flag")
	flag.Parse()

	if *helpFlap || !utils.ValidateFlags(*portFlag, *dirFlag) {
		utils.PrintHelp()
		os.Exit(0)
	}
	config.Dir = "./" + *dirFlag
	err := utils.MakeDir(config.Dir)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	port := fmt.Sprintf(":%d", *portFlag)
	url := "http://localhost" + port + "/"

	fmt.Printf("Starting server on port: %v\nURL: %v\n", *portFlag, url)

	http.HandleFunc("/", handler.Handler)
	log.Fatal(http.ListenAndServe(port, nil))
}

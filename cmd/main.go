package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"triple-s/config"
	"triple-s/internal/handler"
	"triple-s/utils"
)

func main() {
	portFlag := flag.Int("port", 8080, "Port number")
	dirFlag := flag.String("dir", "data", "Path to the directory")
	helpFlap := flag.Bool("help", false, "Help flag")
	flag.Parse()

	if *helpFlap {
		utils.PrintHelp()
		os.Exit(0)
	}

	if (*dirFlag)[0] != '/' {
		config.Dir = "./" + *dirFlag
	} else {
		config.Dir = *dirFlag
	}

	config.Dir = filepath.Clean(config.Dir)
	fmt.Println("Created directory to store buckets in", config.Dir)

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

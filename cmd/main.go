package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"triple-s/internal/utils"
)

var (
	portFlag *int
	dirFlag  *string
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL path: %v", r.URL.Path)
	fmt.Fprintf(w, "PortFlag: %v\nDirFlag: %v", *portFlag, *dirFlag)
}

func main() {
	portFlag = flag.Int("port", 8080, "Port number")
	dirFlag = flag.String("dir", "", "Path to the directory")
	helpFlap := flag.Bool("help", false, "Help flag")
	flag.Parse()

	if *helpFlap || !utils.ValidateFlags(*portFlag, *dirFlag) {
		utils.PrintHelp()
		os.Exit(0)
	}

	address := fmt.Sprintf(":%d", *portFlag)
	url := "http://localhost" + address + "/" + *dirFlag

	fmt.Printf("Starting server on port: %v\nURL: %v\n", *portFlag, url)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(address, nil))
}

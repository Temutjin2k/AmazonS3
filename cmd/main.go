package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"triple-s/internal/utils"
)

// var (
// 	portFlag *int
// 	dirFlag  *string
// )

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello it's triple's\nURL path: %v", r.URL.Path)
	if r.Method == "PUT" { // Create a Bucket:
		folderName := "./data" + r.URL.Path
		err := os.Mkdir(folderName, 0o700) // 0755 is the permission mode
		if err != nil {
			if os.IsExist(err) {
				fmt.Println("Folder already exists")
			} else {
				fmt.Println("Error creating folder:", err)
			}
		} else {
			fmt.Println("Folder created successfully:", folderName)
		}
	}

	// Respond with a 200 OK status
	w.WriteHeader(http.StatusOK) // 200 OK
}

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

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(address, nil))
}

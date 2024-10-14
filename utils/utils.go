package utils

import (
	"fmt"
	"os"
	"time"
	"triple-s/config"
)

func PrintHelp() {
	helpMessage := `Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`

	fmt.Println(helpMessage)
}

func MakeDir(path string) error {
	// Check if the directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If the directory does not exist, create it
		err := os.MkdirAll(path, 0o755)
		if err != nil {
			return err
		}
	}

	// Check if the metadata file exists
	metadataPath := path + "/buckets.csv"
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		// If the metadata file does not exist, create it
		err := os.WriteFile(metadataPath, config.BucketMetadataFields, 0o755)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCurrentTimeStamp() string {
	currentTime := time.Now()                         // Get the current time
	formattedTime := currentTime.Format(time.RFC3339) // Format the time to the desired layout
	return formattedTime
}

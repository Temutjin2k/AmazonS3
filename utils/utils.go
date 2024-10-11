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
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return err
	}

	// Creating buckets.csv to store metadata of bucket
	err = os.WriteFile(path+"/buckets.csv", config.BucketMetadataFields, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentTimeStamp() string {
	currentTime := time.Now()                         // Get the current time
	formattedTime := currentTime.Format(time.RFC3339) // Format the time to the desired layout
	return formattedTime
}

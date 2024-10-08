package utils

import "time"

func GetCurrentTimeStamp() string {
	currentTime := time.Now()                         // Get the current time
	formattedTime := currentTime.Format(time.RFC3339) // Format the time to the desired layout
	return formattedTime
}

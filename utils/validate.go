package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"triple-s/config"
)

// Validates Url and also returns what command is it, is bucketName/ObjectName valid and error message
func ValidateURL(url string) (string, bool, error) {
	commandType := ""
	if url == "/" {
		return config.HandlerBucketList, true, nil
	}

	url = strings.Trim(url, "/")
	endPoints := strings.Split(url, "/")

	switch len(endPoints) {
	case 1:
		commandType = config.HandlerBucket
	case 2:
		commandType = config.HandlerObject

		isValid, err := isValidObjectKey(endPoints[1])
		if !isValid {
			return "", false, err
		}
	default:
		return "", false, config.ErrNotFound
	}

	isValid, err := isValidBucketName(endPoints[0])
	if !isValid {
		return "", false, err
	}
	return commandType, true, nil
}

// isValidBucketName checks if the provided bucket name adheres to the naming conventions.
func isValidBucketName(name string) (bool, error) {
	// Length check
	if len(name) < 3 || len(name) > 63 {
		return false, fmt.Errorf("bucket name must be between 3 and 63 characters")
	}

	// Allowed characters
	if !regexp.MustCompile(`^[a-z0-9.-]+$`).MatchString(name) {
		return false, fmt.Errorf("bucket name can only contain lowercase letters, numbers, hyphens, and dots")
	}

	// Cannot start or end with a hyphen
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return false, fmt.Errorf("bucket name cannot start or end with a hyphen")
	}

	// Cannot contain consecutive periods or hyphens
	if strings.Contains(name, "..") || strings.Contains(name, "--") {
		return false, fmt.Errorf("bucket name cannot contain consecutive periods or hyphens")
	}

	// Cannot be formatted as an IP address
	ipAddressRegex := `^(\d{1,3}\.){3}\d{1,3}$`
	if regexp.MustCompile(ipAddressRegex).MatchString(name) {
		return false, fmt.Errorf("bucket name cannot be formatted as an IP address")
	}

	return true, nil
}

// isValidObjectKey checks if ObjectKey named correctly corresponding to rules
func isValidObjectKey(objectKey string) (bool, error) {
	// Length check
	if len(objectKey) < 1 || len(objectKey) > 1024 {
		return false, errors.New("object key must be between 1 and 1024 characters")
	}

	// Allowed characters regex
	allowedCharacters := `^[a-zA-Z0-9!_.\-'()*\/]+$`
	if !regexp.MustCompile(allowedCharacters).MatchString(objectKey) {
		return false, errors.New("object key can only contain alphanumeric characters and special characters (! - _ . * ' ( ) /)")
	}

	// Check for leading or trailing spaces
	if strings.HasPrefix(objectKey, " ") || strings.HasSuffix(objectKey, " ") {
		return false, errors.New("object key cannot start or end with a space")
	}

	// Check for consecutive slashes
	if strings.Contains(objectKey, "//") {
		return false, errors.New("object key cannot contain consecutive slashes")
	}

	return true, nil
}

/*
#### Bucket Naming Conventions:
    Bucket names must be unique across the system.
    Names should be between 3 and 63 characters long.
    Only lowercase letters, numbers, hyphens (-), and dots (.) are allowed.
    Must not be formatted as an IP address (e.g., 192.168.0.1).
    Must not begin or end with a hyphen and must not contain two consecutive periods or dashes.


### Object key naming guidelines

You can use any UTF-8 character in an object key name. However, using certain characters in key names can cause problems with some applications and protocols. The following guidelines help you maximize compliance with DNS, web-safe characters, XML parsers, and other APIs.
Safe characters

The following character sets are generally safe for use in key names.
Alphanumeric characters
    0-9
    a-z
    A-Z
Special characters
    Exclamation point (!)
    Hyphen (-)
    Underscore (_)
    Period (.)
    Asterisk (*)
    Single quote (')
    Open parenthesis (()
    Close parenthesis ())
*/

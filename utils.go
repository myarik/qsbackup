package qsbackup

import "os"

// FileExists checks if specified file exists.
func IsExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

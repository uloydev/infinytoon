package utils

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

func SaveToken(token string) (err error) {
	var file *os.File

	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// check if .infinitoon file exists in user home directory
	_, err = os.Stat(path.Join(homeDir, ConfigFile))
	if err != nil {
		// create .infinitoon file
		file, err = os.Create(path.Join(homeDir, ConfigFile))
	} else {
		// open .infinitoon file
		file, err = os.OpenFile(path.Join(homeDir, ConfigFile), os.O_RDWR, 0644)
		// clear file content
		if err == nil {
			err = file.Truncate(0)
		}
	}

	defer file.Close()
	// check if there is an error
	if err != nil {
		return
	}

	// write token to .infinitoon file
	_, err = file.WriteString(token)
	return
}

func ReadToken() (token string, err error) {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	// read .infinitoon file
	tokenByte, err := os.ReadFile(path.Join(homeDir, ConfigFile))
	token = string(tokenByte)
	return
}

func DefaultPreRun(cmd *cobra.Command, args []string) {
	cmd.Println(Banner)
}

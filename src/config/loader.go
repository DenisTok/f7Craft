package config

import "os"

func Load() error {
	err := os.MkdirAll(DBDir, 0700)
	if err != nil {
		return err
	}
	return err
}

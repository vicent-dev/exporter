package main

import (
	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load(".env")

	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"testing"
)

func TestInit(t *testing.T) {
	RootDir = "/home/cloud"
	Validators = 1

	err := initFilesWithConfig()
	if err != nil {
		return
	}
}

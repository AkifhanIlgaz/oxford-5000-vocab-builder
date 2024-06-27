package main

import (
	"fmt"

	"github.com/AkifhanIlgaz/vocab-builder/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	_ = cfg
}

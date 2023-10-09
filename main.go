package main

import (
	"fmt"

	"github.com/AkifhanIlgaz/vocab-builder/setup"
)

func main() {
	err := setup.Run(":3000")
	if err != nil {
		fmt.Println(err)
	}
}

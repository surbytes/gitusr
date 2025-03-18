package main

import (
	"fmt"
	"os"

	"github.com/surbytes/gitusr/utils"
)

func main() {
	if len(os.Args) > 2 && os.Args[1] == "switch" {
		err := utils.SwitchUserByName(os.Args[2])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	} else {

		utils.RenderUsers()
	}
}
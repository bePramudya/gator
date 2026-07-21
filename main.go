package main

import (
	"fmt"
	"os"

	"github.com/bePramudya/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = cfg.SetUser("bagus"); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(cfg)
}

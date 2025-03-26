package main

import (
	"fmt"

	"github.com/ChipsAhoyEnjoyer/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(c.DBUrl)
	fmt.Println(c.CurrentUsername)
	err = c.SetUser("ChipsAhoyEnjoyer")
	if err != nil {
		panic(err)
	}
	c, err = config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(c.DBUrl)
	fmt.Println(c.CurrentUsername)
}

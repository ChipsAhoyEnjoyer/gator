package main

import (
	"fmt"

	"github.com/ChipsAhoyEnjoyer/gator/internal/config"
)

func main() {
	fmt.Println("Hello World")
	c, err := config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(c.DBUrl)
	fmt.Println(c.CurrentUsername)
	c.CurrentUsername = "ChipsAhoyEnjoyer"
	err = c.SetUser()
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

package main

import (
	"fmt"

	"github.com/hari-govind/liveserver-go/config"
)

func main() {
	fmt.Println(config.GetConfig().Depth)
}

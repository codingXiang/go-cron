package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())
	c.Start()
	spec := "* * * * * *"
	c.AddFunc(spec, func() {
		fmt.Println("hi")
	})
	fmt.Println(c.Entries())
	select {}
}


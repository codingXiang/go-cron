package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())
	spec := "*/2 * * * * *"
	c.AddFunc(spec, func() {
		fmt.Println("hi")
	})
	fmt.Println(c.Entries())
	c.Start()
	select {}
}

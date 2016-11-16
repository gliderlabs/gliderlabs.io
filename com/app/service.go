package app

import "fmt"

func (c *Component) Serve() {
	fmt.Println("Hello world")
	select {}
}

func (c *Component) Stop() {
	fmt.Println("Stopping")
}

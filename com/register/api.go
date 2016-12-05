package register

import (
	"github.com/gliderlabs/gliderlabs.io/com/api"
	"github.com/progrium/duplex/golang"
)

func (c *Component) RegisterAPI(_ *api.Component) {
	api.Register("register.hello", func(ch *duplex.Channel) error {
		return ch.SendLast("hello")
	})
}

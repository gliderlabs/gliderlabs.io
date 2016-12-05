package api

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/progrium/duplex/golang"
)

func init() {
	com.Register("api", &Component{})
}

type Contributor interface {
	RegisterAPI(api *Component)
}

func Contributors() []Contributor {
	var contributors []Contributor
	for _, com := range com.Enabled(new(Contributor), nil) {
		contributors = append(contributors, com.(Contributor))
	}
	return contributors
}

// Component ...
type Component struct{}

// wrapping package exported Register until all packages migrated
func (c *Component) Register(name string, handler func(*duplex.Channel) error) {
	Register(name, handler)
}

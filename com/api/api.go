package api

import (
	"github.com/progrium/duplex/golang"
	"golang.org/x/net/context"

	"github.com/gliderlabs/comlab/pkg/log"
)

var rpc = duplex.NewRPC(duplex.NewJSONCodec())

const ipKey = 0

// Register API method
func Register(name string, handler func(*duplex.Channel) error) {
	rpc.Register(name, func(ch *duplex.Channel) error {
		// TODO: wrap channel to capture more logging detail
		log.Info(log.Fields{"method": name}, GetIP(ch.Context()))
		return handler(ch)
	})
}

func GetIP(ctx context.Context) string {
	return ctx.Value(ipKey).(string)
}

// ErrorUnexpected ...
func ErrorUnexpected(ch *duplex.Channel, err error) error {
	log.Info(GetIP(ch.Context()), err)
	return ch.SendErr(1000, "Unexpected error", err.Error())
}

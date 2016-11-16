package app

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("app", &Component{})
}

type Component struct{}

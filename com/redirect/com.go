package redirect

import (
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("redirect", &Component{})
}

type Component struct{}

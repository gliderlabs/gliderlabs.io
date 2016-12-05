package register

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("register", &Component{},
		com.Option("key_secret", "764efa883dda1e11db47671c4a3bbd9e", "Secret to sign registration keys with"),
		com.Option("stripe_secret_key", "", "Stripe secret key"),
		com.Option("stripe_pub_key", "", "Stripe publishable key"))

}

// Component component
type Component struct{}

package api

func (c *Component) AppPreStart() error {
	for _, contributor := range Contributors() {
		contributor.RegisterAPI(c)
	}
	return nil
}

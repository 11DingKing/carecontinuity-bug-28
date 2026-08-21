package serviceactivation

type Coordinator struct{ routes *ServiceRoutes }

func NewCoordinator() *Coordinator {
	return &Coordinator{routes: NewServiceRoutes(RoutePolicy{Mode: "eager", RetainFailed: true})}
}
func (c *Coordinator) Activate(region, provider string, commit func() error) error {
	return c.routes.Activate(region, provider, commit)
}
func (c *Coordinator) Owner(region string) (string, bool) { return c.routes.Owner(region) }

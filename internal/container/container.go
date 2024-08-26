package container

import (
	"go.uber.org/fx"
)

type Container struct {
	Modules   []fx.Option
	Servers   []any
	Providers []any

	invokers []any
}

func NewContainer() *Container {
	c := &Container{}
	return c
}

func (a *Container) AddModules(modules ...fx.Option) {
	a.Modules = append(a.Modules, modules...)
}

func (a *Container) AddServers(servers ...any) {
	a.Servers = append(a.Servers, servers...)
}

func (a *Container) AddProviders(providers ...any) {
	a.Providers = append(a.Providers, providers...)
}

func (a *Container) AddInvokers(invokers ...any) {
	a.invokers = append(a.invokers, invokers...)
}

func (a *Container) Run() {
	opts := []fx.Option{}

	opts = append(opts, a.Modules...)
	opts = append(opts, fx.Provide(a.Providers...))
	opts = append(opts, fx.Provide(a.Servers...))
	opts = append(opts, fx.Invoke(a.invokers...))

	fx.New(opts...).Run()
}

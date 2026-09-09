package bootstrap

import "github.com/MontFerret/ferret/v2/pkg/module"

type registrationModule func(module.Bootstrap) error

func (m registrationModule) Name() string {
	return "bootstrap-test"
}

func (m registrationModule) Register(boot module.Bootstrap) error {
	return m(boot)
}

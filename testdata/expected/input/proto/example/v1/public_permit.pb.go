package examplev1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *PublicRequest) GetChecks() pkg.CheckConfig {
	return pkg.CheckConfig{
		Type:   pkg.PUBLIC,
		Checks: []pkg.Check{},
	}
}

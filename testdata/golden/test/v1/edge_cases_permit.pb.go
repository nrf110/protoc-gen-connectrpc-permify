package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *EmptyRequest) GetChecks() pkg.CheckConfig {
	return pkg.CheckConfig{
		Type:   pkg.PUBLIC,
		Checks: []pkg.Check{},
	}
}

func (req *MinimalResource) GetChecks() pkg.CheckConfig {
	return pkg.CheckConfig{
		Type:   pkg.PUBLIC,
		Checks: []pkg.Check{},
	}
}

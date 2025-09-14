package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *EmptyRequest) GetChecks() pkg.CheckConfig {
	return pkg.CheckConfig{
		IsPublic: true,
		Checks:   []pkg.Check{},
	}
}

func (req *MinimalResource) GetChecks() pkg.CheckConfig {
	return pkg.CheckConfig{
		IsPublic: true,
		Checks:   []pkg.Check{},
	}
}

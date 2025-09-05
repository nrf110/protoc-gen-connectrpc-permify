package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *UpdateUserRequest) GetChecks() pkg.CheckConfig {
	permission := "write"
	var checks []pkg.Check
	resource := req
	var id string
	if resource.Id != "" {
		id = resource.Id
	}
	tenantId := "default"
	if resource.CompanyId != "" {
		tenantId = resource.CompanyId
	}
	attributes := make(map[string]any)
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "User",
			ID:         id,
			Attributes: attributes,
		},
	}
	checks = append(checks, check)
	return pkg.CheckConfig{
		Type:   pkg.SINGLE,
		Checks: checks,
	}
}

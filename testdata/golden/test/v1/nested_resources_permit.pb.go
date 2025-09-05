package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *NestedRequest) GetChecks() pkg.CheckConfig {
	permission := "manage"
	var checks []pkg.Check
	resource := req.Organization
	var id string
	if resource.Id != "" {
		id = resource.Id
	}
	tenantId := "default"
	if resource.TenantId != "" {
		tenantId = resource.TenantId
	}
	attributes := make(map[string]any)
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Organization",
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

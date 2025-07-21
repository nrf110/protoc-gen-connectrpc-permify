package examplev1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *AttributesRequest) GetChecks() pkg.CheckConfig {
	permission := "read"
	var checks []pkg.Check
	resource := req
	tenantId := "default"
	if resource.CompanyId != "" {
		tenantId = resource.CompanyId
	}
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type: "Attributes",
			ID:   resource.Id,
		},
	}
	checks = append(checks, check)
	return pkg.CheckConfig{
		Type:   pkg.SINGLE,
		Checks: checks,
	}
}

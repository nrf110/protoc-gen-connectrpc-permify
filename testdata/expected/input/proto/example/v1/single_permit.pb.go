package examplev1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *ResourceRequest) GetChecks() pkg.CheckConfig {
	permission := "create"
	var checks []pkg.Check
	tenantId := "default"
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type: "Flat",
		},
	}
	checks = append(checks, check)
	return pkg.CheckConfig{
		Type:   pkg.SINGLE,
		Checks: checks,
	}
}

func (req *ResourceWithIdRequest) GetChecks() pkg.CheckConfig {
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
			Type: "Flat",
			ID:   resource.Id,
		},
	}
	checks = append(checks, check)
	return pkg.CheckConfig{
		Type:   pkg.SINGLE,
		Checks: checks,
	}
}

func (req *NestedResourceRequest) GetChecks() pkg.CheckConfig {
	permission := "edit"
	var checks []pkg.Check
	resource := req.Resource
	tenantId := "default"
	if resource != nil && resource.NestedIds.CompanyId != "" {
		tenantId = resource.NestedIds.CompanyId
	}
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type: "Nested",
			ID:   resource.NestedIds.Id,
		},
	}
	checks = append(checks, check)
	return pkg.CheckConfig{
		Type:   pkg.SINGLE,
		Checks: checks,
	}
}

package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *ResourceRequest) GetChecks() pkg.CheckConfig {
	permission := "create"
	var checks []pkg.Check
	var id string
	tenantId := "default"
	attributes := make(map[string]any)
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Flat",
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

func (req *ResourceWithIdRequest) GetChecks() pkg.CheckConfig {
	permission := "read"
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
			Type:       "Flat",
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

func (req *NestedResourceRequest) GetChecks() pkg.CheckConfig {
	permission := "edit"
	var checks []pkg.Check
	resource := req.Resource
	var id string
	if resource != nil && resource.NestedIds.Id != "" {
		id = resource.NestedIds.Id
	}
	tenantId := "default"
	if resource != nil && resource.NestedIds.CompanyId != "" {
		tenantId = resource.NestedIds.CompanyId
	}
	attributes := make(map[string]any)
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Nested",
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

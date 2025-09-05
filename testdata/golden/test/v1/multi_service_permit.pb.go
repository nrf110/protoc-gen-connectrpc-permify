package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *Account) GetChecks() pkg.CheckConfig {
	permission := "read"
	var checks []pkg.Check
	resource := req
	var id string
	if resource.Id != "" {
		id = resource.Id
	}
	tenantId := "default"
	if resource.OrgId != "" {
		tenantId = resource.OrgId
	}
	attributes := make(map[string]any)
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Account",
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

func (req *PublicInfo) GetChecks() pkg.CheckConfig {
	return pkg.CheckConfig{
		Type:   pkg.PUBLIC,
		Checks: []pkg.Check{},
	}
}

func (req *Profile) GetChecks() pkg.CheckConfig {
	permission := "read"
	var checks []pkg.Check
	resource := req
	var id string
	if resource.Id != "" {
		id = resource.Id
	}
	tenantId := "default"
	attributes := make(map[string]any)
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Profile",
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

func (req *Settings) GetChecks() pkg.CheckConfig {
	permission := "read"
	var checks []pkg.Check
	resource := req
	var id string
	if resource.Id != "" {
		id = resource.Id
	}
	tenantId := "default"
	attributes := make(map[string]any)
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Settings",
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

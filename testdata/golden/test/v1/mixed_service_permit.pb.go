package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *SimpleRequest) GetChecks() pkg.CheckConfig {
	return pkg.CheckConfig{
		Type:   pkg.PUBLIC,
		Checks: []pkg.Check{},
	}
}

func (req *GetUserResource) GetChecks() pkg.CheckConfig {
	permission := "read"
	var checks []pkg.Check
	resource := req
	var id string
	if resource.UserId != "" {
		id = resource.UserId
	}
	tenantId := "default"
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

func (req *UpdateUserResource) GetChecks() pkg.CheckConfig {
	permission := "write"
	var checks []pkg.Check
	resource := req
	var id string
	if resource.UserId != "" {
		id = resource.UserId
	}
	tenantId := "default"
	if resource.CompanyId != "" {
		tenantId = resource.CompanyId
	}
	attributes := make(map[string]any)
	attributes["email"] = resource.Email
	attributes["role"] = resource.Role
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

func (req *DeleteUserResource) GetChecks() pkg.CheckConfig {
	permission := "admin"
	var checks []pkg.Check
	resource := req
	var id string
	if resource.UserId != "" {
		id = resource.UserId
	}
	tenantId := "default"
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

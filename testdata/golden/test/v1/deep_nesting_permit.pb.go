package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *DeepNestedRequest) GetChecks() pkg.CheckConfig {
	permission := "process"
	var checks []pkg.Check
	resource := req.Container.Level2.Resource
	var id string
	if resource.Id != "" {
		id = resource.Id
	}
	tenantId := "default"
	attributes := make(map[string]any)
	attributes["level3_data"] = resource.Data
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Level3",
			ID:         id,
			Attributes: attributes,
		},
	}
	checks = append(checks, check)
	return pkg.CheckConfig{
		IsPublic: false,
		Checks:   checks,
	}
}

func (req *VeryDeepResource) GetChecks() pkg.CheckConfig {
	permission := "admin"
	var checks []pkg.Check
	resource := req
	var id string
	if resource != nil && resource.Level1 != nil && resource.Level1.Level2 != nil && resource.Level1.Level2.Level3 != nil && resource.Level1.Level2.Level3.Ids.DeepId != "" {
		id = resource.Level1.Level2.Level3.Ids.DeepId
	}
	tenantId := "default"
	if resource != nil && resource.Level1 != nil && resource.Level1.Level2 != nil && resource.Level1.Level2.Level3 != nil && resource.Level1.Level2.Level3.Ids.DeepTenant != "" {
		tenantId = resource.Level1.Level2.Level3.Ids.DeepTenant
	}
	attributes := make(map[string]any)
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "VeryDeep",
			ID:         id,
			Attributes: attributes,
		},
	}
	checks = append(checks, check)
	return pkg.CheckConfig{
		IsPublic: false,
		Checks:   checks,
	}
}

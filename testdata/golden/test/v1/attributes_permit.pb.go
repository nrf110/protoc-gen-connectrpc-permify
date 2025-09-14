package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *AttributesRequest) GetChecks() pkg.CheckConfig {
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
	attributes["complex"] = resource.Complex
	var fooValues []any
	for _, v1 := range resource.Mapped {
		fooValues = append(fooValues, v1.Bar)
	}
	if len(fooValues) > 0 {
		attributes["foo"] = fooValues
	}
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Attributes",
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

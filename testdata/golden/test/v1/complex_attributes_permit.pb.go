package testv1

import (
	pkg "github.com/nrf110/connectrpc-permify/pkg"
)

func (req *ComplexResource) GetChecks() pkg.CheckConfig {
	permission := "process"
	var checks []pkg.Check
	resource := req
	var id string
	if resource.Id != "" {
		id = resource.Id
	}
	tenantId := "default"
	if resource.TenantId != "" {
		tenantId = resource.TenantId
	}
	attributes := make(map[string]any)
	var categoryValues []any
	for _, mp9o25urbv := range resource.Attributes {
		categoryValues = append(categoryValues, mp9o25urbv.Category)
	}
	if len(categoryValues) > 0 {
		attributes["category"] = categoryValues
	}
	var priorityValues []any
	for _, ptwauf := range resource.Attributes {
		priorityValues = append(priorityValues, ptwauf.Priority)
	}
	if len(priorityValues) > 0 {
		attributes["priority"] = priorityValues
	}
	attributes["tags"] = resource.Tags
	attributes["department"] = resource.Department
	check := pkg.Check{
		TenantID:   tenantId,
		Permission: permission,
		Entity: &pkg.Resource{
			Type:       "Document",
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

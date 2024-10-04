package examplev1

import (
	connectrpc_permit "github.com/nrf110/connectrpc-permit"
)

func (req *ResourceRequest) GetChecks() connectrpc_permit.CheckConfig {
	action := "create"
	var checks []connectrpc_permit.Check
	tenantId := "default"
	check := connectrpc_permit.Check{
		Action: action,
		Resource: connectrpc_permit.Resource{
			Type:   "Flat",
			Tenant: tenantId,
		},
	}
	checks = append(checks, check)
	checkType := connectrpc_permit.SINGLE
	if len(checks) > 1 {
		checkType = connectrpc_permit.BULK
	}
	return connectrpc_permit.CheckConfig{
		Type:   checkType,
		Mode:   connectrpc_permit.AllOf,
		Checks: checks,
	}
}

func (req *ResourceWithIdRequest) GetChecks() connectrpc_permit.CheckConfig {
	action := "read"
	var checks []connectrpc_permit.Check
	resource := req
	tenantId := "default"
	if resource.CompanyId != "" {
		tenantId = resource.CompanyId
	}
	check := connectrpc_permit.Check{
		Action: action,
		Resource: connectrpc_permit.Resource{
			Type:   "Flat",
			Key:    resource.Id,
			Tenant: tenantId,
		},
	}
	checks = append(checks, check)
	checkType := connectrpc_permit.SINGLE
	if len(checks) > 1 {
		checkType = connectrpc_permit.BULK
	}
	return connectrpc_permit.CheckConfig{
		Type:   checkType,
		Mode:   connectrpc_permit.AllOf,
		Checks: checks,
	}
}

func (req *NestedResourceRequest) GetChecks() connectrpc_permit.CheckConfig {
	action := "edit"
	var checks []connectrpc_permit.Check
	resource := req.Resource
	tenantId := "default"
	if resource != nil && resource.NestedIds.CompanyId != "" {
		tenantId = resource.NestedIds.CompanyId
	}
	check := connectrpc_permit.Check{
		Action: action,
		Resource: connectrpc_permit.Resource{
			Type:   "Nested",
			Key:    resource.NestedIds.Id,
			Tenant: tenantId,
		},
	}
	checks = append(checks, check)
	checkType := connectrpc_permit.SINGLE
	if len(checks) > 1 {
		checkType = connectrpc_permit.BULK
	}
	return connectrpc_permit.CheckConfig{
		Type:   checkType,
		Mode:   connectrpc_permit.AllOf,
		Checks: checks,
	}
}

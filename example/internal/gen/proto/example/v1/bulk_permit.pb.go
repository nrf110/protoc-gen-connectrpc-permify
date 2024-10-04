package examplev1

import (
	connectrpc_permit "github.com/nrf110/connectrpc-permit"
)

func (req *RepeatedResourceRequest) GetChecks() connectrpc_permit.CheckConfig {
	action := "create"
	var checks []connectrpc_permit.Check
	for range req.Resources {
		tenantId := "default"
		check := connectrpc_permit.Check{
			Action: action,
			Resource: connectrpc_permit.Resource{
				Type:   "Repeated",
				Tenant: tenantId,
			},
		}
		checks = append(checks, check)
	}
	checkType := connectrpc_permit.SINGLE
	if len(checks) > 1 {
		checkType = connectrpc_permit.BULK
	}
	return connectrpc_permit.CheckConfig{
		Type:   checkType,
		Mode:   connectrpc_permit.AnyOf,
		Checks: checks,
	}
}

func (req *RepeatedResourceWithIdsRequest) GetChecks() connectrpc_permit.CheckConfig {
	action := "edit"
	var checks []connectrpc_permit.Check
	for _, l := range req.Resources {
		resource := l
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
	}
	checkType := connectrpc_permit.SINGLE
	if len(checks) > 1 {
		checkType = connectrpc_permit.BULK
	}
	return connectrpc_permit.CheckConfig{
		Type:   checkType,
		Mode:   connectrpc_permit.AnyOf,
		Checks: checks,
	}
}

func (req *RepeatedResourceWithRepeatedIdsRequest) GetChecks() connectrpc_permit.CheckConfig {
	action := "edit"
	var checks []connectrpc_permit.Check
	for _, g := range req.Resources {
		resource := g
		for _, unqf5f17gp := range resource.Repeats {
			for _, ia9v5elfwa4 := range unqf5f17gp.Ids {
				tenantId := "default"
				check := connectrpc_permit.Check{
					Action: action,
					Resource: connectrpc_permit.Resource{
						Type:   "Repeated",
						Key:    ia9v5elfwa4,
						Tenant: tenantId,
					},
				}
				checks = append(checks, check)
			}
		}
	}
	checkType := connectrpc_permit.SINGLE
	if len(checks) > 1 {
		checkType = connectrpc_permit.BULK
	}
	return connectrpc_permit.CheckConfig{
		Type:   checkType,
		Mode:   connectrpc_permit.AnyOf,
		Checks: checks,
	}
}

func (req *MappedResourceRequest) GetChecks() connectrpc_permit.CheckConfig {
	action := "create"
	var checks []connectrpc_permit.Check
	for _, lt55nmzij5r := range req.Mapped {
		resource := lt55nmzij5r
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
	}
	checkType := connectrpc_permit.SINGLE
	if len(checks) > 1 {
		checkType = connectrpc_permit.BULK
	}
	return connectrpc_permit.CheckConfig{
		Type:   checkType,
		Mode:   connectrpc_permit.AnyOf,
		Checks: checks,
	}
}

func (req *MappedResourceWithNestedMappedIdsRequest) GetChecks() connectrpc_permit.CheckConfig {
	action := "edit"
	var checks []connectrpc_permit.Check
	for _, xn00rqv := range req.Mapped {
		for _, soian := range xn00rqv.Resources {
			resource := soian
			for _, frvw := range resource.Nested {
				for _, rf := range frvw.Ids {
					tenantId := "default"
					check := connectrpc_permit.Check{
						Action: action,
						Resource: connectrpc_permit.Resource{
							Type:   "NestedMapped",
							Key:    rf,
							Tenant: tenantId,
						},
					}
					checks = append(checks, check)
				}
			}
		}
	}
	checkType := connectrpc_permit.SINGLE
	if len(checks) > 1 {
		checkType = connectrpc_permit.BULK
	}
	return connectrpc_permit.CheckConfig{
		Type:   checkType,
		Mode:   connectrpc_permit.AnyOf,
		Checks: checks,
	}
}

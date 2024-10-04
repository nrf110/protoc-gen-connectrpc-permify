package examplev1

import (
	connectrpc_permit "github.com/nrf110/connectrpc-permit"
)

func (req *PublicRequest) GetChecks() connectrpc_permit.CheckConfig {
	return connectrpc_permit.CheckConfig{
		Type:   connectrpc_permit.PUBLIC,
		Mode:   connectrpc_permit.AllOf,
		Checks: []connectrpc_permit.Check{},
	}
}

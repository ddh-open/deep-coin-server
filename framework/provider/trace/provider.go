package trace

import (
	framework2 "devops-http/framework"
	"devops-http/framework/contract"
)

type NiceTraceProvider struct {
	c framework2.Container
}

// Register registe a new function for make a service instance
func (provider *NiceTraceProvider) Register(c framework2.Container) framework2.NewInstance {
	return NewNiceTraceService
}

// Boot will called when the service instantiate
func (provider *NiceTraceProvider) Boot(c framework2.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *NiceTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *NiceTraceProvider) Params(c framework2.Container) []interface{} {
	return []interface{}{provider.c}
}

// Name / Name define the name for this service
func (provider *NiceTraceProvider) Name() string {
	return contract.TraceKey
}

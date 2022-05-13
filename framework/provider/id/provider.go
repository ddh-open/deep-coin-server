package id

import (
	framework2 "devops-http/framework"
	"devops-http/framework/contract"
)

type NiceIDProvider struct {
}

// Register register a new function for make a service instance
func (provider *NiceIDProvider) Register(c framework2.Container) framework2.NewInstance {
	return NewNiceIDService
}

// Boot will called when the service instantiate
func (provider *NiceIDProvider) Boot(c framework2.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *NiceIDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *NiceIDProvider) Params(c framework2.Container) []interface{} {
	return []interface{}{}
}

// Name / Name define the name for this service
func (provider *NiceIDProvider) Name() string {
	return contract.IDKey
}

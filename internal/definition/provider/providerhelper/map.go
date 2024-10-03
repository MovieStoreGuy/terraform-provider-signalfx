package providerhelper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDefinitionFunc allows for packages to define their
// resource name and definition as part of the same package.
type ResourceDefinitionFunc func() (string, *schema.Resource)

// New is a convenience function to make it clear of intent
func (fn ResourceDefinitionFunc) New() (string, *schema.Resource) {
	return fn()
}

func MustNewResourceMap(definitions ...ResourceDefinitionFunc) map[string]*schema.Resource {
	rm := make(map[string]*schema.Resource, len(definitions))
	for _, df := range definitions {
		name, resource := df.New()
		if _, exist := rm[name]; exist {
			panic("entry already exists: " + name)
		}
		rm[name] = resource
	}
	return rm
}

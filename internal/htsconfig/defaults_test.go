// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults_test tests module defaults
package htsconfig

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"

	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	d := DefaultConfiguration
	props := d.Container.ServerProps
	reads := d.Container.ReadsConfig
	variants := d.Container.VariantsConfig

	// SERVER PROPS
	assert.Equal(t, props.Host, htsconstants.DfltServerPropsHost)
	assert.Equal(t, props.Port, htsconstants.DfltServerPropsPort)

	// READS DATA SOURCE REGISTRY
	assert.Equal(t, *reads.Enabled, true)

	// VARIANTS DATA SOURCE REGISTRY
	assert.Equal(t, *variants.Enabled, true)
}

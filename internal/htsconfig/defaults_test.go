// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults_test tests module defaults
package htsconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultsTC = []struct {
	key, exp string
}{
	{"port", "3000"},
	{"host", "http://localhost:3000"},
}

var defaultsReadsSourcesRegistryTC = []struct {
	expPattern, expPath string
}{
	{
		"^tabulamuris\\.(?P<accession>10X.*)$",
		"https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/{accession}_possorted_genome.bam",
	},
	{
		"^tabulamuris\\.(?P<accession>.*)$",
		"https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/{accession}.mus.Aligned.out.sorted.bam",
	},
}

var defaultsVariantsSourcesRegistryTC = []struct {
	expPattern, expPath string
}{
	{
		"^1000genomes\\.(?P<accession>.*)$",
		"https://ftp-trace.ncbi.nih.gov/1000genomes/ftp/phase1/analysis_results/integrated_call_sets/{accession}.vcf.gz",
	},
}

func TestDefaults(t *testing.T) {
	d := getDefaults()
	for _, tc := range defaultsTC {
		assert.Equal(t, tc.exp, d[tc.key])
	}
}

func TestDefaultReadsSourcesRegistry(t *testing.T) {
	d := getDefaultReadsSourcesRegistry()
	for i := 0; i < len(d.Sources); i++ {
		tc := defaultsReadsSourcesRegistryTC[i]
		assert.Equal(t, tc.expPattern, d.Sources[i].Pattern)
		assert.Equal(t, tc.expPath, d.Sources[i].Path)
	}
}

func TestDefaultVariantsSourcesRegistry(t *testing.T) {
	d := getDefaultVariantsSourcesRegistry()
	for i := 0; i < len(d.Sources); i++ {
		tc := defaultsVariantsSourcesRegistryTC[i]
		assert.Equal(t, tc.expPattern, d.Sources[i].Pattern)
		assert.Equal(t, tc.expPath, d.Sources[i].Path)
	}
}

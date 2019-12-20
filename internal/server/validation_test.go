package server

import (
	"testing"
)

var params []map[string][]string = []map[string][]string{
	map[string][]string{
		"format":        []string{"Bam"},
		"class":         []string{"Header"},
		"referenceName": []string{"*"},
		"start":         []string{"500"},
		"end":           []string{"9000"},
		"fields":        []string{"QNAME", "SEQ"},
	},
	map[string][]string{
		"format":        []string{"CRAM"},
		"class":         []string{"body"},
		"referenceName": []string{"chr11"},
		"start":         []string{"-100"},
		"end":           []string{"0"},
		"fields":        []string{"QNAME,NotARealField"},
	},
	map[string][]string{
		"format":        []string{"Bam,Cram"},
		"class":         []string{"Header,body"},
		"referenceName": []string{"Chr1"},
		"end":           []string{"900"},
		"fields":        []string{"Tlen,SEQ"},
	},
	map[string][]string{
		"format": []string{"Bam", "Cram"},
		"start":  []string{"100"},
		"end":    []string{"1000"},
	},
	map[string][]string{
		"referenceName": []string{"Chr1"},
	},
	map[string][]string{
		"referenceName": []string{"chr1"},
		"start":         []string{"11000"},
	},
	map[string][]string{
		"referenceName": []string{"chr1"},
		"start":         []string{"100"},
		"end":           []string{"10"},
	},
}

}

func TestParseQueryClass(t *testing.T) {

}

func TestParseRefName(t *testing.T) {

}

func TestParseRange(t *testing.T) {

}

func TestParseFields(t *testing.T) {

}

func TestValidReadFormat(t *testing.T) {

}

func TestValidClass(t *testing.T) {

}

func TestValidRange(t *testing.T) {

}

func TestValidFields(t *testing.T) {

}

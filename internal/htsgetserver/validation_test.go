package htsgetserver

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

func TestParseFormat(t *testing.T) {
	format, err := parseFormat(params[0])
	if format != "BAM" {
		t.Errorf("Got: %v, expected: BAM", format)
	}
	if err != nil {
		t.Errorf("Got non-nil error: %v, expected nil error", err.Error())
	}

	format, err = parseFormat(params[1])
	if err == nil {
		t.Errorf("Got nil error, expected error with message: Unsupported format")
	}

	format, err = parseFormat(params[2])
	if err == nil {
		t.Errorf("Got nil error, expected error with message: Unsupported format")
	}

	format, err = parseFormat(params[3])
	if format != "BAM" {
		t.Errorf("Got: %v, expected: BAM", format)
	}
	if err != nil {
		t.Errorf("Got non-nil error: %v, expected nil error", err.Error())
	}

	format, err = parseFormat(params[4])
	if format != "BAM" {
		t.Errorf("Got: %v, expected: BAM", format)
	}
	if err != nil {
		t.Errorf("Got non-nil error: %v, expected nil error", err.Error())
	}
}

func TestParseQueryClass(t *testing.T) {
	class, err := parseQueryClass(params[0])
	if class != "header" {
		t.Errorf("Got: %v, expected: header", class)
	}
	if err != nil {
		t.Errorf("Got non-nil error: %v, expected nil error", err.Error())
	}

	class, err = parseQueryClass(params[1])
	if err == nil {
		t.Errorf("Got nil error, expected error with message: InvalidInput")
	}

	class, err = parseQueryClass(params[2])
	if err == nil {
		t.Errorf("Got nil error, expected error with message: InvalidInput")
	}

	class, err = parseQueryClass(params[3])
	if class != "" {
		t.Errorf("Got: %v, expected an empty string", class)
	}
	if err != nil {
		t.Errorf("Got non-nil error: %v, expected nil error", err.Error())
	}
}

func TestParseRefName(t *testing.T) {
	name := parseRefName(params[0])
	if name != "*" {
		t.Errorf("Got: %v, expected: *", name)
	}

	name = parseRefName(params[1])
	if name != "chr11" {
		t.Errorf("Got: %v, expected: chr11", name)
	}

	name = parseRefName(params[3])
	if name != "" {
		t.Errorf("Got: %v, expected an empty string", name)
	}
}

func TestParseRange(t *testing.T) {
	start, end, err := parseRange(params[0], "*")
	if err == nil {
		t.Errorf("Got nil error, expected error: InvalidRange")
	}

	start, end, err = parseRange(params[1], "chr11")
	if err == nil {
		t.Errorf("Got nil error, expected error: InvalidRange")
	}

	start, end, err = parseRange(params[2], "Chr1")
	if err == nil {
		t.Errorf("Got nil error, expected error: InvalidRange")
	}

	start, end, err = parseRange(params[3], "")
	if err == nil {
		t.Errorf("Got nil error, expected error: InvalidRange")
	}

	start, end, err = parseRange(params[4], "Chr1")
	if start != "-1" {
		t.Errorf("Got: %v, expected: -1", start)
	}
	if end != "-1" {
		t.Errorf("got: %v, expected: -1", end)
	}
	if err != nil {
		t.Errorf("Got non-nil error: %v, expected nil error", err.Error())
	}

	start, end, err = parseRange(params[5], "")
	if start != "11000" {
		t.Errorf("Got: %v, expected: 11000", start)
	}
	if end != "-1" {
		t.Errorf("Got: %v, expected: -1", end)
	}
	if err != nil {
		t.Errorf("Got non-nil error: %v, expected nil error", err.Error())
	}

	start, end, err = parseRange(params[6], "")
	if err == nil {
		t.Errorf("Got non-nil error, expected error: InvalidInput")
	}
}

func TestParseFields(t *testing.T) {
	fields, err := parseFields(params[0])
	if len(fields) != 1 {
		t.Errorf("Got a slice of fields with length: %v, expected length: 1", len(fields))
	}

	if fields[0] != "QNAME" {
		t.Errorf("First element of fields is: %v, expected: QNAME", fields[0])
	}
	if err != nil {
		t.Errorf("Got non-nil error: %v, expected nil error", err.Error())
	}

	fields, err = parseFields(params[1])
	if err == nil {
		t.Errorf("Got nil error, expected error: InvalidInput")
	}

	fields, err = parseFields(params[2])
	if len(fields) != 2 {
		t.Errorf("Got a slice of fields with length: %v, expected length: 2", len(fields))
	}
	if fields[0] != "TLEN" {
		t.Errorf("First element of fields is: %v, expected: TLEN", fields[0])
	}
	if fields[1] != "SEQ" {
		t.Errorf("Second element of fields is: %v, expected: SEQ", fields[1])
	}
}

func TestValidReadFormat(t *testing.T) {
	valid := validReadFormat("BAM")
	if !valid {
		t.Errorf("Got false with BAM file format, expected true")
	}
}

func TestValidClass(t *testing.T) {
	valid := validClass("header")
	if !valid {
		t.Errorf("Got false with header class, expected true")
	}

	valid = validClass("body")
	if !valid {
		t.Errorf("Got false with body class, expected true")
	}
}

func TestValidRange(t *testing.T) {
	valid := validRange("100", "50", "chr10")
	if valid {
		t.Errorf("Got true, expected false with start > end")
	}

	valid = validRange("100", "1000", "chr10")
	if !valid {
		t.Errorf("Got false, expected true")
	}

	valid = validRange("100", "100", "")
	if valid {
		t.Errorf("Got true, expected false")
	}

	valid = validRange("100", "100", "*")
	if valid {
		t.Errorf("Got true, expected false")
	}

	valid = validRange("-100", "100", "chr10")
	if valid {
		t.Errorf("Got true, expected false")
	}
}

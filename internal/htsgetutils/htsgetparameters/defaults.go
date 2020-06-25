package htsgetparameters

var defaultScalarParameterValues = map[string]string{
	"id":            "",
	"format":        "BAM",
	"class":         "",
	"referenceName": "*",
	"start":         "-1",
	"end":           "-1",
}

var defaultListParameterValues = map[string][]string{
	"fields": {},
	"tags":   {},
	"notags": {},
}

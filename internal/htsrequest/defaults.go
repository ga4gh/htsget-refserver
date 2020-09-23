// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module defaults.go contains default values for each parameter
package htsrequest

var defaultId = ""
var defaultFormat = "BAM"
var defaultClass = ""
var defaultReferenceName = ""
var defaultStart = -1
var defaultEnd = -1
var defaultFields = []string{"ALL"}
var defaultTags = []string{"ALL"}
var defaultNoTags = []string{"NONE"}
var defaultHtsgetBlockClass = ""

var defaultParameterValues = map[string]interface{}{
	"id":               defaultId,
	"format":           defaultFormat,
	"class":            defaultClass,
	"referenceName":    defaultReferenceName,
	"start":            defaultStart,
	"end":              defaultEnd,
	"fields":           defaultFields,
	"tags":             defaultTags,
	"notags":           defaultNoTags,
	"HtsgetBlockClass": defaultHtsgetBlockClass,
}

// Package htsgetrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
package htsgetrequest

import (
	"bufio"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/config"
	"github.com/ga4gh/htsget-refserver/internal/htsgeterror"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
)

// the correct validation function for each request parameter name
var validationByParam = map[string]func(string, *HtsgetRequest) (bool, string){
	"id":               validateID,
	"format":           validateFormat,
	"class":            validateClass,
	"referenceName":    validateReferenceNameExists,
	"start":            validateStart,
	"end":              validateEnd,
	"fields":           validateFields,
	"tags":             noValidation,
	"notags":           validateNoTags,
	"HtsgetBlockClass": validateClass,
	"HtsgetBlockId":    noValidation,
	"HtsgetNumBlocks":  noValidation,
}

// the correct error to raise for each request parameter validation
var errorsByParam = map[string]func(http.ResponseWriter, *string){
	"id":               htsgeterror.NotFound,
	"format":           htsgeterror.UnsupportedFormat,
	"class":            htsgeterror.InvalidInput,
	"referenceName":    htsgeterror.InvalidRange,
	"start":            htsgeterror.InvalidRange,
	"end":              htsgeterror.InvalidRange,
	"fields":           htsgeterror.InvalidInput,
	"tags":             htsgeterror.InvalidInput,
	"notags":           htsgeterror.InvalidInput,
	"HtsgetBlockClass": htsgeterror.InvalidInput,
	"HtsgetBlockId":    htsgeterror.InternalServerError,
	"HtsgetNumBlocks":  htsgeterror.InternalServerError,
}

// helper function, determines if a string can be parsed as an integer
func isInteger(value string) bool {
	_, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return true
}

// helper function, determines if a string can be parsed as an integer, and
// is greater than zero
func isGreaterThanEqualToZero(value string) bool {
	if !isInteger(value) {
		return false
	}
	num, _ := strconv.Atoi(value)
	if num < 0 {
		return false
	}
	return true
}

// empty validation function for request parameters that do not need to be
// validated. always returns true
func noValidation(empty string, htsgetReq *HtsgetRequest) (bool, string) {
	return true, ""
}

// validates the 'id' path parameter. checks if an object matching the 'id'
// could be found from the data source
func validateID(id string, htsgetReq *HtsgetRequest) (bool, string) {
	res, err := http.Head(config.DATA_SOURCE_URL + htsgetutils.FilePath(id))
	if err != nil {
		return false, "The requested resource was not found"
	}
	res.Body.Close()
	if res.Status == "404 Not Found" {
		return false, "The requested resource was not found"
	}
	return true, ""
}

// validates the 'format' query string parameter. checks if the requested
// format is one of the allowed options
func validateFormat(format string, htsgetReq *HtsgetRequest) (bool, string) {
	switch strings.ToUpper(format) {
	case "BAM":
		return true, ""
	case "CRAM":
		return false, "CRAM not supported" // currently not supported
	default:
		return false, "file format: '" + format + "' not supported"
	}
}

// validates the 'class' query string parameter. checks if the requested
// class is one of the allowed options
func validateClass(class string, htsgetReq *HtsgetRequest) (bool, string) {
	switch strings.ToLower(class) {
	case "header":
		return true, ""
	case "body":
		return false, "'body' only requests currently not supported" // currently not supported
	default:
		return false, "class: '" + class + "' not supported"
	}
}

// validates the 'referenceName' query string parameter. checks if the requested
// reference contig/chromosome is in the BAM/CRAM header sequence dictionary
func validateReferenceNameExists(referenceName string, htsgetReq *HtsgetRequest) (bool, string) {
	id := htsgetReq.ID()
	cmd := exec.Command("samtools", "view", "-H", config.DATA_SOURCE_URL+htsgetutils.FilePath(id))
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return false, "Could not access requested file"
	}
	cmd.Start()
	reader := bufio.NewReader(pipe)
	l, _, err := reader.ReadLine()

	for ; err == nil; l, _, err = reader.ReadLine() {
		if strings.Contains(string(l), "SN:"+referenceName) {
			return true, ""
		}
	}
	cmd.Wait()
	return false, "invalid 'referenceName': " + referenceName
}

// validates the 'start' query string parameter. checks that it is a valid,
// non-zero integer, and that it is being used correctly in conjunction with
// 'referenceName'
func validateStart(start string, htsgetReq *HtsgetRequest) (bool, string) {
	referenceName := htsgetReq.ReferenceName()
	if referenceName == "*" || referenceName == "" {
		return false, "'start' cannot be set without 'referenceName'"
	}

	if !isInteger(start) {
		return false, "'start' is not a valid integer"
	}

	if !isGreaterThanEqualToZero(start) {
		return false, "'start' must be greater than or equal to zero"
	}

	return true, ""
}

// validates the 'end' query string parameter. checks that it is a valid,
// non-zero integer, that it's being used correctly in conjunction with
// 'referenceName', and that the end coordinate is greater than the start
// coordinate
func validateEnd(end string, htsgetReq *HtsgetRequest) (bool, string) {
	referenceName := htsgetReq.ReferenceName()
	start := htsgetReq.Start()
	if referenceName == "*" || referenceName == "" {
		return false, "'end' cannot be set without 'referenceName'"
	}

	if !isInteger(end) {
		return false, "'end' is not a valid integer"
	}

	if !isGreaterThanEqualToZero(end) {
		return false, "'end' must be greater than or equal to zero"
	}

	if start != "-1" {
		startNum, startErr := strconv.Atoi(start)
		endNum, endErr := strconv.Atoi(end)
		if startErr != nil || endErr != nil {
			return false, "error converting 'start' and/or 'end' to integers"
		}
		if startNum >= endNum {
			return false, "'end' MUST be higher than 'start'"
		}
	}

	return true, ""
}

// validates the 'fields' query string parameter. checks that every requested
// field is an acceptable value (an expected BAM/CRAM field)
func validateFields(fields string, htsgetReq *HtsgetRequest) (bool, string) {
	fieldsList := splitAndUppercase(fields)
	for _, fieldItem := range fieldsList {
		if _, ok := config.BAM_FIELDS[fieldItem]; !ok {
			return false, "'" + fieldItem + "' not an acceptable field"
		}
	}
	return true, ""
}

// validates the 'notags' query string parameter. checks that there is no
// overlap between tags included by 'tags' and tags excluded by 'notags'
func validateNoTags(notags string, htsgetReq *HtsgetRequest) (bool, string) {
	tagsList := htsgetReq.Tags()
	notagsList := splitOnComma(notags)

	for _, tagItem := range tagsList {
		for _, notagItem := range notagsList {
			if tagItem == notagItem {
				return false, "'" + tagItem + "' cannot be in both 'tags' and 'notags'"
			}
		}
	}
	return true, ""
}

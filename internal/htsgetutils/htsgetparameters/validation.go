package htsgetparameters

import (
	"bufio"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/config"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils/htsgeterror"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils/htsgethttp/htsgetrequest"
)

var validationByParam = map[string]func(string, *htsgetrequest.HtsgetRequest) (bool, string){
	"id":            validateID,
	"format":        validateFormat,
	"class":         validateClass,
	"referenceName": validateReferenceNameExists,
	"start":         validateStart,
	"end":           validateEnd,
	"fields":        validateFields,
	"tags":          noValidation,
	"notags":        validateNoTags,
}

var errorsByParam = map[string]func(http.ResponseWriter, *string){
	"id":            htsgeterror.NotFound,
	"format":        htsgeterror.UnsupportedFormat,
	"class":         htsgeterror.InvalidInput,
	"referenceName": htsgeterror.InvalidRange,
	"start":         htsgeterror.InvalidRange,
	"end":           htsgeterror.InvalidRange,
	"fields":        htsgeterror.InvalidInput,
	"tags":          htsgeterror.InvalidInput,
	"notags":        htsgeterror.InvalidInput,
}

func isInteger(value string) bool {
	_, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return true
}

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

func noValidation(empty string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
	return true, ""
}

func validateID(id string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
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

func validateFormat(format string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
	switch strings.ToUpper(format) {
	case "BAM":
		return true, ""
	case "CRAM":
		return false, "CRAM not supported" // currently not supported
	default:
		return false, "file format: '" + format + "' not supported"
	}
}

func validateClass(class string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
	switch strings.ToLower(class) {
	case "header":
		return true, ""
	case "body":
		return false, "'body' only requests currently not supported" // currently not supported
	default:
		return false, "class: '" + class + "' not supported"
	}
}

func validateReferenceNameExists(referenceName string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
	id := htsgetReq.GetScalar("id")
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

func validateStart(start string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
	referenceName := htsgetReq.GetScalar("referenceName")
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

func validateEnd(end string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
	referenceName := htsgetReq.GetScalar("referenceName")
	start := htsgetReq.GetScalar("start")
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

func validateFields(fields string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
	fieldsList := splitAndUppercase(fields)
	for _, fieldItem := range fieldsList {
		if _, ok := config.BAM_FIELDS[fieldItem]; !ok {
			return false, "'" + fieldItem + "' not an acceptable field"
		}
	}
	return true, ""
}

func validateNoTags(notags string, htsgetReq *htsgetrequest.HtsgetRequest) (bool, string) {
	tagsList := htsgetReq.GetList("tags")
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

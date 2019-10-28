package handler

import (
	"net/url"
	"strconv"
	"strings"
)

func parseFormat(params url.Values) (string, error) {
	if _, ok := params["format"]; ok {
		if validReadFormat(params["format"][0]) {
			return strings.ToUpper(params["format"][0]), nil
		} else {
			panic("UnsupportedFormat")
		}
	} else {
		return "BAM", nil
	}
}

func parseClass(params url.Values) (string, error) {
	if _, ok := params["class"]; ok {
		if validClass(params["class"][0]) {
			return strings.ToLower(params["class"][0]), nil
		} else {
			panic("InvalidInput")
		}
	}
	return "", nil
}
func parseRefName(params url.Values) (string, error) {
	if _, ok := params["referenceName"]; ok {
		return params["referenceName"][0], nil
	}
	return "", nil
}

func parseRange(params url.Values, refName string) (string, string, error) {
	if _, ok := params["start"]; ok {
		if _, ok := params["end"]; ok {
			if validRange(params["start"][0], params["end"][0], refName) {
				start, _ := strconv.ParseUint(params["start"][0], 10, 32)
				end, _ := strconv.ParseUint(params["end"][0], 10, 32)
				return strconv.FormatUint(start, 10), strconv.FormatUint(end, 10), nil
			} else {
				panic("InvalidRange")
			}
		}
	} else if _, ok := params["end"]; ok {
		panic("InvalidRange")
	}
	return "0", "0", nil
}

func parseFields(params url.Values) ([]string, error) {
	if _, ok := params["fields"]; ok {
		fields := strings.Split(params["fields"][0], ",")
		if !validFields(fields) {
			panic("InvalidInput")
		}
		return fields, nil
	}
	return []string{}, nil
}

func validReadFormat(s string) bool {
	switch strings.ToUpper(s) {
	case "BAM":
		return true
	case "CRAM":
		return true
	default:
		return false
	}
}

func validClass(s string) bool {
	switch strings.ToLower(s) {
	case "head":
		return true
	case "body":
		return true
	default:
		return false
	}
}

func validRange(startStr string, endStr string, refName string) bool {
	start, errStart := strconv.ParseUint(startStr, 10, 32)
	end, errEnd := strconv.ParseUint(endStr, 10, 32)

	if errStart != nil || errEnd != nil {
		return false
	}
	if start > end {
		return false
	}
	if refName == "" || refName == "*" {
		return false
	}

	return true
}

func validFields(fields []string) bool {
	for _, field := range fields {
		if _, ok := FIELDS[field]; ok {
			return false
		}
	}
	return true
}

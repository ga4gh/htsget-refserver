package server

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func parseFormat(params url.Values) (string, error) {
	if _, ok := params["format"]; ok {
		if validReadFormat(params["format"][0]) {
			return strings.ToUpper(params["format"][0]), nil
		}
		return "", errors.New("Unsupported format")
	}
	return "BAM", nil
}

func parseQueryClass(params url.Values) (string, error) {
	if _, ok := params["class"]; ok {
		class := strings.ToLower(params["class"][0])
		if class == "header" {
			return class, nil
		}
		return "", errors.New("InvalidInput")
	}
	return "", nil
}

func parseRefName(params url.Values) string {
	if _, ok := params["referenceName"]; ok {
		return params["referenceName"][0]
	}
	return ""
}

func parseRange(params url.Values, refName string) (string, string, error) {
	if _, ok := params["start"]; ok {
		if _, ok := params["end"]; ok {
			if validRange(params["start"][0], params["end"][0], refName) {
				return params["start"][0], params["end"][0], nil
			}
			return "0", "0", errors.New("InvalidRange")
		}
		return params["start"][0], "-1", nil
	} else if _, ok := params["end"]; ok {
		return "0", "0", errors.New("InvalidRange")
	}
	return "-1", "-1", nil
}

func parseFields(params url.Values) ([]string, error) {
	if _, ok := params["fields"]; ok {
		fields := strings.Split(params["fields"][0], ",")
		for i := 0; i < len(fields); i++ {
			fields[i] = strings.ToUpper(fields[i])
		}

		if !validFields(fields) {
			return []string{}, errors.New("InvalidInput")
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
	case "header":
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
		if _, ok := FIELDS[field]; !ok {
			return false
		}
	}
	return true
}

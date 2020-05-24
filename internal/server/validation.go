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

func parseTags(params url.Values) []string {
	if _, ok := params["tags"]; ok {
		var tags []string
		if params["tags"][0] == "" {
			tags = []string{""}
			return tags
		}
		tags = strings.Split(params["tags"][0], ",")
		return tags
	}
	return []string{}
}

func parseNoTags(params url.Values, tags []string) ([]string, error) {
	if _, ok := params["notags"]; ok {
		var notags []string
		if params["notags"][0] == "" {
			return []string{}, nil
		}
		notags = strings.Split(params["notags"][0], ",")
		if validNoTags(tags, notags) {
			return notags, nil
		}
		return []string{}, errors.New("InvalidInput")
	}
	return []string{}, nil
}

func validReadFormat(s string) bool {
	switch strings.ToUpper(s) {
	case "BAM":
		return true
	case "CRAM":
		return false // currently not supported
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
	start, errStart := strconv.ParseInt(startStr, 10, 64)
	end, errEnd := strconv.ParseInt(endStr, 10, 64)

	if errStart != nil || errEnd != nil {
		return false
	}
	if start > end {
		return false
	}
	if refName == "" || refName == "*" {
		return false
	}
	if start < 0 || end < 0 {
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

func validNoTags(tags, notags []string) bool {
	for _, tag := range tags {
		for _, notag := range notags {
			if tag == notag {
				return false
			}
		}
	}
	return true
}

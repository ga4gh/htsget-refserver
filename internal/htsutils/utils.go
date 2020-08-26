// Package htsutils provides general, high-level, reusable functions
//
// Module utils defines general, high-level, reusable functions
package htsutils

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// AddTrailingSlash adds a trailing slash to a url if there isn't one already
func AddTrailingSlash(url string) string {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}

// RemoveTrailingSlash removes a trailing slash from a url path if one is there
func RemoveTrailingSlash(url string) string {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	return url
}

// GetTagName accepts and parses a SAM tag to give only its two-letter identifier
func GetTagName(tag string) string {
	return strings.Split(tag, ":")[0]
}

// IsItemInArray checks if a string is present in an array. Returns true if the
// item is in the array, false if not
func IsItemInArray(item string, array []string) bool {

	for i := 0; i < len(array); i++ {
		if item == array[i] {
			return true
		}
	}
	return false
}

// StringIsEmpty checks if a string is empty. Returns true if string is empty,
// false if not
func StringIsEmpty(item string) bool {

	if item == "" {
		return true
	}
	return false
}

// CreateRegexNamedParameterMap creates a map of captured elements from the
// evaluation of a regex pattern (with named capture groups) against a string.
// the map keys will be named according to the capture group names, and
// associated values will be capture group values
func CreateRegexNamedParameterMap(pattern string, s string) (map[string][]string, error) {
	var regex = regexp.MustCompile(pattern)
	matchBool := regex.MatchString(s)
	if !matchBool {
		return nil, errors.New("the supplied string did not match the pattern")
	}

	match := regex.FindStringSubmatch(s)
	matchMap := make(map[string][]string)
	for i, name := range regex.SubexpNames() {
		matchMap[name] = make([]string, 0)
		matchMap[name] = append(matchMap[name], match[i])
	}
	return matchMap, nil
}

// IsValidURL checks if a passed string is a valid url
func IsValidURL(toTest string) bool {
	_, err1 := url.ParseRequestURI(toTest)
	u, err2 := url.Parse(toTest)
	if err1 != nil || err2 != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

// ParseRangeHeader attempts to parse the canonical range header, converting
// start and end to integers if possible. an error is returned if the range
// header could not be successfully parsed
func ParseRangeHeader(rangeHeader string) (int64, int64, error) {
	pattern := "bytes=(?P<start>\\d+)-(?P<end>\\d+)"
	matchMap, err := CreateRegexNamedParameterMap(pattern, rangeHeader)
	if err != nil {
		return 0, 0, err
	}

	start, err := strconv.Atoi(matchMap["start"][0])
	if err != nil {
		return 0, 0, err
	}
	end, err := strconv.Atoi(matchMap["end"][0])
	if err != nil {
		return 0, 0, err
	}
	return int64(start), int64(end), nil
}

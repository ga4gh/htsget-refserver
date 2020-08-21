// Package htsutils provides general, high-level, reusable functions
//
// Module utils.go defines general, high-level, reusable functions
package htsutils

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// AddTrailingSlash adds a trailing slash to a url if there isn't one already
//
// Arguments
//	url (string): url to add slash to
// Returns
//	(string): url with trailing slash (unmodified if url had a slash already)
func AddTrailingSlash(url string) string {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}

func RemoveTrailingSlash(url string) string {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	return url
}

// FilePath gets the correct S3 directory path based on the id of the requested
// file
//
// Arguments
//	id (string): requested object id
// Returns
//	(string): S3 directory containing the object by the given id
func FilePath(id string) string {
	var path string
	if strings.HasPrefix(id, "10X") {
		path = "10x_bam_files/" + id
	} else {
		path = "facs_bam_files/" + id
	}
	return path
}

// GetTagName deconstructs a SAM tag to give only its two-letter name/identifier
//
// Arguments
//	tag (string): a single tag, including name, data type, and value
// Returns
//	(string): the tag name
func GetTagName(tag string) string {
	return strings.Split(tag, ":")[0]
}

// IsItemInArray checks if a string is present in an array
//
// Arguments
//	item (string): the string to check for presence in the array
//	array ([]string): the array to check against
// Returns
//	(bool): true if item was found in array, false if not
func IsItemInArray(item string, array []string) bool {

	for i := 0; i < len(array); i++ {
		if item == array[i] {
			return true
		}
	}
	return false
}

// StringIsEmpty checks if a string is empty
//
// Arguments
//	item (string): the string to check
// Returns
//	(bool): true if string is empty, false if not
func StringIsEmpty(item string) bool {

	if item == "" {
		return true
	}
	return false
}

func CreateRegexNamedParameterMap(pattern string, s string) map[string][]string {
	var regex = regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(s)
	matchMap := make(map[string][]string)
	for i, name := range regex.SubexpNames() {
		matchMap[name] = make([]string, 0)
		matchMap[name] = append(matchMap[name], match[i])
	}
	return matchMap
}

func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func ParseRangeHeader(rangeHeader string) (int64, int64, error) {
	pattern := "bytes=(?P<start>\\d+)-(?P<end>\\d+)"
	matchMap := CreateRegexNamedParameterMap(pattern, rangeHeader)
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

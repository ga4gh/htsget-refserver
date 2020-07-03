// Package htsgetutils provides general, high-level, reusable functions
//
// Module utils.go defines general, high-level, reusable functions
package htsgetutils

import (
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

package htsgetutils

import (
	"strings"
)

func AddTrailingSlash(url string) string {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}

func FilePath(id string) string {
	var path string
	if strings.HasPrefix(id, "10X") {
		path = "10x_bam_files/" + id
	} else {
		path = "facs_bam_files/" + id
	}
	return path
}

func GetTagName(tag string) string {
	return strings.Split(tag, ":")[0]
}

func IsItemInArray(item string, array []string) bool {

	for i := 0; i < len(array); i++ {
		if item == array[i] {
			return true
		}
	}
	return false
}

func StringIsEmpty(item string) bool {

	if item == "" {
		return true
	}
	return false
}

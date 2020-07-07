// Package htsgetconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module datasources.go allows the program to be configured with multiple
// data sources. request id patterns can be mapped to various file path or url
// endpoints, allowing the htsget service to point to many different sources in
// a structured, predictable manner

package htsgetconfig

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
)

type DataSourceRegistry struct {
	Sources []*DataSource `json:"sources"`
}

type DataSource struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path"`
}

func newDataSourceRegistry() *DataSourceRegistry {
	return new(DataSourceRegistry)
}

func (dataSource *DataSource) evaluatePatternMatch(id string) (bool, error) {
	return regexp.MatchString(dataSource.Pattern, id)
}

func (dataSource *DataSource) evaluatePath(id string) (string, error) {

	// create match map, map of named control groups parsed from the regex
	// evaluation of pattern on id
	idParameterMap := htsgetutils.CreateRegexNamedParameterMap(dataSource.Pattern, id)
	parameterNamesMap := htsgetutils.CreateRegexNamedParameterMap("\\{(?P<paramName>.+?)\\}", dataSource.Path)

	finalPath := dataSource.Path
	for i := 0; i < len(parameterNamesMap["paramName"]); i++ {
		paramName := parameterNamesMap["paramName"][i]
		finalPath = strings.Replace(finalPath, "{"+paramName+"}", idParameterMap[paramName][0], -1)
	}
	return finalPath, nil
}

func newDataSource(pattern string, path string) *DataSource {
	dataSource := new(DataSource)
	dataSource.Pattern = pattern
	dataSource.Path = path
	return dataSource
}

func (registry *DataSourceRegistry) addDataSource(dataSource *DataSource) {
	registry.Sources = append(registry.Sources, dataSource)
}

func (registry *DataSourceRegistry) findFirstMatch(id string) (*DataSource, error) {

	for i := 0; i < len(registry.Sources); i++ {
		match, err := registry.Sources[i].evaluatePatternMatch(id)
		if err != nil {
			return nil, err
		}
		if match {
			return registry.Sources[i], nil
		}
	}
	return nil, errors.New("id: " + id + " did not match any registered data sources")
}

func (registry *DataSourceRegistry) GetMatchingPath(id string) (string, error) {
	matchingDataSource, err := registry.findFirstMatch(id)
	if matchingDataSource == nil || err != nil {
		return "", err
	}
	path, err := matchingDataSource.evaluatePath(id)
	if path == "" || err != nil {
		return "", err
	}
	return path, err
}

func (registry *DataSourceRegistry) String() string {
	var builder strings.Builder
	for i := 0; i < len(registry.Sources); i++ {
		builder.WriteString(
			"Path: " + registry.Sources[i].Path + "\t" +
				"Pattern: " + registry.Sources[i].Pattern + "\n",
		)
	}

	return builder.String()
}

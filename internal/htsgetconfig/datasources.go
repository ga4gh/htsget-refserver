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

// DataSourceRegistry holds all data sources for a particular endpoint
//
// Attributes
//	Sources ([]*DataSource): list of data sources to scan
type DataSourceRegistry struct {
	Sources []*DataSource `json:"sources"`
}

// DataSource references a single source of htsget-related data the service
// points to. requested IDs can be mapped to the correct data source (given the
// id itself)
//
// Attributes
//	Pattern (string): regex pattern indicating criteria for an ID to match the data source
//	Path (string): path template, indicating how matching ids can be resolved to an exact location (path or url)
type DataSource struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path"`
}

// newDataSourceRegistry instantiates a data source registry
//
// Returns
//	(*DataSourceRegistry): unpopulated data source registry
func newDataSourceRegistry() *DataSourceRegistry {
	return new(DataSourceRegistry)
}

// evaluatePatternMatch checks if a requested ID matches the data source pattern
//
// 	Type: DataSource
// Arguments
//	id (string): the requested object id
// Returns
//	(bool): if true, the object id fulfills the data source pattern
//	(error): if not nil, an error was encountered in the evaluation process
func (dataSource *DataSource) evaluatePatternMatch(id string) (bool, error) {
	return regexp.MatchString(dataSource.Pattern, id)
}

// evaluatePath completes a url or file path based on the path template and the passed id
//
//	Type: DataSource
// Arguments
//	id (string): requested object id
// Returns
//	(string): populated resource location based on path template and id
//	(error): if not nil, an error was encountered in the evaluation process
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

// newDataSource creates a data source with the given pattern and path template
//
// Arguments
//	pattern (string): new data source regex pattern
//	path (string): new data source path template
// Returns
//	(*DataSource): data source instance
func newDataSource(pattern string, path string) *DataSource {
	dataSource := new(DataSource)
	dataSource.Pattern = pattern
	dataSource.Path = path
	return dataSource
}

// addDataSource adds a data source to the registry
//
//	Type: DataSourceRegistry
// Arguments
//	dataSource (*DataSource): data source to add
func (registry *DataSourceRegistry) addDataSource(dataSource *DataSource) {
	registry.Sources = append(registry.Sources, dataSource)
}

// findFirstMatch gets the first data source in the registry with a pattern matching the requested id
//
//	Type: DataSourceRegistry
// Arguments
//	id (string): requested object id
// Returns
//	(*DataSource): the data source with a pattern matching the id
//	(error): if not nil, a matching data source was not found
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

// GetMatchingPath gets the correct path to the object from the requested id
// the registry is scanned for the first source matching the pattern. once found,
// the path template is populated with the id
//
//	Type: DataSourceRegistry
// Arguments
//	id (string): requested object id
// Returns
//	(string): location to requested resource
//	(error): if not nil, no suitable resource location could be constructed for the id
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

// String gets the registry representation as a string
//
//	Type: DataSourceRegistry
// Returns
//	(string): data source registry string representation
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

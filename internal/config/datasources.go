// Package config allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module datasources.go allows the program to be configured with multiple
// data sources. request id patterns can be mapped to various file path or url
// endpoints, allowing the htsget service to point to many different sources in
// a structured, predictable manner

package config

import (
	"fmt"
	"regexp"
	"strings"
)

type DataSourceRegistry struct {
	sources []*DataSource
}

type DataSource struct {
	Pattern string
	Path    string
}

func newDataSourceRegistry() *DataSourceRegistry {
	return new(DataSourceRegistry)
}

func (dataSource *DataSource) evaluatePatternMatch(id string) (bool, error) {
	fmt.Println("pattern: " + dataSource.Pattern + "\tid: " + id)
	return regexp.MatchString(dataSource.Pattern, id)
}

func newDataSource(pattern string, path string) *DataSource {
	dataSource := new(DataSource)
	dataSource.Pattern = pattern
	dataSource.Path = path
	return dataSource
}

func (registry *DataSourceRegistry) addDataSource(dataSource *DataSource) {
	registry.sources = append(registry.sources, dataSource)
}

func (registry *DataSourceRegistry) FindFirstMatch(id string) (*DataSource, error) {

	for i := 0; i < len(registry.sources); i++ {
		match, err := registry.sources[i].evaluatePatternMatch(id)
		if err != nil {
			return nil, err
		}
		if match {
			return registry.sources[i], nil
		}
	}
	return nil, nil
}

func (registry *DataSourceRegistry) evaluatePath(id string) (string, error) {
	return "", nil
}

func (registry *DataSourceRegistry) GetMatchingPath(id string) (string, error) {
	return "", nil
}

func (registry *DataSourceRegistry) String() string {
	var builder strings.Builder
	for i := 0; i < len(registry.sources); i++ {
		builder.WriteString(
			"Path: " + registry.sources[i].Path + "\t" +
				"Pattern: " + registry.sources[i].Pattern + "\n",
		)
	}

	return builder.String()
}

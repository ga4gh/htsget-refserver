[![Logo](https://www.ga4gh.org/wp-content/themes/ga4gh-theme/gfx/GA-logo-horizontal-tag-RGB.svg)](https://ga4gh.org)

# htsget Reference Server
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![Go Report](https://goreportcard.com/badge/github.com/ga4gh/htsget-refserver)](https://goreportcard.com/badge/github.com/ga4gh/htsget-refserver)

This is a reference implementation of the version 1.2.0 htsget API protocol for securely streaming genomic data. Click [here](https://academic.oup.com/bioinformatics/article/35/1/119/5040320) for a high level overview of the standard, and [here](https://github.com/samtools/hts-specs/blob/master/htsget.md) to view the specification itself. 

## Setup

### Dependencies

* samtools (tested on version 1.9)
* bcftools (tested on version 1.10)

- Install [Golang(v1.13) and language tools](https://golang.org/dl/). With Homebrew, `$ brew install go`

This project uses [Go modules](https://blog.golang.org/using-go-modules) to manage packages and dependencies.

`$ go build -o ./htsget-refserver ./cmd && ./htsget-refserver` from root directory to start the server on the specified port (`3000` by default, see Configuration section).

## Usage
The API is defined at https://github.com/samtools/hts-specs/blob/master/htsget.md. 

This server is deployed at https://htsget.ga4gh.org/.

### Configuration

The web service can be configured with modifiable runtime parameters via environment variables. The table below indicates these modifiable parameters, their default values, and what environment variables need to be set to override defaults.

| Name | Description | Environment Variable | Default Value | 
|------|-------------|----------------------|---------------|
| port | the port on which the service will run | HTSGET_PORT | 3000 | 
| host | web service hostname. The JSON ticket returned by the server will reference other endpoints, using this hostname/base url to provide a complete url. | HTSGET_HOST | http://localhost:3000 | 

## Todo

* Configurable data sources (local files and urls) instead of hardcoding tabula muris dataset
* Correct handling of `referenceName=*` requests (unmapped, unplaced reads)
* Remove `samtools` dependency, as `samtools view` sometimes automatically fills placeholder values
* Implement `/variants/{id}` endpoint
* Implement `POST` request functionality 

## Issues

Bugs and issues can be submitted via the [Github Issue Tracker](https://github.com/ga4gh/htsget-refserver/issues)

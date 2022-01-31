[![Logo](https://www.ga4gh.org/wp-content/themes/ga4gh-theme/gfx/GA-logo-horizontal-tag-RGB.svg)](https://ga4gh.org)

# htsget Reference Server
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![Go Report](https://goreportcard.com/badge/github.com/ga4gh/htsget-refserver)](https://goreportcard.com/badge/github.com/ga4gh/htsget-refserver)
[![Travis (.org) branch](https://img.shields.io/travis/ga4gh/htsget-refserver/master.svg?style=flat-square)](https://travis-ci.org/ga4gh/htsget-refserver)
[![Coveralls github](https://img.shields.io/coveralls/github/ga4gh/htsget-refserver?style=flat-square)](https://coveralls.io/github/ga4gh/htsget-refserver?branch=master)

Reference server implementation of the htsget API protocol for securely streaming genomic data. For more information about htsget, see the [paper](https://academic.oup.com/bioinformatics/article/35/1/119/5040320) or [specification](http://samtools.github.io/hts-specs/htsget.html).

A GA4GH-hosted instance of this server is running at `https://htsget.ga4gh.org/`. To use, see the [OpenAPI documentation](https://htsget.ga4gh.org/docs/index.html).

## Quickstart - Docker

We suggest running the reference server as a docker container, as the image comes
pre-installed with all dependencies.

With docker installed, run:
```
docker image pull ga4gh/htsget-refserver:${TAG}
```
to pull the image, and:
```
docker container run -d -p 3000:3000 ga4gh/htsget-refserver:${TAG}
```
to spin up a containerized server. Custom config files can also be passed to the application by first mounting the directory containing the config, and specifying the path to config in the run command:
```
docker container run -d -p ${PORT}:${PORT} -v /directory/to/config:/usr/src/app/config ga4gh/htsget-refserver:${TAG} ./htsget-refserver -config /usr/src/app/config/config.json
```
Additional BAM/CRAM/VCF/BCF directories you wish to serve via htsget can also be mounted into the container. See the **Configuration** section below for instructions on how to serve custom datasets.

The full list of tags/versions is available on the [dockerhub repository page](https://hub.docker.com/repository/docker/ga4gh/htsget-refserver).

## Setup - Native

To run and/or develop the server natively on your OS, the following **dependencies** are required: 

* [Golang and language tools](https://golang.org/dl/) (tested on version 1.13) 
* [samtools](http://www.htslib.org/download/) (tested on version 1.9)
* [bcftools](http://www.htslib.org/download/) (tested on version 1.10.2)
* [htsget-refserver-utils](https://github.com/ga4gh/htsget-refserver-utils) (1.0.0+)

This project uses [Go modules](https://blog.golang.org/using-go-modules) to manage packages and dependencies.

With the above dependencies installed, run:
```
git clone https://github.com/ga4gh/htsget-refserver.git
cd htsget-refserver
```
to clone and enter the repository, and:
```
go build -o ./htsget-refserver ./cmd
```
to build the application binary. To start, run:
```
./htsget-refserver
```
A custom config file can also be specified with `-config`:
```
./htsget-refserver -config /path/to/config.json
```

## Configuration

The htsget web service can be configured with runtime parameters via a JSON config file, specified with `-config`. For example:
```
./htsget-refserver -config /path/to/config.json
```

Examples of valid JSON config files are available in this repository:

* [ga4gh instance config](./deployments/ga4gh/prod/config-server.json) - used to run the GA4GH-hosted instance at https://htsget.ga4gh.org
* [local development config](./deployments/ga4gh/prod/config-local.json) - used to run the local instance for development
* [integration tests config](./data/config/integration-tests.config.json) - used for integration testing on Travis CI builds
* [example 0 config](./data/config/example-0.config.json)
* [empty config](./data/config/example-empty.config.json)

In the JSON file, the root object must have a single "htsgetConfig" property, containing all sub-properties. ie:

```
{
    "htsgetConfig": {}
}
```

### Configuration - "props" object

Under the `htsgetConfig` property, the `props` object overrides application-wide settings. The following table indicates the attributes of `props` and what settings they affect.

| Name | Description |  Default Value | 
|------|-------------|----------------|
| port | the port on which the service will run | 3000 | 
| host | web service hostname. The JSON ticket returned by the server will reference other endpoints, using this hostname/base url to provide a complete url. | http://localhost:3000/ | 
| docsDir | path to static file directory containing server documentation (e.g. OpenAPI). the server will serve its contents at the `/docs/` endpoint | NONE |
| tempDir | writes temporary files used in request processing to this directory | . |
| logFile | writes application logs to this file | htsget-refserver.log |
| corsAllowedOrigins | CORS allow client from origins. Use comma to separate for multiple origins. | http://localhost |
| corsAllowedMethods | CORS allow methods. | GET, POST, PUT, DELETE, OPTIONS |
| corsAllowedHeaders | CORS allow headers.  | * |
| corsAllowCredentials | CORS allow credentials.  | false |
| corsMaxAge | CORS max age in seconds.  | 300 |
| awsAssumeRole | Turn on `awsAssumeRole` middleware. See **Private Bucket** section below. | false |

Example `props` object:

```
{
    "htsget": {
        "props": {
            "port": "80",
            "host": "https://htsget.ga4gh.org/",
            "tempdir": "/tmp/",
            "logfile": "/usr/src/app/htsget-refserver.log",
            "corsAllowedOrigins": "https://portal.ga4gh.org, http://intranet.ga4gh.org",
        }
    }
}
```

### Configuration - "reads" object

Under the `htsgetConfig` property, the `reads` object overrides settings for reads-related data and endpoints. The following properties can be set:

* `enabled` (boolean): if true, the server will set up reads-related routes (ie. `/reads/{id}`, `/reads/service-info`). True by default.
* `dataSourceRegistry` (object): allows the server to serve alignment data from multiple cloud or local storage sources by mapping request object id patterns to registered data sources. A single `sources` property contains an array of data sources. For each data source, the following properties are required:
    * `pattern` - a regex pattern that the `id` in `/reads/{id}` is matched against. If an `id` matches the pattern, the server will attempt to load data from the specified source. The pattern should make use of named capture group(s) to populate the path to the file.
    * `path` - the path template (either by url or local file path) to alignment files matching the pattern. The path must indicate how named capture groups in the pattern will populate the path to the file.
* `serviceInfo` (object): specify the attribute values returned in the Service Info response from `/reads/service-info`. Default attributes are supplied if not provided by config. Allows modification of the following properties from the Service Info specification:
    * `id`
    * `name`
    * `description`
    * `organization`
    * `contactUrl`
    * `documentationUrl`
    * `createdAt`
    * `updatedAt`
    * `environment`
    * `version`)

Example `reads` object:

```
{
    "htsgetConfig": {
        "reads": {
            "enabled": true,
            "dataSourceRegistry": {
                "sources": [
                    {
                        "pattern": "^tabulamuris\\.(?P<accession>10X.*)$",
                        "path": "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/{accession}_possorted_genome.bam"
                    },
                    {
                        "pattern": "^tabulamuris\\.(?P<accession>.*)$",
                        "path": "https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/{accession}.mus.Aligned.out.sorted.bam"
                    }
                ]
            }
            "serviceInfo": {
                "id": "demo.reads",
                "name": "htsget demo reads",
                "description": "serve alignment data via htsget",
                "organization": {
                    "name": "Example Org",
                    "url": "https://exampleorg.com"
                },
                "contactUrl": "mailto:nobody@exampleorg.com",
                "documentationUrl": "https://htsget.exampleorg.com/docs",
                "createdAt": "2021-01-01T09:00:00Z",
                "updatedAt": "2021-01-01T09:00:00Z",
                "environment": "test",
                "version": "1.0.0"
            }
        }
    }
}
```

### Configuration - "variants" object

Under the `htsgetConfig` property, the `variants` object overrides settings for variants-related data and endpoints. The following properties can be set:

* `enabled` (boolean): if true, the server will set up variants-related routes (ie. `/variants/{id}`, `/variants/service-info`). True by default.
* `dataSourceRegistry` (object): allows the server to serve variant data from multiple cloud or local storage sources by mapping request object id patterns to registered data sources. A single `sources` property contains an array of data sources. For each data source, the following properties are required:
    * `pattern` - a regex pattern that the `id` in `/variants/{id}` is matched against. If an `id` matches the pattern, the server will attempt to load data from the specified source. The pattern should make use of named capture group(s) to populate the path to the file.
    * `path` - the path template (either by url or local file path) to variant files matching the pattern. The path must indicate how named capture groups in the pattern will populate the path to the file.
* `serviceInfo` (object): specify the attribute values returned in the Service Info response from `/variants/service-info`. Default attributes are supplied if not provided by config. Allows modification of the following properties from the Service Info specification:
    * `id`
    * `name`
    * `description`
    * `organization`
    * `contactUrl`
    * `documentationUrl`
    * `createdAt`
    * `updatedAt`
    * `environment`
    * `version`)

Example `variants` object:

```
{
    "htsgetConfig": {
        "variants": {
            "enabled": true,
            "dataSourceRegistry": {
                "sources": [
                    {
                        "pattern": "^1000genomes\\.(?P<accession>.*)$",
                        "path": "https://ftp-trace.ncbi.nih.gov/1000genomes/ftp/phase1/analysis_results/integrated_call_sets/{accession}.vcf.gz"
                    }
                ]
            }
            "serviceInfo": {
                "id": "demo.variants",
                "name": "htsget demo variants",
                "description": "serve variant data via htsget",
                "organization": {
                    "name": "Example Org",
                    "url": "https://exampleorg.com"
                },
                "contactUrl": "mailto:nobody@exampleorg.com",
                "documentationUrl": "https://htsget.exampleorg.com/docs",
                "createdAt": "2021-01-01T09:00:00Z",
                "updatedAt": "2021-01-01T09:00:00Z",
                "environment": "test",
                "version": "1.0.0"
            }
        }
    }
}
```

## Private Bucket

- Turn on `awsAssumeRole` [middleware](https://github.com/go-chi/chi#middleware-handlers) request interceptor to support AWS [Assume Role](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html) temporary security credentials loading to access S3 private bucket.

- When it is not configured, the default `awsAssumeRole` set to `false` such that execution environment know how to access S3 private bucket through AWS [standard mechanism](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/). In that case, see [htslib AWS S3 plugin](http://www.htslib.org/doc/htslib-s3-plugin.html) for credentials loading requirement.

Say, you have data in private bucket as follows:
```
s3://my-primary-data-prod/Project/PID00115/WGS/PID00115-final.bam
```

Example configuration:
```
{
  "htsgetConfig": {
    "props": {
      ...
      "awsAssumeRole": true
    },
    "reads": {
      ...
      "dataSourceRegistry": {
        "sources": [
          ...
          {
            "pattern": "^my-primary-data(?P<accession>.*)$",
            "path": "s3://my-primary-data{accession}"
          }
        ]
      },
      "serviceInfo": {
        ...
      }
    },
    "variants": {
      ...
      "dataSourceRegistry": {
        "sources": [
          ...
          {
            "pattern": "^my-primary-data(?P<accession>.*)$",
            "path": "s3://my-primary-data{accession}"
          }
        ]
      },
      "serviceInfo": {
        ...
      }
    }
  }
}
```

Then you can call Htsget as follows:
```
curl -s http://localhost:3000/reads/my-primary-data-prod/Project/PID00115/WGS/PID00115-final.bam | jq
```

## Testing

To execute unit and end-to-end tests on the entire package, run:
```
go test ./... -coverprofile=cp.out
```
The go coverage report will be available at `./cp.out`. To execute tests for a specific package (for example the `htsrequest` package) run:
```
go test ./internal/htsrequest -coverprofile=cp.out
```

## Changelog

**v1.5.0**
* Supports configurable CORS headers

**v1.4.0**
* Server supports **experimental** `POST` method for endpoint `/reads/{id}`. Multiple
genomic regions can be requested in a single request. See
[hts-specs PR #285](https://github.com/samtools/hts-specs/pull/285) for more info.

**v1.3.0**
* Server supports reads and/or variants `service-info` endpoints. The attributes of the `service-info` response can be specified via the config file independently for each datatype 

**v1.2.0**

* Server supports htsget `/variants/{id}` endpoint, streams VCFs via htsget protocol
using bcftools dependency

**v1.1.0**

* Added support for configurable data sources via a data source registry specified
in config file
* server can stream reads data via htsget protocol from any **url** or **local file** specified via config 

**v1.0.0**

* Initial release

## Roadmap

* Implement `POST` request functionality 

## Maintainers

* Jeremy Adams (jb-adams) [jeremy.adams@ga4gh.org](mailto:jeremy.adams@ga4gh.org)
* David Liu (xngln)

## Issues

Bugs and issues can be submitted via the [Github Issue Tracker](https://github.com/ga4gh/htsget-refserver/issues)

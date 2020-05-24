# HTSget Reference Server
This is a reference implementation of the version 1.2.0 HTSget API protocol for securely streaming genomic data. Click [here](https://academic.oup.com/bioinformatics/article/35/1/119/5040320) for a high level overview of the standard, and [here](https://github.com/samtools/hts-specs/blob/master/htsget.md) to view the specification itself. 

## Setup
- Install [Golang(v1.13) and language tools](https:/ /golang.org/dl/). With Homebrew, `$ brew install go`

This project uses [Go modules](https://blog.golang.org/using-go-modules) to manage packages and dependencies.

`$ go build & ./htsget-refserver` from root directory to start the server on port 3000.

## Usage
The API is defined at https://github.com/samtools/hts-specs/blob/master/htsget.md. 
This server is deployed at http://htsget.ga4gh.org/.


# HTSget Reference Server #
This is a reference server implementation of the version 1.2.0 HTSget API protocol for securely streaming genomic data. Click [here](https://academic.oup.com/bioinformatics/article/35/1/119/5040320) for a high level overview of the standard, and [here](https://github.com/samtools/hts-specs/blob/master/htsget.md) to view the specification itself. 

  - design
  - API usage
  - setup instructions

## Setup
- Install [Golang(v1.13) and language tools](https://golang.org/dl/)

Note: this project uses the new [Go modules](https://blog.golang.org/using-go-modules) to manage packages and dependencies.

`$ go run main.go` to start the server on port 3000.

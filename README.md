# edgr | Makes SEC filings not terrible

[![Go Report Card](https://goreportcard.com/badge/github.com/piquette/edgr)](https://goreportcard.com/badge/github.com/piquette/edgr)
[![Build Status](https://travis-ci.org/piquette/edgr.svg?branch=master)](https://travis-ci.org/piquette/edgr)
[![GoDoc](https://godoc.org/github.com/piquette/edgr?status.svg)](https://godoc.org/github.com/piquette/qtrn)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This project consists of:
* `edgr` - a cli tool that can populate a postgres db with SEC filings for use in data analysis projects
* `github.com/edgr/core` a go module that requests SEC filers/filings and can be used in other go projects


## edgr CLI
The `edgr` tool makes jumpstarting your data analysis projects that much easier by abstacting away a lot of annoying SEC EDGAR data gathering and parsing. `edgr` can download, parse full text, and persist large quantities of SEC filings in a pre-configured postgres db.

### Installation
`edgr` supports macOS, linux, and windows systems.

#### Brew
The preferred way to install, run the following commands in the terminal after installing homebrew
```sh
brew tap piquette/edgr
brew install edgr
```
#### Download the binary
Distributions for your specific OS can be found on the releases page for this repo.

#### Build from source
Clone this repository and run `make build`

### Usage
`edgr` is a standard cli executable. A running/reachable postgres instance is required for proper usage. You can easily run one locally through docker.
```sh
docker run -p 5432:5432 -d postgres:9.6
```
The default parameters for `edgr` commands will be able to access this instance immediately.
Next, we need to set up the proper db tables to store data. In the terminal, run 
```sh
edgr init
```
Finally, we need to insert some desired filers into the db so we can pull some filings.
```sh
# add the filer.
edgr filers init --symbol AAPL
# download 8-K filings for apple inc.
edgr get --symbol AAPL --form 8-K
```
Now you can access the raw text and meta-data associated with those filings in the database as you normally would, through other analytical tools or a command-line pg client. Run `edgr help` for more information regarding commands.

## `github.com/edgr/core` module

### Installation
Supports Go v1.13+ and modules. Run the usual
```sh
go get -u github.com/piquette/edgr
```

### Usage
There are 3 easy functions
```go
// GetPublicCompanies returns a list of public companies.
func GetPublicCompanies() ([]Company, error) {...}

// GetFiler gets a single filer from the SEC website based on symbol.
func GetFiler(symbol string) (filer *model.Filer, err error) {...}

// GetFilings gets a list of filings for a single CIK.
func GetFilings(cik, formtype, stoptime string) (filings []SECFiling, err error) {...}
```

package main

import "github.com/urfave/cli"

func buildFilersFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:     "all",
			Usage:    "Fetches and stores records for all possible filers",
			EnvVar:   "FILER_ALL",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "symbol",
			Usage:    "Specifies a single symbol to fetch and store a filer record for",
			EnvVar:   "FILER_SYMBOL",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "sic",
			Usage:    "Specifies an industry group to fetch and store filer records for",
			EnvVar:   "FILER_SIC_GROUP",
			FilePath: "EDGRFILE",
		},
	}
}

func buildGetFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "symbol",
			Usage:    "Specify a single symbol to fetch and store filings for",
			EnvVar:   "GET_SYMBOL",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "form",
			Usage:    "Specify a form type to fetch and store filings for",
			EnvVar:   "GET_SYMBOL",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "stop",
			Usage:    "Specify a date to stop retrieving filings records, format is YYYY-MM-DD",
			Value:    "2019-06-01",
			EnvVar:   "GET_TIME",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "sic",
			Usage:    "Specify an industry group to fetch and store filings for",
			EnvVar:   "GET_SIC_GROUP",
			FilePath: "EDGRFILE",
		},
	}
}

func buildGlobalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "pg-addr",
			Value:    "localhost:5432",
			Usage:    "PostgreSQL ~address~",
			EnvVar:   "PG_ADDR",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "pg-user",
			Value:    "postgres",
			Usage:    "PostgreSQL ~username~",
			EnvVar:   "PG_USER",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "pg-pass",
			Value:    "postgres",
			Usage:    "PostgreSQL ~password~",
			EnvVar:   "PG_PASS",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "pg-db",
			Value:    "postgres",
			Usage:    "PostgreSQL ~database~",
			EnvVar:   "PG_DB",
			FilePath: "EDGRFILE",
		},
	}
}

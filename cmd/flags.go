package main

import "github.com/urfave/cli"

func buildGetFlags() []cli.Flag {
	// 	LetterStart string `flag:"letter" env:"LETTER" usage:"LETTER"`
	// 	SymbolStart string `flag:"start" env:"START" usage:"START"`
	// 	FillMode    string `flag:"fill-mode" env:"FILL_MODE" usage:"Fill mode"`
	// 	StopDate    string `flag:"stop-date" env:"STOP_DATE" usage:"Stop date YYYY-MM-DD"`
	// 	LetterStart:   "A",
	// 	SymbolStart:   "",
	// 	FillMode:      "alpha",
	// 	StopDate:      "2019-06-01",
	return []cli.Flag{
		cli.StringFlag{
			Name:     "symbol",
			Usage:    "Specifies a symbol to get filings for",
			EnvVar:   "GET_SYMBOL",
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

package model

import "time"

// Filing represents a sec filing after a successful parsing of an Entry.
type Filing struct {
	ID string

	Created           time.Time
	Updated           time.Time
	Filer             string
	Accession         string
	CIK               string
	EdgarURL          string
	EdgarTime         time.Time
	FilerRelation     string
	FormType          string
	DocumentCount     int64
	TotalSizeEstimate string
	Symbol            string
	AllSymbols        []string
	AllCIKs           []string
}

// Document represent a single document in a filing.
type Document struct {
	ID string

	Created      time.Time
	Sequence     int64
	DocType      string
	EdgarURL     string
	Description  string
	Body         string
	SizeEstimate string
	FilingID     string
}

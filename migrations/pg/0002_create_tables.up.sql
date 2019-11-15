-- filings table: stores filings.
CREATE TABLE filings (
	id TEXT PRIMARY KEY DEFAULT generate_uid(24),

  created TIMESTAMP NOT NULL DEFAULT now(),
  updated TIMESTAMP NOT NULL DEFAULT now(),
	filer TEXT NOT NULL,
  accession TEXT UNIQUE NOT NULL,
  cik TEXT NOT NULL,
	all_ciks TEXT[] NOT NULL,
	edgar_url TEXT NOT NULL,
  edgar_time TIMESTAMP NOT NULL,
	relation TEXT NOT NULL,
	form_type TEXT NOT NULL,
	doc_count BIGINT,
	total_size_est TEXT,
  symbol TEXT,
	all_symbols TEXT[]
);

-- documents table: stores each document of a filing.
CREATE TABLE documents (
	id TEXT PRIMARY KEY DEFAULT generate_uid(24),

  created TIMESTAMP NOT NULL DEFAULT now(),
	seq BIGINT NOT NULL,
	doc_type TEXT NOT NULL,
	edgar_url TEXT NOT NULL,
	doc_desc TEXT NOT NULL,
	data BYTEA NOT NULL,
	size_est TEXT,

	filing_id TEXT REFERENCES filings
);

-- forms table: stores forms.
CREATE TABLE forms (
	id TEXT PRIMARY KEY DEFAULT generate_uid(4),

  form_type TEXT NOT NULL,
  form_desc TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT false
);

-- filers table: stores filers.
CREATE TABLE filers (
	id TEXT PRIMARY KEY DEFAULT generate_uid(8),

  symbol TEXT,
  cik TEXT NOT NULL,
  name TEXT NOT NULL,
  sic TEXT,
  sic_desc TEXT
);

--ALTER TABLE filers ALTER COLUMN name DROP NOT NULL;
ALTER TABLE filers ADD CONSTRAINT u_cik UNIQUE (cik);

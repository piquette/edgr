-- filings table indexes
CREATE INDEX filings_created_idx ON filings (created DESC);
CREATE INDEX filings_edgar_time_idx ON filings (edgar_time DESC);
CREATE INDEX filings_symbol_idx ON filings (symbol);
CREATE INDEX filings_cik_idx ON filings (cik);
CREATE INDEX filings_accession_idx ON filings (accession);
CREATE INDEX filings_form_idx ON filings (form_type);
-- documents table indexes
CREATE INDEX documents_type_idx ON documents (doc_type);
CREATE INDEX documents_filing_id_idx ON documents (filing_id);

-- forms table indexes
CREATE INDEX forms_id_idx ON forms (id);
CREATE INDEX forms_type_idx ON forms (form_type);

-- filers table indexes
CREATE INDEX filers_id_idx ON filers (id);
CREATE INDEX filers_symbol_idx ON filers (symbol);
CREATE INDEX filers_cik_idx ON filers (cik);
CREATE INDEX filers_name_idx ON filers (name);
CREATE INDEX filers_search_idx ON filers using gin ((name || ' ' || cik || ' ' || symbol) gin_trgm_ops);

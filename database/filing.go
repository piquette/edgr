package database

import (
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/piquette/edgr/core/model"
)

// Filing is an ORM object for the `filings` table.
type Filing struct {
	ID string

	Created           time.Time
	Updated           time.Time
	Filer             string
	Accession         string
	CIK               string
	EdgarURL          string    `sql:"edgar_url"`
	EdgarTime         time.Time `sql:"edgar_time"`
	FilerRelation     string    `sql:"relation"`
	FormType          string    `sql:"form_type"`
	DocumentCount     int64     `sql:"doc_count"`
	TotalSizeEstimate string    `sql:"total_size_est"`
	Symbol            string
	AllSymbols        []string `sql:"all_symbols,array"`
	AllCIKs           []string `sql:"all_ciks,array"`
}

// QueryParams holds api query parameters.
type QueryParams struct {
	Term  string
	Page  int
	Size  int
	Sort  int
	Form  string
	Start string
	End   string
}

// QueryResult holds db query result.
type QueryResult struct {
	Filings    []*model.Filing
	TotalCount int
}

// FilingDao provides access to filing storage.
type FilingDao interface {
	StrictExists(cik, accession, relation string) (bool, error)
	Add(filing *model.Filing) (created *model.Filing, existed bool, err error)
	Update(filing *model.Filing) (bool, error)
	AddDocuments(filing *model.Filing, docs []*model.Document) error
	Get(accession string) (exists bool, filing *model.Filing, err error)

	GetByID(filingID string) (exists bool, filing *model.Filing, err error)
	GetBySymbol(params *QueryParams) (exists bool, result *QueryResult, err error)
	GetByCIK(params *QueryParams) (exists bool, result *QueryResult, err error)
	GetByFiler(params *QueryParams) (exists bool, result *QueryResult, err error)
	GetAll(params *QueryParams) (exists bool, result *QueryResult, err error)
}

// FilingDaoImpl implements FilingDao.
type FilingDaoImpl struct {
	db *Handle
}

// StrictExists checks strictly.
func (dao *FilingDaoImpl) StrictExists(cik, accession, relation string) (bool, error) {
	f := &Filing{}
	return dao.db.Model(f).Where("cik = ? and accession = ? and relation = ?", cik, accession, relation).Exists()
}

// Add adds a new filing.
func (dao *FilingDaoImpl) Add(filing *model.Filing) (created *model.Filing, existed bool, err error) {
	f := filingAsRecord(filing)
	res, err := dao.db.Model(f).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return nil, false, err
	}
	if res.RowsAffected() > 0 {
		return f.export(), false, nil
	}
	return f.export(), true, nil
}

// Update updates an existing filing.
func (dao *FilingDaoImpl) Update(filing *model.Filing) (bool, error) {
	f := filingAsRecord(filing)
	f.Updated = time.Now()
	res, dberr := dao.db.Model(f).Column("updated", "filer", "cik", "symbol", "relation", "all_symbols", "all_ciks").WherePK().Update()
	if dberr != nil {
		return false, dberr
	}
	return res.RowsAffected() != 0, nil
}

// AddDocuments adds docs to an existing filing.
func (dao *FilingDaoImpl) AddDocuments(filing *model.Filing, docs []*model.Document) error {
	if filing == nil || filing.ID == "" {
		return fmt.Errorf("no filing to add docs to")
	}
	if len(docs) == 0 {
		return fmt.Errorf("no docs")
	}
	f := filingAsRecord(filing)
	var totalSize int64
	records := []*Document{}
	for _, doc := range docs {
		rec := documentPartial(doc)
		rec.FilingID = filing.ID
		rec.Data = []byte(doc.Body)
		size := int64(len(rec.Data))
		rec.SizeEstimate = byteCountDecimal(size)
		totalSize += size
		records = append(records, rec)
	}
	f.DocumentCount = int64(len(records))
	f.TotalSizeEstimate = byteCountDecimal(totalSize)
	f.Updated = time.Now()

	tx, err := dao.db.Begin()
	if err != nil {
		return err
	}

	// Insert docs.
	err = dao.db.Insert(&records)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Model(f).Column("doc_count", "total_size_est", "updated").WherePK().Update()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Get queries a single filing.
func (dao *FilingDaoImpl) Get(accession string) (exists bool, filing *model.Filing, err error) {
	f := &Filing{}
	err = dao.db.Model(f).Where("accession = ?", accession).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return false, nil, nil
		}
		return
	}
	return true, f.export(), nil
}

// GetByID gets a filing by id.
func (dao *FilingDaoImpl) GetByID(filingID string) (exists bool, filing *model.Filing, err error) {
	f := &Filing{ID: filingID}
	err = dao.db.Select(f)
	if err != nil {
		if err == pg.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, f.export(), nil
}

// GetBySymbol gets filings by symbol.
func (dao *FilingDaoImpl) GetBySymbol(params *QueryParams) (exists bool, result *QueryResult, err error) {

	// Validate params.
	if params == nil {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Term == "" {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Size <= 0 || params.Size > 50 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Page < 0 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Sort != -1 && params.Sort != 1 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	// Validate times.
	if len(params.Start) != 0 && len(params.Start) != 10 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if len(params.End) != 0 && len(params.End) != 10 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if len(params.Form) > 7 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}

	var results []Filing
	q := dao.db.Model(&results).
		Where("? =ANY(all_symbols)", params.Term)

	if params.Start != "" {
		q = q.Where("edgar_time >= ?", params.Start)
	}
	if params.End != "" {
		q = q.Where("edgar_time <= ?", params.End)
	}
	if params.Form != "" {
		q = q.Where("form_type = ?", params.Form)
	}

	order := "desc"
	if params.Sort == -1 {
		order = "asc"
	}
	q = q.Limit(params.Size).
		Offset(params.Page * params.Size).
		Order("edgar_time " + order)

	totalCount, err := q.SelectAndCount()
	if err != nil {
		return false, nil, err
	}
	if totalCount == 0 {
		return false, nil, nil
	}

	// map.
	var filings []*model.Filing
	for _, result := range results {
		filings = append(filings, result.export())
	}

	return true, &QueryResult{TotalCount: totalCount, Filings: filings}, nil
}

// GetByCIK gets filings by cik.
func (dao *FilingDaoImpl) GetByCIK(params *QueryParams) (exists bool, result *QueryResult, err error) {

	// Validate params.
	if params == nil {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Term == "" {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Size <= 0 || params.Size > 50 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Page < 0 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Sort != -1 && params.Sort != 1 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	// Validate times.
	if len(params.Start) != 0 && len(params.Start) != 10 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if len(params.End) != 0 && len(params.End) != 10 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if len(params.Form) > 7 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}

	var results []Filing
	q := dao.db.Model(&results).
		Where("? =ANY(all_ciks)", params.Term)

	if params.Start != "" {
		q = q.Where("edgar_time >= ?", params.Start)
	}
	if params.End != "" {
		q = q.Where("edgar_time <= ?", params.End)
	}
	if params.Form != "" {
		q = q.Where("form_type = ?", params.Form)
	}

	order := "desc"
	if params.Sort == -1 {
		order = "asc"
	}
	q = q.Limit(params.Size).
		Offset(params.Page * params.Size).
		Order("edgar_time " + order)

	totalCount, err := q.SelectAndCount()
	if err != nil {
		return false, nil, err
	}
	if totalCount == 0 {
		return false, nil, nil
	}

	// map.
	var filings []*model.Filing
	for _, result := range results {
		filings = append(filings, result.export())
	}

	return true, &QueryResult{TotalCount: totalCount, Filings: filings}, nil
}

// GetByFiler gets filings by filer name.
func (dao *FilingDaoImpl) GetByFiler(params *QueryParams) (exists bool, result *QueryResult, err error) {

	// Validate params.
	if params == nil {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Term == "" {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Size <= 0 || params.Size > 50 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Page < 0 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Sort != -1 && params.Sort != 1 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	// Validate times.
	if len(params.Start) != 0 && len(params.Start) != 10 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if len(params.End) != 0 && len(params.End) != 10 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if len(params.Form) > 7 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}

	var results []Filing
	q := dao.db.Model(&results).
		Where("filer ILIKE ?", params.Term+"%")

	if params.Start != "" {
		q = q.Where("edgar_time >= ?", params.Start)
	}
	if params.End != "" {
		q = q.Where("edgar_time <= ?", params.End)
	}
	if params.Form != "" {
		q = q.Where("form_type = ?", params.Form)
	}

	order := "desc"
	if params.Sort == -1 {
		order = "asc"
	}
	q = q.Limit(params.Size).
		Offset(params.Page * params.Size).
		Order("edgar_time " + order)

	totalCount, err := q.SelectAndCount()
	if err != nil {
		return false, nil, err
	}
	if totalCount == 0 {
		return false, nil, nil
	}

	// map.
	var filings []*model.Filing
	for _, result := range results {
		filings = append(filings, result.export())
	}

	return true, &QueryResult{TotalCount: totalCount, Filings: filings}, nil
}

// GetAll gets filings.
func (dao *FilingDaoImpl) GetAll(params *QueryParams) (exists bool, result *QueryResult, err error) {

	// Validate params.
	if params == nil {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Size <= 0 || params.Size > 50 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Page < 0 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if params.Sort != -1 && params.Sort != 1 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	// Validate times.
	if len(params.Start) != 0 && len(params.Start) != 10 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if len(params.End) != 0 && len(params.End) != 10 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}
	if len(params.Form) > 7 {
		return false, nil, fmt.Errorf("incorrect query parameters")
	}

	var results []Filing
	q := dao.db.Model(&results)

	if params.Start != "" {
		q = q.Where("edgar_time >= ?", params.Start)
	}
	if params.End != "" {
		q = q.Where("edgar_time <= ?", params.End)
	}
	if params.Form != "" {
		q = q.Where("form_type = ?", params.Form)
	}

	order := "desc"
	if params.Sort == -1 {
		order = "asc"
	}
	q = q.Limit(params.Size).
		Offset(params.Page * params.Size).
		Order("edgar_time " + order)

	totalCount, err := q.SelectAndCount()
	if err != nil {
		return false, nil, err
	}
	if totalCount == 0 {
		return false, nil, nil
	}

	// map.
	var filings []*model.Filing
	for _, result := range results {
		filings = append(filings, result.export())
	}

	return true, &QueryResult{TotalCount: totalCount, Filings: filings}, nil
}

func documentPartial(d *model.Document) *Document {
	return &Document{
		ID:          d.ID,
		Created:     d.Created,
		Sequence:    d.Sequence,
		DocType:     d.DocType,
		EdgarURL:    d.EdgarURL,
		Description: d.Description,
	}
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func filingAsRecord(f *model.Filing) *Filing {
	return &Filing{
		ID:                f.ID,
		Created:           f.Created,
		Updated:           f.Updated,
		Filer:             f.Filer,
		Accession:         f.Accession,
		CIK:               f.CIK,
		EdgarURL:          f.EdgarURL,
		EdgarTime:         f.EdgarTime,
		FilerRelation:     f.FilerRelation,
		FormType:          f.FormType,
		DocumentCount:     f.DocumentCount,
		TotalSizeEstimate: f.TotalSizeEstimate,
		Symbol:            f.Symbol,
		AllSymbols:        f.AllSymbols,
		AllCIKs:           f.AllCIKs,
	}
}
func (f *Filing) export() *model.Filing {
	return &model.Filing{
		ID:                f.ID,
		Created:           f.Created,
		Updated:           f.Updated,
		Filer:             f.Filer,
		Accession:         f.Accession,
		CIK:               f.CIK,
		EdgarURL:          f.EdgarURL,
		EdgarTime:         f.EdgarTime,
		FilerRelation:     f.FilerRelation,
		FormType:          f.FormType,
		DocumentCount:     f.DocumentCount,
		TotalSizeEstimate: f.TotalSizeEstimate,
		Symbol:            f.Symbol,
		AllSymbols:        f.AllSymbols,
		AllCIKs:           f.AllCIKs,
	}
}

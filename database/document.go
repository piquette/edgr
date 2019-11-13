package database

import (
	"time"

	pg "github.com/go-pg/pg"
	"github.com/piquette/edgr/core/model"
)

// Document is an ORM object for the `documents` table.
type Document struct {
	ID string

	Created      time.Time
	Sequence     int64  `sql:"seq"`
	DocType      string `sql:"doc_type"`
	EdgarURL     string `sql:"edgar_url"`
	Description  string `sql:"doc_desc"`
	Data         []byte `sql:"data"`
	SizeEstimate string `sql:"size_est"`
	FilingID     string `sql:"filing_id"`
}

// DocumentDao provides access to document storage.
type DocumentDao interface {
	GetByFilingID(filingID string) (exists bool, docs []*model.Document, err error)
	GetIDs(filingID string) (exists bool, ids []string, err error)
	Get(documentID string) (exists bool, doc *model.Document, err error)
}

// DocumentDaoImpl implements DocumentDao.
type DocumentDaoImpl struct {
	db *Handle
}

// GetByFilingID gets docs by filing id.
func (dao *DocumentDaoImpl) GetByFilingID(filingID string) (exists bool, docs []*model.Document, err error) {
	var results []Document
	err = dao.db.Model(&results).Where("filing_id = ?", filingID).Order("seq ASC").Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, err
	}

	docs = []*model.Document{}
	for _, result := range results {
		docs = append(docs, result.export())
	}

	return true, docs, nil
}

// GetIDs gets doc ids by filing id.
func (dao *DocumentDaoImpl) GetIDs(filingID string) (exists bool, ids []string, err error) {
	var results []Document
	err = dao.db.Model(&results).Where("filing_id = ?", filingID).
		Column("id").Order("seq ASC").Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, err
	}

	ids = []string{}
	for _, result := range results {
		ids = append(ids, result.ID)
	}

	return true, ids, nil
}

// Get gets a doc by id.
func (dao *DocumentDaoImpl) Get(documentID string) (exists bool, doc *model.Document, err error) {
	d := &Document{ID: documentID}
	err = dao.db.Select(d)
	if err != nil {
		if err == pg.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, d.export(), nil
}

func (d *Document) export() *model.Document {
	return &model.Document{
		ID:           d.ID,
		Created:      d.Created,
		Sequence:     d.Sequence,
		DocType:      d.DocType,
		EdgarURL:     d.EdgarURL,
		Description:  d.Description,
		Body:         string(d.Data),
		SizeEstimate: d.SizeEstimate,
		FilingID:     d.FilingID,
	}
}

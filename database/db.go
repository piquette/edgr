package database

import (
	pg "github.com/go-pg/pg"
)

// Handle is a database handle.
type Handle struct {
	*pg.DB
	opts Options
}

// Options is a set of options for connecting to PostgreSQL.
type Options struct {
	Addr     string
	User     string
	Password string
	Database string
}

// Open opens a new database connection.
func Open(opts Options) *Handle {
	db := pg.Connect(&pg.Options{
		Addr:     opts.Addr,
		User:     opts.User,
		Password: opts.Password,
		Database: opts.Database,
	})

	return &Handle{db, opts}
}

// NewFilerDao returns a filer dao.
func (h *Handle) NewFilerDao() *FilerDaoImpl {
	return &FilerDaoImpl{db: h}
}

// NewFormDao returns a form dao.
func (h *Handle) NewFormDao() *FormDaoImpl {
	return &FormDaoImpl{db: h}
}

// NewFilingDao returns a filing dao.
func (h *Handle) NewFilingDao() *FilingDaoImpl {
	return &FilingDaoImpl{db: h}
}

// NewDocumentDao returns a document dao.
func (h *Handle) NewDocumentDao() *DocumentDaoImpl {
	return &DocumentDaoImpl{db: h}
}

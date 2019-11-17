package database

import (
	pg "github.com/go-pg/pg"
	"github.com/piquette/edgr/core/model"
)

type (
	// FilerDao abstracts access to sec filers.
	FilerDao interface {
		Search(condition string) (bool, []*model.Filer, error)
		GetSet(term string) (bool, []*model.Filer, error)
		Put(filer *model.Filer) (bool, error)
	}
	// Filer is an ORM object for the `filers` table.
	Filer struct {
		ID             string
		Symbol         string
		CIK            string
		Name           string `sql:"name"`
		SIC            string
		SICDescription string `sql:"sic_desc"`
	}
	// FilerDaoImpl implements the FilerDao.
	FilerDaoImpl struct {
		db *Handle
	}
)

// Search attempts to match a list of filers from a query term.
func (dao *FilerDaoImpl) Search(condition string) (found bool, filers []*model.Filer, err error) {
	var results []Filer
	condition = "%" + condition + "%"

	err = dao.db.Model(&results).
		Where("name || ' ' || cik || ' ' || symbol ILIKE ?", condition).
		Limit(10).
		Order("symbol ASC").
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return false, filers, nil
		}
		return false, filers, err
	}

	for _, result := range results {
		filers = append(filers, result.export())
	}

	return true, filers, nil
}

// GetSet attempts to retrieve a set of filers by symbol.
func (dao *FilerDaoImpl) GetSet(term string) (found bool, filers []*model.Filer, err error) {
	var results []Filer

	err = dao.db.Model(&results).
		Where("symbol LIKE ?", term).
		Order("symbol ASC").
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return false, filers, nil
		}
		return false, filers, err
	}

	for _, result := range results {
		filers = append(filers, result.export())
	}

	return true, filers, nil
}

// GetSetBySIC get filers by sic.
func (dao *FilerDaoImpl) GetSetBySIC(term string) (found bool, filers []*model.Filer, err error) {
	var results []Filer

	err = dao.db.Model(&results).
		Where("sic LIKE ?", term).
		Order("symbol ASC").
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return false, filers, nil
		}
		return false, filers, err
	}

	for _, result := range results {
		filers = append(filers, result.export())
	}

	return true, filers, nil
}

// Put puts a potential new filer and returns whether or not it already existed..
func (dao *FilerDaoImpl) Put(filer *model.Filer) (bool, error) {
	f := filerAsRecord(filer)

	res, err := dao.db.Model(f).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return false, err
	}
	if res.RowsAffected() > 0 {
		return false, nil
	}
	return true, nil
}

func filerAsRecord(f *model.Filer) *Filer {
	return &Filer{
		ID:             f.ID,
		CIK:            f.CIK,
		Name:           f.Name,
		Symbol:         f.Symbol,
		SIC:            f.SIC,
		SICDescription: f.SICDescription,
	}
}

func (f *Filer) export() *model.Filer {
	return &model.Filer{
		ID:             f.ID,
		CIK:            f.CIK,
		Name:           f.Name,
		Symbol:         f.Symbol,
		SIC:            f.SIC,
		SICDescription: f.SICDescription,
	}
}

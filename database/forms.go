package database

import "github.com/piquette/edgr/core/model"

type (
	// FormDao abstracts access to sec forms.
	FormDao interface {
		GetActiveForms() ([]*model.Form, error)
	}
	// Form is an ORM object for the `forms` table.
	Form struct {
		ID          string
		FormType    string `sql:"form_type"`
		Description string `sql:"form_desc"`
		Active      bool
	}
	// FormDaoImpl implements the FormDao.
	FormDaoImpl struct {
		db *Handle
	}
)

// GetActiveForms gets all the active forms.
func (dao *FormDaoImpl) GetActiveForms() (forms []*model.Form, err error) {
	var results []Form

	err = dao.db.Model(&results).
		Where("active = true").
		Order("form_type ASC").
		Select()

	for _, result := range results {
		forms = append(forms, mapFormResult(&result))
	}

	return
}

func mapFormResult(f *Form) *model.Form {
	return &model.Form{
		ID:          f.ID,
		FormType:    f.FormType,
		Description: f.Description,
		Active:      f.Active,
	}
}

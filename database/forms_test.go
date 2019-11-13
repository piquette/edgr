package database

import "github.com/stretchr/testify/assert"

func (s *APISuite) TestFormDao() {
	// validate impl.
	dao := s.db.NewFormDao()
	assert.Implements(s.T(), (*FormDao)(nil), dao)

	// Smoke test.
	{
		forms, err := dao.GetActiveForms()
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), 16, len(forms))
		for _, form := range forms {
			assert.True(s.T(), form.Active)
		}
	}
}

func (s *EdgrSuite) TestFormDao() {
	// validate impl.
	dao := s.db.NewFormDao()
	assert.Implements(s.T(), (*FormDao)(nil), dao)

	// Smoke test.
	{
		forms, err := dao.GetActiveForms()
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), 16, len(forms))
		for _, form := range forms {
			assert.True(s.T(), form.Active)
		}
	}
}

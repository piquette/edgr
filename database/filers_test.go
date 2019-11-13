package database

import (
	"github.com/piquette/edgr/core/model"
	"github.com/stretchr/testify/assert"
)

func (s *APISuite) TestFilerDao() {
	// up
	s.exec("./testdata/filers.up.sql")

	// validate impl.
	dao := s.db.NewFilerDao()
	assert.Implements(s.T(), (*FilerDao)(nil), dao)

	// Symbol.
	{
		found, filers, err := dao.Search("TWT")
		assert.Nil(s.T(), err)
		assert.True(s.T(), found)
		assert.Equal(s.T(), 1, len(filers))
		assert.Equal(s.T(), "TWTR", filers[0].Symbol)
	}
	// Name.
	{
		found, filers, err := dao.Search("twi")
		assert.Nil(s.T(), err)
		assert.True(s.T(), found)
		assert.Equal(s.T(), 1, len(filers))
		assert.Equal(s.T(), "TWTR", filers[0].Symbol)
	}
	// CIK.
	{
		found, filers, err := dao.Search("0965458")
		assert.Nil(s.T(), err)
		assert.True(s.T(), found)
		assert.Equal(s.T(), 6, len(filers))
		// Abc by symbol, so apple comes first.
		assert.Equal(s.T(), "AAPL", filers[0].Symbol)
	}
	// CIK 2.
	{
		found, filers, err := dao.Search("09654584")
		assert.Nil(s.T(), err)
		assert.True(s.T(), found)
		assert.Equal(s.T(), 2, len(filers))
		// Abc by symbol, so ibm comes first.
		assert.Equal(s.T(), "IBM", filers[0].Symbol)
	}

	// down
	s.exec("./testdata/filers.down.sql")
}

func (s *APISuite) TestFilerPut() {
	dao := s.db.NewFilerDao()

	// Put empty.
	{
		aapl := &model.Filer{Name: "AAPL Corp"}
		existed, err := dao.Put(aapl)
		assert.NotNil(s.T(), err)
		assert.False(s.T(), existed)
	}
	// Put with CIK.
	{
		aapl := &model.Filer{Name: "AAPL Corp", CIK: "78904"}
		existed, err := dao.Put(aapl)
		assert.Nil(s.T(), err)
		assert.False(s.T(), existed)
	}
	// Put duplicate.
	{
		twtr := &model.Filer{Name: "TWTR Corp", CIK: "1234"}
		dupe := &model.Filer{Name: "TWTW Corp", CIK: "1234"}
		existed, err := dao.Put(twtr)
		assert.Nil(s.T(), err)
		assert.False(s.T(), existed)

		existed, err = dao.Put(dupe)
		assert.Nil(s.T(), err)
		assert.True(s.T(), existed)
	}
}

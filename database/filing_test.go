package database

import (
	"testing"
	"time"

	"github.com/piquette/edgr/core/model"
	"github.com/stretchr/testify/assert"
)

func (s *EdgrSuite) TestFilingDao() {
	// validate impl.
	dao := s.db.NewFilingDao()
	assert.Implements(s.T(), (*FilingDao)(nil), dao)
}

func (s *EdgrSuite) TestStrictExists() {
	dao := s.db.NewFilingDao()
	{
		candidate := &model.Filing{
			Filer:         "ABC Corp",
			Accession:     "123",
			CIK:           "123",
			AllCIKs:       []string{"123"},
			EdgarURL:      "example.com",
			EdgarTime:     time.Now(),
			FilerRelation: "Issuer",
			FormType:      "8-K",
		}
		filing, _, err := dao.Add(candidate)
		assert.NotEqual(s.T(), "", filing.ID)
		assert.NoError(s.T(), err)
	}
	{
		exists, err := dao.StrictExists("123", "123", "Issuer")
		assert.True(s.T(), exists)
		assert.NoError(s.T(), err)
	}
	{
		exists, err := dao.StrictExists("123", "123", "Reporting")
		assert.False(s.T(), exists)
		assert.NoError(s.T(), err)
	}
	{
		exists, err := dao.StrictExists("123", "1234", "Issuer")
		assert.False(s.T(), exists)
		assert.NoError(s.T(), err)
	}
	{
		exists, err := dao.StrictExists("1234", "123", "Issuer")
		assert.False(s.T(), exists)
		assert.NoError(s.T(), err)
	}
}

func (s *EdgrSuite) TestAdd() {
	dao := s.db.NewFilingDao()
	{
		candidate := &model.Filing{
			Filer:         "ABC Corp",
			Accession:     "123",
			CIK:           "123",
			AllCIKs:       []string{"123"},
			EdgarURL:      "example.com",
			EdgarTime:     time.Now(),
			FilerRelation: "Issuer",
			FormType:      "8-K",
		}
		filing, _, err := dao.Add(candidate)
		assert.NotEqual(s.T(), "", filing.ID)
		assert.NoError(s.T(), err)
	}
	{
		candidate := &model.Filing{
			Filer:         "ABC Corp",
			Accession:     "1234",
			CIK:           "123",
			AllCIKs:       []string{"123"},
			EdgarURL:      "example.com",
			EdgarTime:     time.Now(),
			FilerRelation: "Issuer",
			FormType:      "8-K",
		}
		filing, _, err := dao.Add(candidate)
		assert.NotEqual(s.T(), "", filing.ID)
		assert.NoError(s.T(), err)
	}
	{
		candidate := &model.Filing{
			Filer:         "ABC Inc",
			Accession:     "123",
			CIK:           "1234",
			AllCIKs:       []string{"1234"},
			EdgarURL:      "example.net",
			EdgarTime:     time.Now(),
			FilerRelation: "Reporting",
			FormType:      "6-K",
		}
		_, _, err := dao.Add(candidate)
		assert.Error(s.T(), err)
		//assert.NotEqual(s.T(), "", filing.ID)
	}
}

func (s *EdgrSuite) TestGet() {
	dao := s.db.NewFilingDao()
	{
		candidate := &model.Filing{
			Filer:         "ABC Corp",
			Accession:     "123",
			CIK:           "123",
			AllCIKs:       []string{"123"},
			EdgarURL:      "example.com",
			EdgarTime:     time.Now(),
			FilerRelation: "Issuer",
			FormType:      "8-K",
		}
		filing, _, err := dao.Add(candidate)
		assert.NotEqual(s.T(), "", filing.ID)
		assert.NoError(s.T(), err)
	}
	{
		exists, filing, err := dao.Get("123")
		assert.True(s.T(), exists)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), filing)
		assert.Equal(s.T(), "ABC Corp", filing.Filer)
		assert.Equal(s.T(), "123", filing.CIK)
		assert.Equal(s.T(), "123", filing.Accession)
		assert.Equal(s.T(), "example.com", filing.EdgarURL)
		assert.Equal(s.T(), "Issuer", filing.FilerRelation)
		assert.Equal(s.T(), "8-K", filing.FormType)
		assert.WithinDuration(s.T(), time.Now(), filing.Created, time.Second)
		assert.WithinDuration(s.T(), time.Now(), filing.Updated, time.Second)
		assert.WithinDuration(s.T(), filing.Created, filing.Updated, 2*time.Millisecond)
	}
	{
		exists, filing, err := dao.Get("")
		assert.False(s.T(), exists)
		assert.NoError(s.T(), err)
		assert.Nil(s.T(), filing)
	}
	{
		exists, filing, err := dao.Get("1")
		assert.False(s.T(), exists)
		assert.NoError(s.T(), err)
		assert.Nil(s.T(), filing)
	}
}

func TestByteSize(t *testing.T) {
	str := "Lorem ipsum dolor amet"
	b := []byte(str)
	s := int64(len(b))
	byteSize := byteCountDecimal(s)
	assert.Equal(t, "22 B", byteSize)
}

func (s *EdgrSuite) TestAddDocuments() {
	dao := s.db.NewFilingDao()
	{
		candidate := &model.Filing{
			Filer:         "ABC Corp",
			Accession:     "123",
			CIK:           "123",
			AllCIKs:       []string{"123"},
			EdgarURL:      "example.com",
			EdgarTime:     time.Now(),
			FilerRelation: "Issuer",
			FormType:      "8-K",
		}
		filing, _, err := dao.Add(candidate)
		assert.NotEqual(s.T(), "", filing.ID)
		assert.NoError(s.T(), err)
	}
	{
		exists, filing, err := dao.Get("123")
		assert.True(s.T(), exists)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), filing)

		docs := []*model.Document{
			&model.Document{
				Sequence:    1,
				DocType:     "4",
				Description: "Form-4",
				EdgarURL:    "https://edgar.sec.gov/filings/0594395-18-52394582",
				Body:        "<body><p>Test text.</p></body>",
			},
			&model.Document{
				Sequence:    2,
				DocType:     "Addtl Info.",
				Description: "Form-4 Info",
				EdgarURL:    "https://edgar.sec.gov/filings/0594395-18-52394582",
				Body:        "<body><p>More info.</p></body>",
			},
		}
		err = dao.AddDocuments(filing, docs)
		assert.NoError(s.T(), err)
	}
	{
		exists, filing, err := dao.Get("123")
		assert.True(s.T(), exists)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), filing)
		assert.Equal(s.T(), int64(2), filing.DocumentCount)
		assert.Equal(s.T(), "60 B", filing.TotalSizeEstimate)
		// Need a way to test that updated got updated..
		//assert.WithinDuration(s.T(), filing.Created, filing.Updated, 2*time.Millisecond)
	}
}

func (s *EdgrSuite) TestUpdate() {
	dao := s.db.NewFilingDao()
	{
		candidate := &model.Filing{
			Filer:         "ABC Corp Owner",
			Accession:     "123",
			CIK:           "002",
			AllCIKs:       []string{"002"},
			Symbol:        "ABCO",
			EdgarURL:      "example.com",
			EdgarTime:     time.Now(),
			FilerRelation: "Reporting",
			FormType:      "4/A",
		}
		filing, _, err := dao.Add(candidate)
		assert.NotEqual(s.T(), "", filing.ID)
		assert.NoError(s.T(), err)
	}

	{
		exists, filing, err := dao.Get("123")
		assert.True(s.T(), exists)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), filing)
		assert.Equal(s.T(), "002", filing.CIK)

		filing.Filer = "ABC Corp"
		oldCIK := filing.CIK
		oldSymbol := filing.Symbol

		filing.CIK = "001"
		filing.Symbol = "ABC"
		filing.FilerRelation = "Issuer"
		filing.AllCIKs = append(filing.AllCIKs, oldCIK)
		filing.AllSymbols = append(filing.AllSymbols, oldSymbol)

		updated, err := dao.Update(filing)
		assert.True(s.T(), updated)
		assert.NoError(s.T(), err)
	}
	{
		exists, filing, err := dao.Get("123")
		assert.True(s.T(), exists)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), filing)
		assert.Equal(s.T(), "001", filing.CIK)
	}
}

func (s *EdgrSuite) TestGetBySymbols() {
	dao := s.db.NewFilingDao()
	{
		candidate := &model.Filing{
			Filer:         "ABC Corp",
			Accession:     "123",
			CIK:           "123",
			AllCIKs:       []string{"123"},
			Symbol:        "ABC",
			AllSymbols:    []string{"ABC"},
			EdgarURL:      "example.com",
			EdgarTime:     time.Now(),
			FilerRelation: "Issuer",
			FormType:      "8-K",
		}
		filing, _, err := dao.Add(candidate)
		assert.NotEqual(s.T(), "", filing.ID)
		assert.NoError(s.T(), err)
	}
	{
		found, result, err := dao.GetBySymbol(&QueryParams{
			Term: "ABC",
			Size: 1,
			Sort: 1,
		})
		if assert.NoError(s.T(), err) {
			assert.True(s.T(), found)
			assert.Equal(s.T(), 1, result.TotalCount)
		}
	}
	{
		found, result, err := dao.GetBySymbol(&QueryParams{
			Term: "AB",
			Size: 1,
			Sort: 1,
		})
		if assert.NoError(s.T(), err) {
			assert.False(s.T(), found)
			assert.Nil(s.T(), result)
		}
	}
	{
		found, result, err := dao.GetBySymbol(&QueryParams{
			Term: "ABC",
			Form: "8-K",
			Size: 1,
			Sort: 1,
		})
		if assert.NoError(s.T(), err) {
			assert.True(s.T(), found)
			assert.Equal(s.T(), 1, result.TotalCount)
		}
	}
	{
		found, result, err := dao.GetBySymbol(&QueryParams{
			Term: "ABC",
			Form: "8",
			Size: 1,
			Sort: 1,
		})
		if assert.NoError(s.T(), err) {
			assert.False(s.T(), found)
			assert.Nil(s.T(), result)
		}
	}
}

package database

import (
	"time"

	"github.com/piquette/edgr/core/model"
	"github.com/stretchr/testify/assert"
)

func (s *EdgrSuite) TestDocumentDao() {
	// validate impl.
	dao := s.db.NewDocumentDao()
	assert.Implements(s.T(), (*DocumentDao)(nil), dao)
}

func (s *EdgrSuite) TestGetDocuments() {
	dao := s.db.NewFilingDao()
	ddao := s.db.NewDocumentDao()

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

		found, docs, err := ddao.GetByFilingID(filing.ID)
		assert.True(s.T(), found)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), docs, 2)
	}
	{
		exists, filing, err := dao.Get("123")
		assert.True(s.T(), exists)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), filing)
		assert.Equal(s.T(), int64(2), filing.DocumentCount)
		assert.Equal(s.T(), "60 B", filing.TotalSizeEstimate)

		found, ids, err := ddao.GetIDs(filing.ID)
		assert.True(s.T(), found)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), ids, 2)

		// check each document.
		for idx, id := range ids {
			docFound, doc, docErr := ddao.Get(id)
			assert.True(s.T(), docFound)
			assert.NoError(s.T(), docErr)
			assert.Equal(s.T(), int64(idx+1), doc.Sequence)
		}

	}
}

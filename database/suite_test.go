package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	// TestDBAddr swaps the test db address.
	TestDBAddr = "localhost"
)

// APISuite tests the api db.
type APISuite struct {
	suite.Suite
	db *Handle
}

type EdgrSuite struct {
	suite.Suite
	db *Handle
}

func (s *APISuite) SetupSuite() {
	s.db = Open(Options{
		Addr:     fmt.Sprintf("%s:5432", TestDBAddr),
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})

	reset := `
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;

	GRANT ALL ON SCHEMA public TO postgres;
	GRANT ALL ON SCHEMA public TO public;`

	_, err := s.db.Exec(reset)
	s.Assert().NoError(err)
	defer s.db.Close()
}

func (s *EdgrSuite) SetupSuite() {
	s.db = Open(Options{
		Addr:     fmt.Sprintf("%s:5432", TestDBAddr),
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})

	reset := `
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;

	GRANT ALL ON SCHEMA public TO postgres;
	GRANT ALL ON SCHEMA public TO public;`

	_, err := s.db.Exec(reset)
	s.Assert().NoError(err)
	defer s.db.Close()
}

// BeforeTest is run before.
func (s *APISuite) BeforeTest(suiteName, testName string) {
	s.db = Open(Options{
		Addr:     fmt.Sprintf("%s:5432", TestDBAddr),
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})

	//s.db.Exec("SET search_path TO pg_temp")
	//s.migrate("../../migrations/*.down.sql", true)
	for _, stmt := range getMigrations("../migrations/api/*.up.sql", false) {
		_, err := s.db.Exec(stmt)
		if err != nil {
			s.Assert().NoError(err)
		}
	}
}

// AfterTest is run after.
func (s *APISuite) AfterTest(suiteName, testName string) {
	for _, stmt := range getMigrations("../migrations/api/*.down.sql", true) {
		_, err := s.db.Exec(stmt)
		if err != nil {
			s.Assert().NoError(err)
		}
	}
	assert.Nil(s.T(), s.db.Close())
}

func (s *APISuite) exec(filepath string) {
	for _, stmt := range getMigrations(filepath, false) {
		_, err := s.db.Exec(stmt)
		if err != nil {
			s.Assert().NoError(err)
		}
	}
}

// BeforeTest is run before.
func (s *EdgrSuite) BeforeTest(suiteName, testName string) {
	s.db = Open(Options{
		Addr:     fmt.Sprintf("%s:5432", TestDBAddr),
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})

	//s.db.Exec("SET search_path TO pg_temp")
	//s.migrate("../../migrations/*.down.sql", true)
	for _, stmt := range getMigrations("../migrations/pipeline/*.up.sql", false) {
		_, err := s.db.Exec(stmt)
		if err != nil {
			s.Assert().NoError(err)
		}
	}
}

// AfterTest is run after.
func (s *EdgrSuite) AfterTest(suiteName, testName string) {
	for _, stmt := range getMigrations("../migrations/pipeline/*.down.sql", true) {
		_, err := s.db.Exec(stmt)
		if err != nil {
			s.Assert().NoError(err)
		}
	}

	assert.Nil(s.T(), s.db.Close())
}

func (s *EdgrSuite) exec(filepath string) {
	for _, stmt := range getMigrations(filepath, false) {
		_, err := s.db.Exec(stmt)
		if err != nil {
			s.Assert().NoError(err)
		}
	}
}

// TestSuite runs the db test suites.
func TestSuite(t *testing.T) {
	// API Database.
	apiDB := new(APISuite)
	suite.Run(t, apiDB)

	// Pipeline Database.
	pipelineDB := new(EdgrSuite)
	suite.Run(t, pipelineDB)
}

func getMigrations(pattern string, invert bool) []string {
	migrations, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}
	if len(migrations) == 0 {
		log.Fatalf(fmt.Sprintf("cant find migrations %s", pattern))
	}

	if invert {
		migrations = reverse(migrations)
	}

	result := []string{}
	for _, migration := range migrations {
		sql, err := ioutil.ReadFile(migration)
		if err != nil {
			log.Fatal(err)
		}

		stmt := string(sql)
		if stmt == "" {
			continue
		}

		result = append(result, stmt)
	}
	return result
}

func reverse(files []string) []string {
	for i := 0; i < len(files)/2; i++ {
		j := len(files) - i - 1
		files[i], files[j] = files[j], files[i]
	}
	return files
}

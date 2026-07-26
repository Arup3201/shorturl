package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type URLShortenerTestSuite struct {
	suite.Suite
	ctx    context.Context
	db     *gorm.DB
	userID string
}

func TestURLShortenerSuite(t *testing.T) {
	suite.Run(t, new(URLShortenerTestSuite))
}

func (s *URLShortenerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	pg, err := CreatePostgresContainer(s.ctx)
	s.Require().NoError(err, "could not start postgres container")

	s.T().Cleanup(func() {
		s.Require().NoError(pg.Terminate(s.ctx), "could not terminate postgres container")
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	s.Require().NoError(err, "failed to open gorm db")

	s.Require().NoError(db.AutoMigrate(&ShortenedUrl{}), "failed to migrate schema")

	s.db = db
}

func (s *URLShortenerTestSuite) TestShortenURL() {
	tests := []struct {
		name    string
		longURL string
	}{
		{
			name:    "shorten chatgpt url",
			longURL: "https://chatgpt.com/c/6a5baca9-ed74-83e8-9b49-881a49bd8d7c",
		},
	}

	var pgDB = &PostgresDB{s.db}
	var url string
	var err error
	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			url, err = ShortenURL(s.ctx, test.longURL, pgDB)
			s.Require().NoError(err)

			s.Require().NotEqual(url, test.longURL)
			// https://short.ly/lZyD
			s.Require().LessOrEqual(17+4, len(url))
			s.Require().GreaterOrEqual(17+14, len(url))

			row, err := gorm.G[ShortenedUrl](s.db).Where("url = ?", url).First(s.ctx)
			s.Require().NoError(err)
			s.Require().Equal(row.Original, test.longURL)
		})
	}
}

func (s *URLShortenerTestSuite) TestURLIsCollisionFree() {
	test_long_url := "https://chatgpt.com/c/6a5baca9-ed74-83e8-9b49-881a49bd8d7c"

	var results = []string{}
	var pgDB = &PostgresDB{s.db}
	var url string
	for range 1000 {
		url, _ = ShortenURL(s.ctx, test_long_url, pgDB)
		results = append(results, url)
	}

	isCollisionFree := true
	for i, r1 := range results {
		for j, r2 := range results {
			if i != j && r1 == r2 {
				isCollisionFree = false
				break
			}
		}
	}

	s.Require().True(isCollisionFree)
}

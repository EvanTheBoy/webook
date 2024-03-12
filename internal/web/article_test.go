package web

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct{
		name string
	} {
		{

		}
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

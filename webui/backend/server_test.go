package backend

import (
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

type ServerTestSuite struct {
	suite.Suite
}

func (s *ServerTestSuite) SetupTest() {}

func (s *ServerTestSuite) BeforeTest(_, _ string) {}

func (s *ServerTestSuite) TestIndex() {
	assert := assert2.New(s.T())

	res := httptest.NewRecorder()
	// Should return the file `index.html`.
	req := httptest.NewRequest("GET", "/", nil)

	newRouter(false, false).ServeHTTP(res, req)

	// Read the body of the response.
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Get the absolute path of the frontend.
	fd, err := filepath.Abs("../frontend/dist/index.html")
	if err != nil {
		panic(err)
	}

	// Open the index file and read the content.
	indexFileContent, err := ioutil.ReadFile(fd)
	if err != nil {
		panic(err)
	}

	assert.Equal(indexFileContent, responseBody, "a request to the root path should return the file "+
		"index.html")
	assert.Equal(200, res.Code, "the response should be with the code 200")
}

func (s *ServerTestSuite) AfterTest(_, _ string) {}

func (s *ServerTestSuite) TearDownTest() {}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

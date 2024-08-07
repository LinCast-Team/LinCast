package api

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"lincast/api/handlers"

	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
}

func (s *ServerTestSuite) SetupTest() {}

func (s *ServerTestSuite) BeforeTest(_, _ string) {}

// TODO Update this test
func (s *ServerTestSuite) TestIndex() {
	assert := assert2.New(s.T())

	res := httptest.NewRecorder()
	// Should return the file `index.html`.
	req := httptest.NewRequest("GET", "/", nil)

	newRouter(false, false, &handlers.Manager{}).ServeHTTP(res, req)

	// Read the body of the response.
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Get the absolute path of the frontend.
	fd, err := filepath.Abs("../webui/frontend/dist/index.html")
	if err != nil {
		panic(err)
	}

	// Open the index file and read the content.
	indexFileContent, err := os.ReadFile(fd)
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

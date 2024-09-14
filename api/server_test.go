package api

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"lincast/api/handlers"
	"lincast/models"

	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
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

	createRouter(&handlers.Manager{}).ServeHTTP(res, req)

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

func (s *ServerTestSuite) TestNewServer() {
	assert := assert2.New(s.T())

	// Mock dependencies
	db := &gorm.DB{}
	manualUpdate := make(chan *models.Podcast)

	// Test cases
	tests := []struct {
		port        uint
		localServer bool
		devMode     bool
		logRequests bool
	}{
		{8080, true, false, false},
		{8081, false, true, true},
		{8082, false, false, true},
		{8083, true, true, false},
	}

	for _, tt := range tests {
		server := New(tt.port, tt.localServer, tt.devMode, tt.logRequests, db, manualUpdate)

		expectedAddr := "127.0.0.1:" + strconv.Itoa(int(tt.port))
		if !tt.localServer {
			expectedAddr = ":" + strconv.Itoa(int(tt.port))
		}

		assert.Equal(expectedAddr, server.Addr, "server address should match the expected address")
		assert.Equal(time.Second*15, server.ReadTimeout, "server read timeout should be 15 seconds")
		assert.Equal(time.Second*15, server.WriteTimeout, "server write timeout should be 15 seconds")
		assert.Equal(8000, server.MaxHeaderBytes, "server max header bytes should be 8KB")
	}
}

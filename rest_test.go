package leaseweb

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = ctxT{}
var testApiKey = "test-api-key"

type serverErrorTest struct {
	Title         string
	MockServer    func(http.ResponseWriter, *http.Request)
	FunctionCall  func() (interface{}, error)
	ExpectedError LeasewebError
}

type ctxT struct {
	ts            *httptest.Server
	oldHost       string
	oldHttpClient *http.Client
}

func assertServerErrorTests(t *testing.T, serverErrorTests []serverErrorTest) {
	for _, serverErrorTest := range serverErrorTests {
		setup(serverErrorTest.MockServer)
		defer teardown()
		resp, err := serverErrorTest.FunctionCall()
		assert := assert.New(t)
		assert.Empty(resp)
		assert.NotNil(err)
		assert.Equal(err.Error(), serverErrorTest.ExpectedError.ErrorMessage)
		lswErr, ok := err.(*LeasewebError)
		assert.Equal(true, ok)
		assert.Equal(lswErr.ErrorMessage, serverErrorTest.ExpectedError.ErrorMessage)
		assert.Equal(lswErr.CorrelationId, serverErrorTest.ExpectedError.CorrelationId)
		assert.Equal(lswErr.ErrorCode, serverErrorTest.ExpectedError.ErrorCode)
		assert.Equal(lswErr.Reference, serverErrorTest.ExpectedError.Reference)
		assert.Equal(lswErr.UserMessage, serverErrorTest.ExpectedError.UserMessage)
	}
}

func setup(handler http.HandlerFunc) *httptest.Server {
	ts := httptest.NewTLSServer(handler)
	ctx.ts = ts
	_url, err := url.Parse(ts.URL)
	if err != nil {
		panic(err)
	}

	ctx.oldHost = lswClient.baseUrl
	ctx.oldHttpClient = lswClient.client

	lswClient.baseUrl = "https://" + _url.Host
	lswClient.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return ts
}

func teardown() {
	ctx.ts.Close()
	lswClient.baseUrl = ctx.oldHost
	lswClient.client = ctx.oldHttpClient
}

func TestMain(m *testing.M) {
	InitLeasewebClient(testApiKey)
	os.Exit(m.Run())
}

package leaseweb

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRanges(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{"_metadata":{"limit": 10, "offset": 0, "totalCount": 2}, "ranges": [
			{
				"id": "85.17.0.0_17",
				"range": "85.17.0.0/17",
				"customerId": "10001234",
				"salesOrgId": "2000",
				"location": "AMS-01",
				"type": "SITE"
			},
			{
				"id": "86.17.0.1_17",
				"range": "86.17.0.1/17",
				"customerId": "10001234",
				"salesOrgId": "2000",
				"location": "AMS",
				"type": "METRO"
			}
		]}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.ListRanges()

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(response.Metadata.TotalCount, 2)
	assert.Equal(response.Metadata.Offset, 0)
	assert.Equal(response.Metadata.Limit, 10)
	assert.Equal(len(response.Ranges), 2)

	range1 := response.Ranges[0]
	assert.Equal(range1.Id, "85.17.0.0_17")
	assert.Equal(range1.Range, "85.17.0.0/17")
	assert.Equal(range1.CustomerId, "10001234")
	assert.Equal(range1.SalesOrgId, "2000")
	assert.Equal(range1.Location, "AMS-01")
	assert.Equal(range1.Type, "SITE")

	range2 := response.Ranges[1]
	assert.Equal(range2.Id, "86.17.0.1_17")
	assert.Equal(range2.Range, "86.17.0.1/17")
	assert.Equal(range2.CustomerId, "10001234")
	assert.Equal(range2.SalesOrgId, "2000")
	assert.Equal(range2.Location, "AMS")
	assert.Equal(range2.Type, "METRO")
}

func TestListRangesPaginateAndFilter(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{"_metadata":{"limit": 10, "offset": 1, "totalCount": 11}, "ranges": [
			{
				"id": "85.17.0.0_17",
				"range": "85.17.0.0/17",
				"customerId": "10001234",
				"salesOrgId": "2000",
				"location": "AMS-01",
				"type": "SITE"
			}
		]}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.ListRanges(1, 10, []string{"SITE", "METRO"}, "AMS-01")

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(response.Metadata.TotalCount, 11)
	assert.Equal(response.Metadata.Offset, 1)
	assert.Equal(response.Metadata.Limit, 10)
	assert.Equal(len(response.Ranges), 1)

	range1 := response.Ranges[0]
	assert.Equal(range1.Id, "85.17.0.0_17")
	assert.Equal(range1.Range, "85.17.0.0/17")
	assert.Equal(range1.CustomerId, "10001234")
	assert.Equal(range1.SalesOrgId, "2000")
	assert.Equal(range1.Location, "AMS-01")
	assert.Equal(range1.Type, "SITE")
}

func TestListRangesServerErrors(t *testing.T) {
	serverErrorTests := []serverErrorTest{
		{
			Title: "error 403",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "ACCESS_DENIED", "errorMessage": "The access token is expired or invalid."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.ListRanges()
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "ACCESS_DENIED",
				ErrorMessage:  "The access token is expired or invalid.",
			},
		},
		{
			Title: "error 500",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "500", "errorMessage": "The server encountered an unexpected condition that prevented it from fulfilling the request."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.ListRanges()
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "500",
				ErrorMessage:  "The server encountered an unexpected condition that prevented it from fulfilling the request.",
			},
		},
		{
			Title: "error 503",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "503", "errorMessage": "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.ListRanges()
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "503",
				ErrorMessage:  "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
			},
		},
	}
	assertServerErrorTests(t, serverErrorTests)
}

func TestGetRange(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{
			"id": "85.17.0.0_17",
			"range": "88.17.0.0/17",
			"customerId": "10001234",
			"salesOrgId": "2000",
			"location": "AMS-01",
			"type": "SITE"
		}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.GetRange("123456789")

	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(response.Id, "85.17.0.0_17")
	assert.Equal(response.Range, "88.17.0.0/17")
	assert.Equal(response.CustomerId, "10001234")
	assert.Equal(response.SalesOrgId, "2000")
	assert.Equal(response.Location, "AMS-01")
	assert.Equal(response.Type, "SITE")
}

func TestGetRangeServerErrors(t *testing.T) {
	serverErrorTests := []serverErrorTest{
		{
			Title: "error 403",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"errorCode": "ACCESS_DENIED", "errorMessage": "The access token is expired or invalid."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.GetRange("123456789")
			},
			ExpectedError: LeasewebError{
				ErrorCode:    "ACCESS_DENIED",
				ErrorMessage: "The access token is expired or invalid.",
			},
		},
		{
			Title: "error 404",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, `{"correlationId": "39e010ed-0e93-42c3-c28f-3ffc373553d5", "errorCode": "404", "errorMessage": "Range with id 88.17.0.0_17 does not exist"}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.GetRange("123456789")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "39e010ed-0e93-42c3-c28f-3ffc373553d5",
				ErrorCode:     "404",
				ErrorMessage:  "Range with id 88.17.0.0_17 does not exist",
			},
		},
		{
			Title: "error 500",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "500", "errorMessage": "The server encountered an unexpected condition that prevented it from fulfilling the request."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.GetRange("123456789")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "500",
				ErrorMessage:  "The server encountered an unexpected condition that prevented it from fulfilling the request.",
			},
		},
		{
			Title: "error 503",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "503", "errorMessage": "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.GetRange("123456789")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "503",
				ErrorMessage:  "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
			},
		},
	}
	assertServerErrorTests(t, serverErrorTests)
}

func TestListRangeDefinitions(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{"_metadata":{"limit": 10, "offset": 0, "totalCount": 2}, "floatingIpDefinitions": [
			{
				"id": "88.17.34.108_32",
				"rangeId": "88.17.0.0_17",
				"location": "AMS-01",
				"type": "SITE",
				"customerId": "10001234",
				"salesOrgId": "2000",
				"floatingIp": "88.17.34.108/32",
				"anchorIp": "95.10.126.1",
				"status": "ACTIVE",
				"createdAt": "2019-03-13T09:10:02+0000",
				"updatedAt": "2019-03-13T09:10:02+0000"
			},
			{
				"id": "88.17.34.109_32",
				"rangeId": "88.17.0.0_17",
				"location": "AMS-01",
				"type": "SITE",
				"customerId": "10001234",
				"salesOrgId": "2000",
				"floatingIp": "88.17.34.109/32",
				"anchorIp": "95.10.126.12",
				"status": "ACTIVE",
				"createdAt": "2019-03-13T09:10:02+0000",
				"updatedAt": "2019-03-13T09:10:02+0000"
			}
		]}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.ListRangeDefinitions("123456789")

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(response.Metadata.TotalCount, 2)
	assert.Equal(response.Metadata.Offset, 0)
	assert.Equal(response.Metadata.Limit, 10)
	assert.Equal(len(response.FloatingIpDefinitions), 2)

	floatingIpDefinition1 := response.FloatingIpDefinitions[0]
	assert.Equal(floatingIpDefinition1.Id, "88.17.34.108_32")
	assert.Equal(floatingIpDefinition1.RangeId, "88.17.0.0_17")
	assert.Equal(floatingIpDefinition1.CustomerId, "10001234")
	assert.Equal(floatingIpDefinition1.SalesOrgId, "2000")
	assert.Equal(floatingIpDefinition1.Location, "AMS-01")
	assert.Equal(floatingIpDefinition1.Type, "SITE")
	assert.Equal(floatingIpDefinition1.FloatingIp, "88.17.34.108/32")
	assert.Equal(floatingIpDefinition1.AnchorIp, "95.10.126.1")
	assert.Equal(floatingIpDefinition1.Status, "ACTIVE")
	assert.Equal(floatingIpDefinition1.CreatedAt, "2019-03-13T09:10:02+0000")
	assert.Equal(floatingIpDefinition1.UpdatedAt, "2019-03-13T09:10:02+0000")

	floatingIpDefinition2 := response.FloatingIpDefinitions[1]
	assert.Equal(floatingIpDefinition2.Id, "88.17.34.109_32")
	assert.Equal(floatingIpDefinition2.RangeId, "88.17.0.0_17")
	assert.Equal(floatingIpDefinition2.CustomerId, "10001234")
	assert.Equal(floatingIpDefinition2.SalesOrgId, "2000")
	assert.Equal(floatingIpDefinition2.Location, "AMS-01")
	assert.Equal(floatingIpDefinition2.Type, "SITE")
	assert.Equal(floatingIpDefinition2.FloatingIp, "88.17.34.109/32")
	assert.Equal(floatingIpDefinition2.AnchorIp, "95.10.126.12")
	assert.Equal(floatingIpDefinition2.Status, "ACTIVE")
	assert.Equal(floatingIpDefinition2.CreatedAt, "2019-03-13T09:10:02+0000")
	assert.Equal(floatingIpDefinition2.UpdatedAt, "2019-03-13T09:10:02+0000")
}

func TestListRangeDefinitionsPaginateAndFilter(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{"_metadata":{"limit": 10, "offset": 1, "totalCount": 11}, "floatingIpDefinitions": [
			{
				"id": "88.17.34.108_32",
				"rangeId": "88.17.0.0_17",
				"location": "AMS-01",
				"type": "SITE",
				"customerId": "10001234",
				"salesOrgId": "2000",
				"floatingIp": "88.17.34.108/32",
				"anchorIp": "95.10.126.1",
				"status": "ACTIVE",
				"createdAt": "2019-03-13T09:10:02+0000",
				"updatedAt": "2019-03-13T09:10:02+0000"
			}
		]}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.ListRangeDefinitions("123456789", 1, 10, []string{"SITE", "METRO"}, "AMS-01")

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(response.Metadata.TotalCount, 11)
	assert.Equal(response.Metadata.Offset, 1)
	assert.Equal(response.Metadata.Limit, 10)

	floatingIpDefinition1 := response.FloatingIpDefinitions[0]
	assert.Equal(floatingIpDefinition1.Id, "88.17.34.108_32")
	assert.Equal(floatingIpDefinition1.RangeId, "88.17.0.0_17")
	assert.Equal(floatingIpDefinition1.CustomerId, "10001234")
	assert.Equal(floatingIpDefinition1.SalesOrgId, "2000")
	assert.Equal(floatingIpDefinition1.Location, "AMS-01")
	assert.Equal(floatingIpDefinition1.Type, "SITE")
	assert.Equal(floatingIpDefinition1.FloatingIp, "88.17.34.108/32")
	assert.Equal(floatingIpDefinition1.AnchorIp, "95.10.126.1")
	assert.Equal(floatingIpDefinition1.Status, "ACTIVE")
	assert.Equal(floatingIpDefinition1.CreatedAt, "2019-03-13T09:10:02+0000")
	assert.Equal(floatingIpDefinition1.UpdatedAt, "2019-03-13T09:10:02+0000")
}

func TestListRangeDefinitionsServerErrors(t *testing.T) {
	serverErrorTests := []serverErrorTest{
		{
			Title: "error 403",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"errorCode": "ACCESS_DENIED", "errorMessage": "The access token is expired or invalid."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.ListRangeDefinitions("123456789")
			},
			ExpectedError: LeasewebError{
				ErrorCode:    "ACCESS_DENIED",
				ErrorMessage: "The access token is expired or invalid.",
			},
		},
		{
			Title: "error 404",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, `{"correlationId": "39e010ed-0e93-42c3-c28f-3ffc373553d5", "errorCode": "404", "errorMessage": "Range with id 88.17.0.0_17 does not exist"}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.ListRangeDefinitions("123456789")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "39e010ed-0e93-42c3-c28f-3ffc373553d5",
				ErrorCode:     "404",
				ErrorMessage:  "Range with id 88.17.0.0_17 does not exist",
			},
		},
		{
			Title: "error 500",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "500", "errorMessage": "The server encountered an unexpected condition that prevented it from fulfilling the request."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.ListRangeDefinitions("123456789")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "500",
				ErrorMessage:  "The server encountered an unexpected condition that prevented it from fulfilling the request.",
			},
		},
		{
			Title: "error 503",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "503", "errorMessage": "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.ListRangeDefinitions("123456789")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "503",
				ErrorMessage:  "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
			},
		},
	}
	assertServerErrorTests(t, serverErrorTests)
}

func TestCreateRangeDefinition(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{
			"id": "88.17.34.108_32",
			"rangeId": "88.17.0.0_17",
			"location": "AMS-01",
			"type": "SITE",
			"customerId": "10001234",
			"salesOrgId": "2000",
			"floatingIp": "88.17.34.108/32",
			"anchorIp": "95.10.126.1",
			"status": "ACTIVE",
			"createdAt": "2019-03-13T09:10:02+0000",
			"updatedAt": "2019-03-13T09:10:02+0000"
		}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.CreateRangeDefinition("10.0.0.0_29", "88.17.0.5/32", "95.10.126.1")

	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(response.Id, "88.17.34.108_32")
	assert.Equal(response.RangeId, "88.17.0.0_17")
	assert.Equal(response.CustomerId, "10001234")
	assert.Equal(response.SalesOrgId, "2000")
	assert.Equal(response.Location, "AMS-01")
	assert.Equal(response.Type, "SITE")
	assert.Equal(response.FloatingIp, "88.17.34.108/32")
	assert.Equal(response.AnchorIp, "95.10.126.1")
	assert.Equal(response.Status, "ACTIVE")
	assert.Equal(response.CreatedAt, "2019-03-13T09:10:02+0000")
	assert.Equal(response.UpdatedAt, "2019-03-13T09:10:02+0000")
}

func TestCreateRangeDefinitionServerError(t *testing.T) {
	serverErrorTests := []serverErrorTest{
		{
			Title: "error 400",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"errorCode": "400", "errorMessage": "Validation Failed"}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.CreateRangeDefinition("10.0.0.0_29", "88.17.0.5/32", "95.10.126.1")
			},
			ExpectedError: LeasewebError{
				ErrorCode:    "400",
				ErrorMessage: "Validation Failed",
			},
		},
		{
			Title: "error 403",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"errorCode": "ACCESS_DENIED", "errorMessage": "The access token is expired or invalid."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.CreateRangeDefinition("10.0.0.0_29", "88.17.0.5/32", "95.10.126.1")
			},
			ExpectedError: LeasewebError{
				ErrorCode:    "ACCESS_DENIED",
				ErrorMessage: "The access token is expired or invalid.",
			},
		},
		{
			Title: "error 500",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "500", "errorMessage": "The server encountered an unexpected condition that prevented it from fulfilling the request."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.CreateRangeDefinition("10.0.0.0_29", "88.17.0.5/32", "95.10.126.1")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "500",
				ErrorMessage:  "The server encountered an unexpected condition that prevented it from fulfilling the request.",
			},
		},
		{
			Title: "error 503",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "503", "errorMessage": "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.CreateRangeDefinition("10.0.0.0_29", "88.17.0.5/32", "95.10.126.1")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "503",
				ErrorMessage:  "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
			},
		},
	}
	assertServerErrorTests(t, serverErrorTests)
}

func TestGetRangeDefinition(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{
			"id": "88.17.34.108_32",
			"rangeId": "88.17.0.0_17",
			"location": "AMS-01",
			"type": "SITE",
			"customerId": "10001234",
			"salesOrgId": "2000",
			"floatingIp": "88.17.34.108/32",
			"anchorIp": "95.10.126.1",
			"status": "ACTIVE",
			"createdAt": "2019-03-13T09:10:02+0000",
			"updatedAt": "2019-03-13T09:10:02+0000"
		}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.GetRangeDefinition("88.17.0.0_17", "88.17.34.108_32")

	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(response.Id, "88.17.34.108_32")
	assert.Equal(response.RangeId, "88.17.0.0_17")
	assert.Equal(response.CustomerId, "10001234")
	assert.Equal(response.SalesOrgId, "2000")
	assert.Equal(response.Location, "AMS-01")
	assert.Equal(response.Type, "SITE")
	assert.Equal(response.FloatingIp, "88.17.34.108/32")
	assert.Equal(response.AnchorIp, "95.10.126.1")
	assert.Equal(response.Status, "ACTIVE")
	assert.Equal(response.CreatedAt, "2019-03-13T09:10:02+0000")
	assert.Equal(response.UpdatedAt, "2019-03-13T09:10:02+0000")
}

func TestGetRangeDefinitionServerErrors(t *testing.T) {
	serverErrorTests := []serverErrorTest{
		{
			Title: "error 403",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"errorCode": "ACCESS_DENIED", "errorMessage": "The access token is expired or invalid."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.GetRangeDefinition("88.17.0.0_17", "88.17.34.108_32")
			},
			ExpectedError: LeasewebError{
				ErrorCode:    "ACCESS_DENIED",
				ErrorMessage: "The access token is expired or invalid.",
			},
		},
		{
			Title: "error 500",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "500", "errorMessage": "The server encountered an unexpected condition that prevented it from fulfilling the request."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.GetRangeDefinition("88.17.0.0_17", "88.17.34.108_32")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "500",
				ErrorMessage:  "The server encountered an unexpected condition that prevented it from fulfilling the request.",
			},
		},
		{
			Title: "error 503",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "503", "errorMessage": "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.GetRangeDefinition("88.17.0.0_17", "88.17.34.108_32")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "503",
				ErrorMessage:  "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
			},
		},
	}
	assertServerErrorTests(t, serverErrorTests)
}

func TestUpdateRangeDefinition(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{
			"id": "88.17.34.108_32",
			"rangeId": "88.17.0.0_17",
			"location": "AMS-01",
			"type": "SITE",
			"customerId": "10001234",
			"salesOrgId": "2000",
			"floatingIp": "88.17.34.108/32",
			"anchorIp": "95.10.126.1",
			"status": "ACTIVE",
			"createdAt": "2019-03-13T09:10:02+0000",
			"updatedAt": "2019-03-13T09:10:02+0000"
		}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.UpdateRangeDefinition("88.17.0.0_17", "88.17.34.108_32", "95.10.126.1")

	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(response.Id, "88.17.34.108_32")
	assert.Equal(response.RangeId, "88.17.0.0_17")
	assert.Equal(response.CustomerId, "10001234")
	assert.Equal(response.SalesOrgId, "2000")
	assert.Equal(response.Location, "AMS-01")
	assert.Equal(response.Type, "SITE")
	assert.Equal(response.FloatingIp, "88.17.34.108/32")
	assert.Equal(response.AnchorIp, "95.10.126.1")
	assert.Equal(response.Status, "ACTIVE")
	assert.Equal(response.CreatedAt, "2019-03-13T09:10:02+0000")
	assert.Equal(response.UpdatedAt, "2019-03-13T09:10:02+0000")
}

func TestUpdateRangeDefinitionServerErrors(t *testing.T) {
	serverErrorTests := []serverErrorTest{
		{
			Title: "error 400",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"correlationId": "945bef2e-1caf-4027-bd0a-8976848f3dee", "errorCode": "400", "errorMessage": "Validation Failed"}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.UpdateRangeDefinition("wrong 1", "88.17.34.108_32", "95.10.126.1")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "945bef2e-1caf-4027-bd0a-8976848f3dee",
				ErrorCode:     "400",
				ErrorMessage:  "Validation Failed",
			},
		},
		{
			Title: "error 403",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"errorCode": "ACCESS_DENIED", "errorMessage": "The access token is expired or invalid."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.UpdateRangeDefinition("wrong 1", "88.17.34.108_32", "95.10.126.1")
			},
			ExpectedError: LeasewebError{
				ErrorCode:    "ACCESS_DENIED",
				ErrorMessage: "The access token is expired or invalid.",
			},
		},
		{
			Title: "error 500",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "500", "errorMessage": "The server encountered an unexpected condition that prevented it from fulfilling the request."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.UpdateRangeDefinition("wrong 1", "88.17.34.108_32", "95.10.126.1")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "500",
				ErrorMessage:  "The server encountered an unexpected condition that prevented it from fulfilling the request.",
			},
		},
		{
			Title: "error 503",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "503", "errorMessage": "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.UpdateRangeDefinition("wrong 1", "88.17.34.108_32", "95.10.126.1")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "503",
				ErrorMessage:  "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
			},
		},
	}
	assertServerErrorTests(t, serverErrorTests)
}

func TestRemoveRangeDefinition(t *testing.T) {
	setup(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
		fmt.Fprintf(w, `{
			"id": "88.17.34.108_32",
			"rangeId": "88.17.0.0_17",
			"location": "AMS-01",
			"type": "SITE",
			"customerId": "10001234",
			"salesOrgId": "2000",
			"floatingIp": "88.17.34.108/32",
			"anchorIp": "95.10.126.1",
			"status": "ACTIVE",
			"createdAt": "2019-03-13T09:10:02+0000",
			"updatedAt": "2019-03-13T09:10:02+0000"
		}`)
	})
	defer teardown()

	floatingIpApi := FloatingIpApi{}
	response, err := floatingIpApi.RemoveRangeDefinition("88.17.0.0_17", "88.17.34.108_32")

	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(response.Id, "88.17.34.108_32")
	assert.Equal(response.RangeId, "88.17.0.0_17")
	assert.Equal(response.CustomerId, "10001234")
	assert.Equal(response.SalesOrgId, "2000")
	assert.Equal(response.Location, "AMS-01")
	assert.Equal(response.Type, "SITE")
	assert.Equal(response.FloatingIp, "88.17.34.108/32")
	assert.Equal(response.AnchorIp, "95.10.126.1")
	assert.Equal(response.Status, "ACTIVE")
	assert.Equal(response.CreatedAt, "2019-03-13T09:10:02+0000")
	assert.Equal(response.UpdatedAt, "2019-03-13T09:10:02+0000")
}

func TestRemoveRangeDefinitionServerErrors(t *testing.T) {

	serverErrorTests := []serverErrorTest{
		{
			Title: "error 403",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"errorCode": "ACCESS_DENIED", "errorMessage": "The access token is expired or invalid."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.RemoveRangeDefinition("88.17.0.0_17", "88.17.34.108_32")
			},
			ExpectedError: LeasewebError{
				ErrorCode:    "ACCESS_DENIED",
				ErrorMessage: "The access token is expired or invalid.",
			},
		},
		{
			Title: "error 500",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "500", "errorMessage": "The server encountered an unexpected condition that prevented it from fulfilling the request."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.RemoveRangeDefinition("88.17.0.0_17", "88.17.34.108_32")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "500",
				ErrorMessage:  "The server encountered an unexpected condition that prevented it from fulfilling the request.",
			},
		},
		{
			Title: "error 503",
			MockServer: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, testApiKey, r.Header.Get("x-lsw-auth"))
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, `{"correlationId": "289346a1-3eaf-4da4-b707-62ef12eb08be", "errorCode": "503", "errorMessage": "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server."}`)
			},
			FunctionCall: func() (interface{}, error) {
				return FloatingIpApi{}.RemoveRangeDefinition("88.17.0.0_17", "88.17.34.108_32")
			},
			ExpectedError: LeasewebError{
				CorrelationId: "289346a1-3eaf-4da4-b707-62ef12eb08be",
				ErrorCode:     "503",
				ErrorMessage:  "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
			},
		},
	}
	assertServerErrorTests(t, serverErrorTests)
}

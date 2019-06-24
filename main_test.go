// main_test.go

package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
   "github.com/cocooma/host_catalog"
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize(
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM hosts")
	a.DB.Exec("ALTER SEQUENCE hosts_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS hosts
(
    id SERIAL,
    name TEXT UNIQUE NOT NULL,
    ip VARCHAR(15) NOT NULL,
    CONSTRAINT hosts_pkey PRIMARY KEY (id)
)`

func TestHealthEndpoint(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/health", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != `{"health":"true"}` {
		t.Errorf("Expected Health Return. Got %s", body)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/hosts", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentHost(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/host/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Host not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Host not found'. Got '%s'", m["error"])
	}
}

func TestCreateHost(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test host","ip":"11.22.33.44"}`)

	req, _ := http.NewRequest("POST", "/host", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test host" {
		t.Errorf("Expected host name to be 'test host'. Got '%v'", m["name"])
	}

	if m["ip"] != "11.22.33.44" {
		t.Errorf("Expected host ip to be '11.22.33.44'. Got '%v'", m["ip"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected host ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetHost(t *testing.T) {
	clearTable()
	addHosts(1)

	req, _ := http.NewRequest("GET", "/host/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addHosts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO hosts(name, ip) VALUES($1, $2)", "Host "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestUpdateHost(t *testing.T) {
	clearTable()
	addHosts(1)

	req, _ := http.NewRequest("GET", "/host/1", nil)
	response := executeRequest(req)
	var originalHost map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalHost)

	payload := []byte(`{"name":"test host - updated name","ip":"11.22.33.44"}`)

	req, _ = http.NewRequest("PUT", "/host/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalHost["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalHost["id"], m["id"])
	}

	if m["name"] == originalHost["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalHost["name"], m["name"], m["name"])
	}

	if m["ip"] == originalHost["ip"] {
		t.Errorf("Expected the ip to change from '%v' to '%v'. Got '%v'", originalHost["ip"], m["ip"], m["ip"])
	}
}

func TestDeleteHost(t *testing.T) {
	clearTable()
	addHosts(1)

	req, _ := http.NewRequest("GET", "/host/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/host/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/host/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
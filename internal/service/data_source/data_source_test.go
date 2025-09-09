package datasources_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDatabricksDataSource(t *testing.T) {
	payload := map[string]interface{}{
		"friendlyDatabaseName":         "My Databricks Connection",
		"databaseType":                 "databricks",
		"databasePassword":             "jose17mysvcprincipaloauthsecret10690",
		"hostname":                     "https://6my-45wrkbook34-host32.cloud.databricks.com/",
		"databricksServicePrincipalId": "3bda4my-serv-princ-uuid-f016f8248c6a",
		"databricksClusterId":          "23my-3dbx23-cluster3432",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/databases", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	// Replace with actual handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": true}`))
	})

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"success": true}`, rec.Body.String())
}

func TestCreateSnowflakeDataSource(t *testing.T) {
	payload := map[string]interface{}{
		"friendlyDatabaseName": "My Snowflake Database",
		"databaseType":         "snowflake_external_functions",
		"databaseUsername":     "mysnowflakeusername",
		"databasePassword":     "mysnowflakepassword",
		"databaseName":         "MY_DATABASE_NAME",
		"hostname":             "dfr432245.us-west.snowflakecomputing.com",
		"databasePort":         3305,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/databases", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	// Replace with actual handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": true}`))
	})

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"success": true}`, rec.Body.String())
}

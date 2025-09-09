package client

import (
	"context"
	"fmt"
	"net/http"
)

// Database represents a database object returned by the ALTR API.
type DataSource struct {
	ID                   int64  `json:"id"`
	ClientID             string `json:"clientId"`
	FriendlyDatabaseName string `json:"friendlyDatabaseName"`
	DatabaseType         string `json:"databaseType"`
	DatabaseName         string `json:"databaseName"`
	DatabaseUsername     string `json:"databaseUsername"`
	SFCount              int64  `json:"SFCount"`
	InProgress           int64  `json:"inProgress"`
	LastConnectedTime    string `json:"lastConnectedTime"`
	DatabasePort         int64  `json:"databasePort"`
	Hostname             string `json:"hostname"`
}

type CreateDataSourceInput struct {
	ID                   int64  `json:"id,omitempty"`
	NumberConnections    int    `json:"maxNumberOfConnections,omitempty"`
	NumberBatches        int    `json:"maxNumberOfBatches,omitempty"`
	Warehouse            string `json:"warehouseName,omitempty"`
	ShouldClassify       bool   `json:"shouldClassify,omitempty"`
	ClassificationType   string `json:"classificationType,omitempty"`
	ImportHistoryEnabled bool   `json:"dataUsageHistory,omitempty"`
	Role                 string `json:"snowflakeRole,omitempty"`
	Hostname             string `json:"hostname"`
	ConnectionStringUsed bool   `json:"connectionString,omitempty"`
	SPCAccountID         int64  `json:"accountId,omitempty"`
	DatabaseUsername     string `json:"databaseUsername,omitempty"`
	DatabasePassword     string `json:"databasePassword,omitempty"`
	DatabasePort         int64  `json:"databasePort,omitempty"`
	DatabaseName         string `json:"databaseName,omitempty"`
	DatabaseType         string `json:"databaseType,omitempty"`
	FriendlyDatabaseName string `json:"friendlyDatabaseName,omitempty"`
}

// GetDatasources fetches the list of datasources from the ALTR API.
func (c *Client) GetDataSources(ctx context.Context) ([]DataSource, error) {
	resp, err := c.makeRequest(http.MethodGet, "/databases", nil, "external")
	if err != nil {
		return nil, fmt.Errorf("failed to get impersonation policy: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var response struct {
		Data []DataSource `json:"data"`
	}

	if err := handleAPIResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to get impersonation policy: %w", err)
	}

	return response.Data, nil
}

func (c *Client) CreateDataSource(input CreateDataSourceInput) (*DataSource, error) {
	resp, err := c.makeRequest(http.MethodPost, "/databases", input, "external")
	if err != nil {
		return nil, fmt.Errorf("failed to create impersonation policy: %w", err)
	}

	// Define a temporary structure to parse the response
	var response struct {
		Data DataSource `json:"data"`
	}

	// Parse the response into the temporary structure
	if err := handleAPIResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to parse impersonation policy response: %w", err)
	}

	// Return the parsed policy
	return &response.Data, nil
}

// DeleteDataSource deletes a datasource by ID
func (c *Client) DeleteDataSource(dataSourceID int64) error {
	resp, err := c.makeRequest(http.MethodDelete, fmt.Sprintf("/databases/%d", dataSourceID), nil, "external")
	if err != nil {
		return fmt.Errorf("failed to delete datasource: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if err := handleAPIResponse(resp, nil); err != nil {
		return fmt.Errorf("failed to delete datasource: %w", err)
	}

	return nil
}

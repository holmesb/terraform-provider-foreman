package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HanseMerkur/terraform-provider-utils/log"
)

const (
	OverrideValueEndpointPrefix = "smart_class_parameters/%d/override_values"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

// The FormanOverrideValue API model represents an override value on top of a
// smart class parameter
type ForemanOverrideValue struct {
	// Inherits the base object's attributes
	ForemanObject

	// The Smart Class Parameter that is modifed by this override value
	SmartClassParameterID string

	// Override match
	Match string `json:"match"`

	// Override value, required if omit is false
	Value string `json:"value"`

	// Foreman will not send this parameter in classification output
	Omit bool `json:"omit"`
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// CreateOverrideValue creates a new ForemanOverrideValue with the attributes
// of the supplied ForemanOverrideValue reference and returns the created
// ForemanOverrideValue reference.  The returned reference will have its ID
// and other API default values set by this function.
func (c *Client) CreateOverrideValue(t *ForemanOverrideValue) (*ForemanOverrideValue, error) {
	log.Tracef("foreman/api/overridevalue.go#Create")

	reqEndpoint := fmt.Sprintf(OverrideValueEndpointPrefix, t.SmartClassParameterID)

	tJSONBytes, jsonEncErr := c.WrapJSON("override_value", t)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("overrideValuesJSONBytes: [%s]", tJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPost,
		reqEndpoint,
		bytes.NewBuffer(tJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var createdOverrideValue ForemanOverrideValue
	sendErr := c.SendAndParse(req, &createdOverrideValue)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("createdOverrideValue: [%+v]", createdOverrideValue)

	return &createdOverrideValue, nil
}

// ReadOverrideValue reads the attributes of a ForemanOverrideValue
// identified by the supplied ID and returns a ForemanOverrideValue reference.
func (c *Client) ReadOverrideValue(id int) (*ForemanOverrideValue, error) {
	log.Tracef("foreman/api/overridevalue.go#Read")

	reqEndpoint := fmt.Sprintf(OverrideValueEndpointPrefix+"/%d", t.SmartClassParameterID, id)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readOverrideValue ForemanOverrideValue
	sendErr := c.SendAndParse(req, &readOverrideValue)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("readOverrideValue: [%+v]", readOverrideValue)

	return &readOverrideValue, nil
}

// UpdateOverrideValue updates a ForemanOverrideValue's attributes.  The
// partition table with the ID of the supplied ForemanOverrideValue will be
// updated. A new ForemanOverrideValue reference is returned with the
// attributes from the result of the update operation.
func (c *Client) UpdateOverrideValue(t *ForemanOverrideValue) (*ForemanOverrideValue, error) {
	log.Tracef("foreman/api/overridevalue.go#Update")

	reqEndpoint := fmt.Sprintf(OverrideValueEndpointPrefix+"/%d", t.SmartClassParameterID, id)

	tJSONBytes, jsonEncErr := c.WrapJSON("ptable", t)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("overridevalueJSONBytes: [%s]", tJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPut,
		reqEndpoint,
		bytes.NewBuffer(tJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var updatedOverrideValue ForemanOverrideValue
	sendErr := c.SendAndParse(req, &updatedOverrideValue)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("updatedOverrideValue: [%+v]", updatedOverrideValue)

	return &updatedOverrideValue, nil
}

// DeleteOverrideValue deletes the ForemanOverrideValue identified by the
// supplied ID
func (c *Client) DeleteOverrideValue(id int) error {
	log.Tracef("foreman/api/overridevalue.go#Delete")

	reqEndpoint := fmt.Sprintf(OverrideValueEndpointPrefix+"/%d", t.SmartClassParameterID, id)

	req, reqErr := c.NewRequest(
		http.MethodDelete,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return reqErr
	}

	return c.SendAndParse(req, nil)
}

// -----------------------------------------------------------------------------
// Query Implementation
// -----------------------------------------------------------------------------

// QueryOverrideValue queries for a ForemanOverrideValue based on the
// attributes of the supplied ForemanOverrideValue reference and returns a
// QueryResponse struct containing query/response metadata and the matching
// partition tables.
func (c *Client) QueryOverrideValue(t *ForemanOverrideValue) (QueryResponse, error) {
	log.Tracef("foreman/api/overridevalue.go#Search")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s/%d", OverrideValueEndpointPrefix, t.Id)
	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return queryResponse, reqErr
	}

	// dynamically build the query based on the attributes
	reqQuery := req.URL.Query()
	name := `"` + t.Name + `"`
	reqQuery.Set("search", "name="+name)

	req.URL.RawQuery = reqQuery.Encode()
	sendErr := c.SendAndParse(req, &queryResponse)
	if sendErr != nil {
		return queryResponse, sendErr
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	// Results will be Unmarshaled into a []map[string]interface{}
	//
	// Encode back to JSON, then Unmarshal into []ForemanOverrideValue for
	// the results
	results := []ForemanOverrideValue{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	jsonDecErr := json.Unmarshal(resultsBytes, &results)
	if jsonDecErr != nil {
		return queryResponse, jsonDecErr
	}
	// convert the search results from []ForemanOverrideValue to []interface
	// and set the search results on the query
	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}

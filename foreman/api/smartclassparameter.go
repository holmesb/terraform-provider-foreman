package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HanseMerkur/terraform-provider-utils/log"
)

const (
	ParameterEndpointPrefix = "/%s/%d/parameters"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

// ForemanSmartClassOverrideValue reprents overriden Smart Class Parameters as
// stored in Foreman
type ForemanSmartClassOverrideValue struct {
	Match string `json:"match"`
	Value string `json:"value"`
	Omit  bool   `json:"omit"`
}

// The ForemanSmartClassParameter API model represents a Smart Class Parameter
// defined by a Puppet Installation. It can only be read and updated by
// Foreman. Neither Deleted or Created
type ForemanSmartClassParameter struct {
	ForemanObject

	HostID        int
	HostGroupID   int
	EnvironmentID int

	ParameterID int

	Override           bool                             `json:"override"`
	Description        string                           `json:"description"`
	DefaultValue       string                           `json:"default_value"`
	HiddenValue        bool                             `json:"hidden_value"`
	Omit               bool                             `json:"omit"`
	Path               string                           `json:"path"`
	ValidatorType      string                           `json:"validator_type"`
	ValidatorRule      string                           `json:"validator_rule"`
	OverrideValues     []ForemanSmartClassOverrideValue `json:"override_values"`
	OverrideValueOrder string                           `json:"override_value_order"`
	ParameterType      string                           `json:"parameter_type"`
	Required           bool                             `json:"required"`
	MergeOverrides     bool                             `json:"merge_overrides"`
	MergeDefault       bool                             `json:"merge_default"`
	AvoidDuplicates    bool                             `json:"avoid_duplicates"`
}

// -----------------------------------------------------------------------------
// Read/Update Implementation
// -----------------------------------------------------------------------------

// ReadParameter reads the attributes of a ForemanSmartClassParameter identified by the
// supplied ID and returns a ForemanSmartClassParameter reference.
func (c *Client) ReadParameter(d *ForemanSmartClassParameter, id int) (*ForemanSmartClassParameter, error) {
	log.Tracef("foreman/api/parameter.go#Read")

	selEndA, selEndB := d.apiEndpoint()
	reqEndpoint := fmt.Sprintf(ParameterEndpointPrefix+"/%d", selEndA, selEndB, id)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readParameter ForemanSmartClassParameter
	sendErr := c.SendAndParse(req, &readParameter)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("readParameter: [%+v]", readParameter)

	d.Id = readParameter.Id
	d.Parameter = readParameter.Parameter
	return d, nil
}

// UpdateParameter deletes all parameters for the subject resource and re-creates them
// as we look at them differently on either side this is the safest way to reach sync
func (c *Client) UpdateParameter(d *ForemanSmartClassParameter, id int) (*ForemanSmartClassParameter, error) {
	log.Tracef("foreman/api/parameter.go#Update")

	selEndA, selEndB := d.apiEndpoint()
	reqEndpoint := fmt.Sprintf(ParameterEndpointPrefix+"/%d", selEndA, selEndB, id)

	parameterJSONBytes, jsonEncErr := json.Marshal(d)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("parameterJSONBytes: [%s]", parameterJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPut,
		reqEndpoint,
		bytes.NewBuffer(parameterJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var updatedParameter ForemanSmartClassParameter
	sendErr := c.SendAndParse(req, &updatedParameter)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("updatedParameter: [%+v]", updatedParameter)

	d.Id = updatedParameter.Id
	d.Parameter = updatedParameter.Parameter
	return d, nil
}

// DeleteParameter deletes the ForemanSmartClassParameters for the given resource
func (c *Client) DeleteParameter(d *ForemanSmartClassParameter, id int) error {
	log.Tracef("foreman/api/parameter.go#Delete")

	selEndA, selEndB := d.apiEndpoint()
	reqEndpoint := fmt.Sprintf(ParameterEndpointPrefix+"/%d", selEndA, selEndB, id)

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

// QueryParameter queries for a ForemanSmartClassParameter based on the attributes of the
// supplied ForemanSmartClassParameter reference and returns a QueryResponse struct
// containing query/response metadata and the matching parameters.
func (c *Client) QueryParameter(d *ForemanSmartClassParameter) (QueryResponse, error) {
	log.Tracef("foreman/api/parameter.go#Search")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", ParameterEndpointPrefix)
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
	name := `"` + d.Name + `"`
	reqQuery.Set("search", "name="+name)

	req.URL.RawQuery = reqQuery.Encode()
	sendErr := c.SendAndParse(req, &queryResponse)
	if sendErr != nil {
		return queryResponse, sendErr
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	// Results will be Unmarshaled into a []map[string]interface{}
	//
	// Encode back to JSON, then Unmarshal into []ForemanSmartClassParameter for
	// the results
	results := []ForemanSmartClassParameter{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	jsonDecErr := json.Unmarshal(resultsBytes, &results)
	if jsonDecErr != nil {
		return queryResponse, jsonDecErr
	}
	// convert the search results from []ForemanSmartClassParameter to []interface
	// and set the search results on the query
	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}

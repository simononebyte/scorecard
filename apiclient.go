package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// APIClient ...
type APIClient struct {
	authToken string
	client    *http.Client
}

// NewAPICLient returns a new APIClient
func NewAPICLient(token string) *APIClient {
	api := APIClient{}
	api.authToken = base64.StdEncoding.EncodeToString([]byte(token))
	api.client = &http.Client{
		Timeout: time.Second * 30,
	}
	return &api
}

// Get performs the requested API Get returning the results as JSON
func (api *APIClient) Get(url string, out interface{}) error {

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return reqErr
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", api.authToken))
	req.Header.Set("Content-Type", "application/json")

	if err := api.doJSONRequest(req, out); err != nil {
		return err
	}

	res, getErr := api.client.Do(req)
	if getErr != nil {
		return getErr
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status: %s", res.Status)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	jsonErr := json.Unmarshal(body, out)
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

// Post performs the requested API Post returning the results as JSON
func (api *APIClient) Post(url string, query interface{}, out interface{}) error {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(query)

	req, reqErr := http.NewRequest(http.MethodPost, url, b)
	if reqErr != nil {
		return reqErr
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", api.authToken))
	req.Header.Set("Content-Type", "application/json")

	if err := api.doJSONRequest(req, out); err != nil {
		return err
	}

	res, getErr := api.client.Do(req)
	if getErr != nil {
		return getErr
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status: %s", res.Status)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	jsonErr := json.Unmarshal(body, out)
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

func (api *APIClient) doJSONRequest(req *http.Request, out interface{}) error {

	res, getErr := api.client.Do(req)
	if getErr != nil {
		return getErr
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status: %s", res.Status)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	jsonErr := json.Unmarshal(body, out)
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

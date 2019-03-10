// ScreenshotLayer API Library for Go
// https://screenshotlayer.com/documentation
// This is a pre-release version and is subject to change

// Copyright 2019 Bryan Anderson (https://www.bganderson.com)
// Relesed under a BSD-style license which can be found in the LICENSE file

package screenshotlayergo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

const (
	// Version is the global package version
	Version = "0.1.0"
	// DefaultAPIEndpoint is the default screenshotlayer REST API endpoint
	DefaultAPIEndpoint = "http://api.screenshotlayer.com/api/capture"
)

// Client API connection settings
type Client struct {
	AccessKey string // AccessKey for API authentication
	URL       string // URL endpoint for API
	HTTPS     bool   // Use https endpoint
}

// APIRequest to screenshotlayer API
type APIRequest struct {
	URL            string `param:"url"`
	FullPage       string `param:"fullpage"`
	Width          string `param:"width"`
	Viewport       string `param:"viewport"`
	Format         string `param:"format"`
	SecretKey      string `param:"secret_key"`
	CSSURL         string `param:"css_url"`
	Delay          string `param:"delay"`
	TTL            string `param:"ttl"`
	Force          string `param:"force"`
	Placeholder    string `param:"placeholder"`
	UserAgent      string `param:"user_agent"`
	AcceptLanguage string `param:"accept_lang"`
	Export         string `param:"export"`
}

// APIResponse from screenshotlayer API
type APIResponse struct {
	Bytes    []byte
	APIError APIError
}

// APIError is returned if the API returns an error
type APIError struct {
	Success bool `json:"success"`
	Error   struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
}

// Screenshot queries the screenshotlayer.com API and returns APIResponse
func (c *Client) Screenshot(req *APIRequest) (*APIResponse, error) {
	r, err := c.queryAPI(req)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}
	if !r.APIError.Success {
		return nil, fmt.Errorf(
			"[ERROR] API responded with %d (%s). %s",
			r.APIError.Error.Code,
			r.APIError.Error.Type,
			r.APIError.Error.Info,
		)
	}
	return r, nil
}

// Build URL for screenshotlayer.com API query
func (c *Client) buildURL(req *APIRequest) string {
	if c.URL == "" {
		c.URL = DefaultAPIEndpoint
	}
	u, err := url.Parse(c.URL)
	if err != nil {
		panic("unable to parse API URL")
	}
	if c.HTTPS {
		u.Scheme = "https"
	}
	param := url.Values{}
	param.Add("access_key", c.AccessKey)
	/*
		Inefficient but allows me iterate over the struct fields.
	*/
	t, v := reflect.TypeOf(*req), reflect.ValueOf(req)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("param")
		if v := reflect.Indirect(v).FieldByName(f.Name).String(); v != "" {
			param.Add(tag, v)
		}
	}
	u.RawQuery = param.Encode()
	return u.String()
}

// Query screenshotlayer.com API
func (c *Client) queryAPI(req *APIRequest) (*APIResponse, error) {
	r, err := http.Get(c.buildURL(req))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var resp APIResponse
	resp.APIError.Success = true
	/*
		This is the most reliable way I could determine whether or not the API
		returned an error.
	*/
	if r.Header.Get("Content-Type") == "application/json; Charset=UTF-8" {
		if err := json.Unmarshal(body, &resp.APIError); err != nil {
			return nil, err
		}
		return &resp, nil
	}
	resp.Bytes = body
	return &resp, nil
}

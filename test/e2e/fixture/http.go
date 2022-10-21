package fixture

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/argoproj/argo-cd/v2/common"
)

// DoHttpRequest executes a http request against the Argo CD API server
func DoHttpRequest(method string, path string, data ...byte) (*http.Response, error) {
	reqUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	reqUrl.Scheme = "http"
	reqUrl.Host = apiServerAddress
	var body io.Reader
	if data != nil {
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, reqUrl.String(), body)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Http Client with TLS verification set to false.
	// Required when dealing with self-signed certificates.
	hc := &http.Client{Transport: tr}

	req.AddCookie(&http.Cookie{Name: common.AuthCookieName, Value: token})
	return hc.Do(req)
}

// DoHttpJsonRequest executes a http request against the Argo CD API server and unmarshals the response body as JSON
func DoHttpJsonRequest(method string, path string, result interface{}, data ...byte) error {
	resp, err := DoHttpRequest(method, path, data...)
	if err != nil {
		return err
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(responseData, result)
}

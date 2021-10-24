package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type graphErrors []struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}

func (e graphErrors) Error() string {
	b := strings.Builder{}
	for _, err := range e {
		b.WriteString(fmt.Sprintf("Message: %s, Locations: %+v", err.Message, err.Locations))
	}
	return b.String()
}

type GraphQLClient struct {
	Client  *http.Client
	BaseURL string
}

func (c GraphQLClient) Do(query string, variables map[string]interface{}, v interface{}) (string, error) {
	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return "", err
	}

	resp, err := c.Client.Post(c.BaseURL, "application/json", &buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}
	var out struct {
		Data   *json.RawMessage
		Errors graphErrors
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &out)
	if err != nil {
		return "", err
	}

	if out.Data != nil && v != nil {
		err := json.Unmarshal(*out.Data, &v)
		if err != nil {
			return string(body), err
		}
	}

	if len(out.Errors) > 0 {
		b, _ := json.Marshal(out.Errors)
		return string(b), out.Errors
	}

	return string(body), nil
}

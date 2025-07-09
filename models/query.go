package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nakurai/gourl/db"
	"gorm.io/gorm"
)

var client = &http.Client{}

type Query struct {
	ID        uint    `gorm:"primaryKey"`
	Data      JSONMap `gorm:"type:json"`
	Header    JSONMap `gorm:"type:json"`
	IsJson    bool
	Method    string
	Name      string // if the query is a saved query. For example: demo/post/message
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Query) Send() (string, error) {
	// just in case
	q.Method = strings.ToUpper(q.Method)
	var body io.Reader
	urlToUse := q.Url
	isPost := q.Method == http.MethodConnect
	isPut := q.Method == http.MethodPut
	isPatch := q.Method == http.MethodPatch

	if len(q.Data) > 0 {
		var err error
		if isPost || isPut || isPatch {
			// this required to add the parameters in a body
			if q.IsJson {
				body, err = q.GetJsonParam()
				if err != nil {
					return "", err
				}
			} else {
				body, err = q.GetFormParam()
				if err != nil {
					return "", err
				}
			}
		} else {
			// otherwise the parameters are part of the URL
			urlToUse, err = q.GetQueryUrl()
			if err != nil {
				return "", err
			}
		}
	}

	req, err := http.NewRequest(q.Method, urlToUse, body)
	if err != nil {
		return "", err
	}
	for headerKey, headerValue := range q.Header {
		req.Header.Set(headerKey, headerValue)
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(resBody)
	headerJson, err := json.Marshal(res.Header)
	if err != nil {
		return "", err
	}
	headerString := string(headerJson)
	finalString := fmt.Sprintf("%s\n\nbody:\n%s\n\nheaders:\n%s\n", res.Status, bodyString, headerString)
	return finalString, nil
}

// this function is used to append the query parameters to the
// the url.
// It returns the full url that should be used in the http request
func (q Query) GetQueryUrl() (string, error) {
	u, err := url.Parse(q.Url)
	if err != nil {
		return "", err
	}
	urlQuery := u.Query()
	for queryKey, queryValue := range q.Data {
		urlQuery.Set(queryKey, queryValue)
	}
	u.RawQuery = urlQuery.Encode()
	return u.String(), nil
}

// create the form encoded body that needs to be sent in the http request
func (q Query) GetFormParam() (io.Reader, error) {
	form := url.Values{}
	for k, v := range q.Data {
		form.Set(k, v)
	}
	return strings.NewReader(form.Encode()), nil
}

// create the JSON encoded body that needs to be sent in the http request
func (q Query) GetJsonParam() (io.Reader, error) {
	paramBytes, err := json.Marshal(q.Data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(paramBytes), nil
}

func (q *Query) Save() error {
	var existingQuery Query
	if q.Name == "" {
		return fmt.Errorf("failed to save query. The query does not have a name")
	}
	res := db.Db.Where("name = ?", q.Name).First(&existingQuery)
	if res.Error != nil {
		// if the record was not found, we can create the new query
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			res := db.Db.Create(q)
			if res.Error != nil {
				return fmt.Errorf("error while saving the new query %s: %v", q.Name, res.Error)
			}
			return nil
		} else {
			return fmt.Errorf("error while fetching existing query %s: %v", q.Name, res.Error)
		}
	}
	if res.RowsAffected > 0 {
		return fmt.Errorf("EXIST-ALREADY")
	}
	return nil
}

func BuildQueryTree() error {
	queries := []Query{}
	res := db.Db.Order("name").Find(&queries)
	if res.Error != nil {
		return fmt.Errorf("error while querying all queries: %v", res.Error)
	}

	for _, query := range queries {
		nameParts := strings.Split(query.Name, "/")
		QueryTree.AddQueryName(nameParts, query.Method)
	}

	return nil
}

package tdr

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/philips-software/go-hsdp-api/fhir"
)

// DataItemsService provides operations on TDR data items
type DataItemsService struct {
	client *Client
}

// GetDataItemOptions describes the fields on which you can search for data items
type GetDataItemOptions struct {
	Organization *string `url:"organization,omitempty"`
	DataType     *string `url:"dataType,omitempty"`
	Count        *int    `url:"_count,omitempty"`
}

// KeyValue is backed by a string hash map
type KeyValue map[string]string

// DataSearch builds a custom query for TDR searches on the data element
func DataSearch(kv KeyValue) OptionFunc {
	return func(req *http.Request) error {
		q := url.Values{}
		for k, v := range kv {
			q.Add(k, v)
		}
		custom := q.Encode()
		if req.URL.RawQuery != "" && custom != "" {
			custom = "&" + custom
		}
		req.URL.RawQuery += custom
		return nil
	}
}

// GetDataItem searches for data items in TDR
// Use the DataSearch OptionFunc to search in the data part. When using this the
// DataType must added as part of the options
func (d *DataItemsService) GetDataItem(opt *GetDataItemOptions, options ...OptionFunc) ([]*DataItem, *Response, error) {
	var dataItems []*DataItem

	req, err := d.client.newTDRRequest("GET", "store/tdr/DataItem", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Api-Version", TDRAPIVersion)

	var bundleResponse fhir.Bundle

	resp, err := d.client.Do(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	if bundleResponse.Total == 0 {
		return dataItems, resp, ErrEmptyResult
	}
	for _, e := range bundleResponse.Entry {
		item := new(DataItem)
		if err := json.Unmarshal(e.Resource, item); err == nil {
			dataItems = append(dataItems, item)
		} else {
			return nil, resp, err
		}
	}
	return dataItems, resp, err
}

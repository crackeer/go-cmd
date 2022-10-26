package http

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// GetSourceList
//  @param workCode
//  @return string
//  @return error
func GetSourceList() ([]map[string]interface{}, error) {

	client := resty.New()
	resp, err := client.R().Get("http://i.shepherd.realsee.com/util/source/list")

	if err != nil {
		return nil, err
	}

	data := gjson.Get(resp.String(), "data").String()
	list := []map[string]interface{}{}
	err = json.Unmarshal([]byte(data), &list)
	return list, err
}

package http

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// DecodeWorkCode
//  @param workCode
//  @return string
//  @return error
func DecodeWorkCode(workCode string) (string, error) {

	client := resty.New()
	resp, err := client.R().SetQueryParams(map[string]string{
		"work_code": workCode,
	}).Get("http://i.svc.open.realsee.com/util/decode/work_code")

	if err != nil {
		return "", fmt.Errorf("Decode失败(error=%s)", err.Error())
	}

	gr := gjson.Get(resp.String(), "data.work_id")
	if gr.Exists() {
		return gr.String(), nil
	}

	return "", errors.New("未解码成功")
}

// EncodeWorkID
//  @param workID
//  @return string
//  @return error
func EncodeWorkID(workID string) (string, error) {

	client := resty.New()
	resp, err := client.R().SetQueryParams(map[string]string{
		"work_id": workID,
	}).Get("http://i.svc.open.realsee.com/util/decode/work_code")

	if err != nil {
		return "", fmt.Errorf("Encode失败(error=%s)", err.Error())
	}

	gr := gjson.Get(resp.String(), "data.work_code")
	if gr.Exists() {
		return gr.String(), nil
	}

	return "", errors.New("未Encode成功")
}

// QueryWorkData
//  @param workID
//  @param table
//  @return string
//  @return error
func QueryWorkData(query map[string]interface{}, table string) ([]map[string]interface{}, error) {
	client := resty.New()
	resp, err := client.R().SetBody(map[string]interface{}{
		"table":      table,
		"query":      query,
		"order_by":   "id desc",
		"page":       1,
		"page_size":  100,
		"query_type": "list",
	}).Post("http://i.svc.open.realsee.com/util/database/vrapi")

	if err != nil {
		return nil, fmt.Errorf("QueryError(error=%s)", err.Error())
	}

	gr := gjson.Get(resp.String(), "data.list")
	if gr.Exists() {
		retData := []map[string]interface{}{}
		json.Unmarshal([]byte(gr.String()), &retData)
		return retData, nil
	}

	return nil, nil
}

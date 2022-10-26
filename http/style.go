package http

import (
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// ForceUpdateStyleStatus
//  @param styleID
//  @param status
//  @return int64
//  @return error
func ForceUpdateStyleStatus(styleID int64, status int64) (int64, error) {
	client := resty.New()
	result, err := client.R().SetBody(map[string]interface{}{
		"style_id": styleID,
		"data": map[string]interface{}{
			"status": status,
		},
	}).SetHeader("Content-Type", "application/json").SetHeader("X-Realsee-UserID", "888888").Post("http://i.sems-model.api.realsee.com/style/v1/force_modify.json")

	if err != nil {
		return 0, err
	}

	if gjson.Get(result.String(), "code").Int() != 0 {
		return 0, errors.New(gjson.Get(result.String(), "status").String())
	}

	return gjson.Get(result.String(), "data.result").Int(), nil

}

package http

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type CosData struct {
}

// GetRenderData
//  @param cosKey
//  @param keys
//  @return string
//  @return error
func GetRenderData(cosKey string, keys []string) (string, error) {

	client := resty.New()
	resp, err := client.R().SetQueryParams(map[string]string{
		"cos_keys": cosKey,
		"keys":     strings.Join(keys, ","),
	}).Get("http://i.shepherd.realsee.com/util/render/cos-get")

	if err != nil {
		return "", fmt.Errorf("获取数据失败(error=%s)", err.Error())
	}

	gr := gjson.Get(resp.String(), "data.data")
	if gr.Exists() {
		return gr.String(), nil
	}

	return "", errors.New("未解码成功")
}

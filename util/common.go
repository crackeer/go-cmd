package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/speps/go-hashids"
)

// GetStringValFromMap
//
//	@param container
//	@param key
//	@return string
func GetStringValFromMap(container map[string]interface{}, key string) string {
	val, ok := container[key]

	if !ok {
		return ""
	}

	if v, ok := val.(json.Number); ok {
		val, _ := v.Int64()
		return fmt.Sprintf("%d", val)
	}

	if v, ok := val.(int64); ok {
		return fmt.Sprintf("%d", v)
	}

	if v, ok := val.(int); ok {
		return fmt.Sprintf("%d", v)
	}

	if v, ok := val.(float64); ok {
		return strconv.FormatFloat(v, 'f', -1, 64)
	}

	if v, ok := val.(float32); ok {
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	if v, ok := val.(string); ok {
		return v
	}

	bytes, _ := json.Marshal(val)
	return string(bytes)
}

// DiffMapSlice
//
//	@param oldFields
//	@param newFields
//	@param keys
//	@return []map
func DiffMapSlice(oldFields []map[string]interface{}, newFields []map[string]interface{}, keys []string) []map[string]interface{} {
	mapping := map[string]map[string]interface{}{}
	for _, item := range oldFields {
		values := []string{}
		for _, key := range keys {
			values = append(values, GetStringValFromMap(item, key))
		}
		mapping[strings.Join(values, "_")] = item
	}

	retData := []map[string]interface{}{}
	for _, item := range newFields {
		values := []string{}
		for _, key := range keys {
			values = append(values, GetStringValFromMap(item, key))
		}
		if _, ok := mapping[strings.Join(values, "_")]; !ok {
			retData = append(retData, item)
		}
	}

	return retData
}

// DeleteMapSliceKeys
//
//	@param list
//	@param keys
//	@return []map
func DeleteMapSliceKeys(list []map[string]interface{}, keys ...string) []map[string]interface{} {
	retData := []map[string]interface{}{}
	for _, item := range list {
		for _, key := range keys {
			delete(item, key)
		}
		retData = append(retData, item)
	}

	return retData
}

// GetInt64ValFromMap
//
//	@param container
//	@param key
//	@return int64
func GetInt64ValFromMap(container map[string]interface{}, key string) int64 {

	val, ok := container[key]

	if !ok {
		return 0
	}

	if v, ok := val.(int64); ok {
		return v
	}

	if v, ok := val.(json.Number); ok {
		val, _ := v.Int64()
		return val
	}

	if v, ok := val.(int); ok {
		return int64(v)
	}

	if v, ok := val.(uint32); ok {
		return int64(v)
	}

	if v, ok := val.(int32); ok {
		return int64(v)
	}

	if v, ok := val.(uint64); ok {
		return int64(v)
	}

	if v, ok := val.(float64); ok {
		return int64(v)
	}

	if v, ok := val.(string); ok {
		if rv, err := strconv.Atoi(v); err == nil {
			return int64(rv)
		}
	}
	return 0
}

// DecodeResourceCode
//
//	@param resourceCode
//	@return int64
//	@return int64
//	@return int64
//	@return error
func DecodeResourceCode(resourceCode string) (int64, int64, int64, error) {
	hd := hashids.NewData()
	hd.Salt = "#resource-parameter-hashids-salt@localized#"
	hd.MinLength = 18

	h, _ := hashids.NewWithData(hd)
	result, err := h.DecodeInt64WithError(resourceCode)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(result) < 3 {
		return 0, 0, 0, fmt.Errorf("invalid resource code")
	}

	return result[0], result[1], result[2], nil
}

// EncodeResourceCode
//
//	@param customID
//	@param teamID
//	@param resourceType
//	@return string
//	@return error
func EncodeResourceCode(customID, teamID, resourceType int64) (string, error) {
	hd := hashids.NewData()
	hd.Salt = "#resource-parameter-hashids-salt@localized#"
	hd.MinLength = 18

	h, _ := hashids.NewWithData(hd)

	return h.EncodeInt64([]int64{customID, teamID, resourceType})
}

func ToString(src interface{}) string {
	return asString(src)
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func GetInnerIP() []string {
	ipList := []string{}
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return ipList
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("Error fetching addresses for %q: %v\n", iface.Name, err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.IsLoopback() || ip.IsLinkLocalUnicast() || !ip.IsGlobalUnicast() {
				// 忽略回环地址、链路本地地址和非全局单播地址
				continue
			}

			if ip.To4() != nil && ip.IsPrivate() {
				ipList = append(ipList, ip.String())
			}
			// 注意：如果你也需要IPv6私有地址，可以取消上面的IP版本检查
		}
	}
	return ipList
}

// DownloadTo
//
//	@param urlString
//	@param target
//	@return error
func DownloadTo(urlString string, target string) error {
	dir, _ := filepath.Split(target)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	response, err := http.Get(urlString)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(target, bytes, os.ModePerm)
}

// ExtractSomeString
//
//	@param value
//	@param matchFunc
//	@return []string
func Extract(value interface{}, matchFunc func(string) bool) []string {
	if strValue, ok := value.(string); ok {
		if data, err := JsonDecodeX(strValue); err == nil {
			return Extract(data, matchFunc)
		}
		if matchFunc(strValue) {
			return []string{strValue}
		}
		return []string{}
	}

	var retData []string
	if mapValue, ok := value.(map[string]interface{}); ok {
		for _, value := range mapValue {
			retData = append(retData, Extract(value, matchFunc)...)
		}
		return retData
	}

	if listValue, ok := value.([]interface{}); ok {
		for _, value := range listValue {
			retData = append(retData, Extract(value, matchFunc)...)
		}
		return retData
	}
	return nil
}

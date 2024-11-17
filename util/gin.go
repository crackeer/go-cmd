package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultMultipartMemory = 32 << 20

// OTAFailure
//
//	@param ctx
//	@param code
//	@param message
func OTAFailure(ctx *gin.Context, code int64, message string) {
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"Success": false,
		"Error": map[string]interface{}{
			"Code":   code,
			"String": message,
		},
	})
	ctx.Abort()
}

// OTASuccess
//
//	@param ctx
//	@param data
func OTASuccess(ctx *gin.Context, data interface{}) {
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"Success": true,
		"Data":    data,
	})
	ctx.Abort()
}

// QueryAllParams
//
//	@param ctx
//	@return map
func QueryAllParams(ctx *gin.Context) map[string]string {

	if ctx == nil {
		return map[string]string{}
	}

	query := ctx.Request.URL.Query()
	retData := make(map[string]string, len(query))
	for k, v := range query {
		retData[k] = v[0]
	}
	return retData
}

// AllGetParams ...
//
//	@param ctx
//	@return map
func AllGetParams(ctx *gin.Context) map[string]string {

	if ctx == nil {
		return map[string]string{}
	}

	querys := ctx.Request.URL.Query()
	retData := make(map[string]string, len(querys))
	for k, v := range querys {
		retData[k] = strings.Join(v, ",")
	}
	return retData
}

// AllPostParams  ...
//
//	@param ctx
//	@return map
func AllPostParams(ctx *gin.Context) map[string]interface{} {

	if ctx == nil {
		return map[string]interface{}{}
	}

	retData := map[string]interface{}{}
	contentType := ctx.ContentType()
	switch contentType {
	case gin.MIMEPOSTForm:
		ctx.Request.ParseForm()
		for k, v := range ctx.Request.Form {
			retData[k] = strings.Join(v, ",")
		}
	case gin.MIMEMultipartPOSTForm:
		if err := ctx.Request.ParseMultipartForm(defaultMultipartMemory); err == nil {
			for k, v := range ctx.Request.MultipartForm.Value {
				retData[k] = strings.Join(v, ",")
			}
			for k, v := range ctx.Request.MultipartForm.File {
				retData[k] = v
			}
		}
	default:
		if raw, err := ctx.GetRawData(); err == nil {
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(raw)) // 再写入body
			jsonDecoder := json.NewDecoder(bytes.NewReader(raw))
			jsonDecoder.UseNumber()
			if err := jsonDecoder.Decode(&retData); err != nil {
				retData["raw_data"] = raw
			}
		}
	}

	return retData
}

// AllHeader
//
//	@param ctx
//	@return map
func AllHeader(ctx *gin.Context) map[string]string {
	if ctx == nil {
		return map[string]string{}
	}

	retData := make(map[string]string)
	for k, v := range ctx.Request.Header {
		retData[k] = strings.Join(v, ",")
	}
	return retData
}

// Query
//
//	@param ctx
//	@param arge
//	@return string
func Query(ctx *gin.Context, arge ...string) string {
	for _, k := range arge {
		if val, e := ctx.GetQuery(k); e {
			return val
		}
	}
	return ""
}

// AllParams
//
//	@param ctx
//	@return map
func AllParams(ctx *gin.Context) map[string]interface{} {

	retData := map[string]interface{}{}
	for k, v := range AllGetParams(ctx) {
		retData[k] = v
	}

	if http.MethodPost == ctx.Request.Method {
		postData := AllPostParams(ctx)
		for k, v := range postData {
			retData[k] = v
		}
	}
	return retData
}

// Success  ...
//
//	@param ctx
//	@param data
func Success(ctx *gin.Context, data interface{}) {
	SetResponse(ctx, CodeSuccess, data, MessageSuccess)
}

// Failure ...
//
//	@param ctx
//	@param code
//	@param message
func Failure(ctx *gin.Context, code int64, message string) {
	SetResponse(ctx, code, nil, message)
}

// SetResponse ...
//
//	@param ctx
//	@param code
//	@param data
//	@param message
func SetResponse(ctx *gin.Context, code int64, data interface{}, message string) {
	if response, exists := GetJSONResponse(ctx); exists {
		response.Code = code
		response.Message = message
		response.Data = data
	} else {
		ctx.Set(keyResponse, &JSONResponse{
			Code:    code,
			Message: message,
			Data:    data,
		})
	}
}

// JSONResspone
type JSONResponse struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const keyResponse = "__RESPONSE__"

const (

	// CodeSuccess
	CodeSuccess int64 = 0

	CodeDefaultError int64 = -1
)

const (

	// MessageSuccess
	MessageSuccess = "success"

	// MessageEmptyResponse
	MessageEmptyResponse = "no router match or no data response"
)

// getJSONResponse
//
//	@param ctx
//	@return *JSONResponse
//	@return bool
func GetJSONResponse(ctx *gin.Context) (*JSONResponse, bool) {
	if body, exists := ctx.Get(keyResponse); exists {
		response, flag := body.(*JSONResponse)
		return response, flag
	}
	return nil, false
}

// ResponseJSON ...
//
//	@param ctx N
func ResponseJSON(ctx *gin.Context) {
	if ctx.IsAborted() {
		return
	}
	if response, found := GetJSONResponse(ctx); found {
		ctx.Abort()
		ctx.PureJSON(http.StatusOK, response)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, &JSONResponse{
		Code:    CodeDefaultError,
		Message: MessageEmptyResponse,
		Data:    nil,
	})
}

// DoResponseJSON
//
//	@return gin.HandlerFunc
func DoResponseJSON() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		ResponseJSON(ctx)
	}
}

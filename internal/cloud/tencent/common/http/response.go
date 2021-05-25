package http

import (
	"encoding/json"
	"fmt"
	errors2 "github.com/eadydb/k8s-aim/internal/cloud/tencent/common/errors"
	"io/ioutil"
	"net/http"
)

type Response interface {
	ParseErrorFromHTTPResponse(body []byte) error
}

type BaseResponse struct {
}

type ErrorResponse struct {
	Response struct {
		Error struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error,omitempty"`
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

type DeprecatedAPIErrorResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CodeDesc string `json:"codeDesc"`
}

func (r *BaseResponse) ParseErrorFromHTTPResponse(body []byte) (err error) {
	resp := &ErrorResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
		return errors2.NewTencentCloudSDKError("ClientError.ParseJsonError", msg, "")
	}
	if resp.Response.Error.Code != "" {
		return errors2.NewTencentCloudSDKError(resp.Response.Error.Code, resp.Response.Error.Message, resp.Response.RequestId)
	}

	deprecated := &DeprecatedAPIErrorResponse{}
	err = json.Unmarshal(body, deprecated)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
		return errors2.NewTencentCloudSDKError("ClientError.ParseJsonError", msg, "")
	}
	if deprecated.Code != 0 {
		return errors2.NewTencentCloudSDKError(deprecated.CodeDesc, deprecated.Message, "")
	}
	return nil
}

func ParseFromHttpResponse(hr *http.Response, response Response) (err error) {
	defer hr.Body.Close()
	body, err := ioutil.ReadAll(hr.Body)
	if err != nil {
		msg := fmt.Sprintf("Fail to read response body because %s", err)
		return errors2.NewTencentCloudSDKError("ClientError.IOError", msg, "")
	}
	if hr.StatusCode != 200 {
		msg := fmt.Sprintf("Request fail with http status code: %s, with body: %s", hr.Status, body)
		return errors2.NewTencentCloudSDKError("ClientError.HttpStatusCodeError", msg, "")
	}
	//log.Printf("[DEBUG] Response Body=%s", body)
	err = response.ParseErrorFromHTTPResponse(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json content: %s, because: %s", body, err)
		return errors2.NewTencentCloudSDKError("ClientError.ParseJsonError", msg, "")
	}
	return
}

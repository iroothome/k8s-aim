package cvm

import (
	"github.com/eadydb/k8s-aim/internal/cloud/tencent/common"
	tcHttp "github.com/eadydb/k8s-aim/internal/cloud/tencent/common/http"
)

const APIVersion = "2017-03-12"

// ZoneInfo 可用区信息
type ZoneInfo struct {
	Zone      string `json:"Zone,omitempty"` // 可用区名称
	ZoneName  string `json:"ZoneName"`       // 可用区描述
	ZoneId    string `json:"ZoneId"`         // 可用区ID
	ZoneState string `json:"ZoneState"`      // 可用区状态，包含AVAILABLE和UNAVAILABLE。AVAILABLE代表可用，UNAVAILABLE代表不可用。
}

// RegionInfo 地域信息
type RegionInfo struct {
	Region      string `json:"Region"`      // 地域名称
	RegionName  string `json:"RegionName"`  // 地域描述
	RegionState string `json:"RegionState"` // 地域是否可用状态,包含AVAILABLE和UNAVAILABLE。AVAILABLE代表可用，UNAVAILABLE代表不可用。
}

// Client 客户端
type Client struct {
	common.Client
}

// ZonesRequest 可用区
type ZonesRequest struct {
	*tcHttp.BaseRequest
}

// ZonesResponse 可用区响应结果
type ZonesResponse struct {
	*tcHttp.BaseResponse
	Response *struct {
		TotalCount uint64      `json:"TotalCount,omitempty"` // 可用区数量
		ZoneSet    []*ZoneInfo `json:"ZoneSet,omitempty"`    // 可用区列表
		RequestId  string      `json:"RequestId,omitempty"`  // 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId
	} `json:"Response"`
}

// NewZonesRequest 实例化
func NewZonesRequest() *ZonesRequest {
	z := &ZonesRequest{BaseRequest: &tcHttp.BaseRequest{}}
	z.Init().WithApiInfo("cvm", APIVersion, "DescribeZones")
	return z
}

// NewZonesResponse 实例化
func NewZonesResponse() *ZonesResponse {
	response := &ZonesResponse{BaseResponse: &tcHttp.BaseResponse{}}
	return response
}

// DescribeZones 查询可用区信息
func (c *Client) DescribeZones(req *ZonesRequest) (*ZonesResponse, error) {
	if req == nil {
		req = NewZonesRequest()
	}
	resp := NewZonesResponse()
	err := c.Send(req, resp)
	return resp, err
}

// RegionRequest 地域请求参数
type RegionRequest struct {
	*tcHttp.BaseRequest
}

// RegionResponse 地域响应结果
type RegionResponse struct {
	*tcHttp.BaseResponse
	Response *struct {
		TotalCount uint64        `json:"TotalCount,omitempty"` // 可用区数量
		RegionSet  []*RegionInfo `json:"RegionSet,omitempty"`  // 地域列表
		RequestId  string        `json:"RequestId,omitempty"`  // 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId
	} `json:"Response"`
}

// NewRegionRequest 例化
func NewRegionRequest() *RegionRequest {
	req := &RegionRequest{BaseRequest: &tcHttp.BaseRequest{}}
	req.Init().WithApiInfo("cvm", APIVersion, "DescribeRegions")
	return req
}

// NewRegionResponse 实例化
func NewRegionResponse() *RegionResponse {
	resp := &RegionResponse{BaseResponse: &tcHttp.BaseResponse{}}
	return resp
}

// DescribeRegions 查询地域信息
func (c *Client) DescribeRegions(req *RegionRequest) (*RegionResponse, error) {
	if req == nil {
		req = NewRegionRequest()
	}
	resp := NewRegionResponse()
	err := c.Send(req, resp)
	return resp, err
}

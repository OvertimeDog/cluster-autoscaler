/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vpc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/baiducloud/baiducloud-sdk-go/bce"
)

// RouteRule define route
type RouteRule struct {
	RouteRuleID        string `json:"routeRuleId"`
	RouteTableID       string `json:"routeTableId"`
	SourceAddress      string `json:"sourceAddress"`
	DestinationAddress string `json:"destinationAddress"`
	NexthopID          string `json:"nexthopId"`
	NexthopType        string `json:"nexthopType"`
	Description        string `json:"description"`
}

// ListRouteArgs define listroute args
type ListRouteArgs struct {
	RouteTableID string `json:"routeTableId"`
	VpcID        string `json:"vpcId"`
}

// ListRouteResponse define response of list route
type ListRouteResponse struct {
	RouteTableID string      `json:"routeTableId"`
	VpcID        string      `json:"vpcId"`
	RouteRules   []RouteRule `json:"routeRules"`
}

func (args *ListRouteArgs) validate() error {
	if args == nil {
		return fmt.Errorf("ListRouteArgs need args")
	}
	if args.RouteTableID == "" && args.VpcID == "" {
		return fmt.Errorf("ListRouteArgs need RouteTableID or VpcID")
	}

	return nil
}

// ListRouteTable list all routes
func (c *Client) ListRouteTable(args *ListRouteArgs) ([]RouteRule, error) {
	err := args.validate()
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"routeTableId": args.RouteTableID,
		"vpcId":        args.VpcID,
	}
	req, err := bce.NewRequest("GET", c.GetURL("v1/route", params), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}
	var routesResp *ListRouteResponse
	err = json.Unmarshal(bodyContent, &routesResp)

	if err != nil {
		return nil, err
	}
	return routesResp.RouteRules, nil
}

// DeleteRoute delete a route
// http://gollum.baidu.com/Logical-Network-API#??????????????????
func (c *Client) DeleteRoute(routeID string) error {
	if routeID == "" {
		return fmt.Errorf("DeleteRoute need routeID")
	}
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	req, err := bce.NewRequest("DELETE", c.GetURL("v1/route/rule"+"/"+routeID, params), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// CreateRouteRuleArgs define args create route
// http://gollum.baidu.com/Logical-Network-API#??????????????????
type CreateRouteRuleArgs struct {
	RouteTableID  string `json:"routeTableId"`
	SourceAddress string `json:"sourceAddress"`
	// ??????????????????????????????0.0.0.0/0???
	// ????????????????????????????????????????????????VPC???????????????
	// ????????????????????????????????????????????????????????????????????????
	DestinationAddress string `json:"destinationAddress"`
	// ????????????????????????????????????0.0.0.0/0???
	// ??????????????????????????????VPC cidr??????
	// ?????????????????????VPC cidr???0.0.0.0/0????????????
	NexthopID string `json:"nexthopId,omitempty"`
	// ?????????id??????nexthopType???????????????????????????
	// ?????????????????????
	NexthopType string `json:"nexthopType"`
	// ???????????????Bcc?????????"custom"???
	// VPN?????????"vpn"???NAT?????????"nat"????????????????????????"defaultGateway"
	Description string `json:"description"`
}

// CreateRouteResponse define response of creating route
type CreateRouteResponse struct {
	RouteRuleID string `json:"routeRuleId"`
}

func (args *CreateRouteRuleArgs) validate() error {
	if args == nil {
		return fmt.Errorf("CreateRouteRuleArgs need args")
	}
	if args.RouteTableID == "" {
		return fmt.Errorf("CreateRouteRuleArgs need RouteTableID")
	}
	if args.SourceAddress == "" || args.DestinationAddress == "" {
		return fmt.Errorf("CreateRouteRuleArgs need address")
	}
	if args.NexthopID == "" || args.NexthopType == "" {
		return fmt.Errorf("CreateRouteRuleArgs need NexthopID and NexthopType")
	}
	return nil
}

// CreateRouteRule create a route rule
func (c *Client) CreateRouteRule(args *CreateRouteRuleArgs) (string, error) {
	err := args.validate()
	if err != nil {
		return "", err
	}
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/route/rule", params), bytes.NewBuffer(postContent))
	if err != nil {
		return "", err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return "", err
	}
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return "", err
	}
	var crResp *CreateRouteResponse
	err = json.Unmarshal(bodyContent, &crResp)

	if err != nil {
		return "", err
	}
	return crResp.RouteRuleID, nil
}

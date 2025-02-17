package webull

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	// "fmt"
	"github.com/google/uuid"
	"net/url"

	model "quantfu.com/webull/openapi"
)

// GetOrders returns orders.
func (c *Client) GetOrders(accountID string, status model.OrderStatus, count int32) ([]*model.GetOrdersItem, error) {
	var (
		u, _        = url.Parse(TradeEndpoint + "/v2/option/list")
		response    []model.GetOrdersItem
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	queryParams["secAccountId"] = accountID
	queryParams["startTime"] = "1970-0-1"
	queryParams["dateType"] = "ORDER"
	queryParams["pageSize"] = fmt.Sprintf("%d", count)
	queryParams["status"] = string(status)

	ords := make([]*model.GetOrdersItem, 0)
	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return ords, err
	}

	for _, item := range response {
		o := &model.GetOrdersItem{}
		*o = item
		ords = append(ords, o)
	}
	return ords, err
}

// IsTradeable returns information on where a specific ticker is traded
func (c *Client) IsTradeable(tickerID string) (*model.GetIsTradeableResponse, error) {
	var (
		u, _        = url.Parse(TradeEndpoint + "/ticker/broker/permissionV2")
		response    model.GetIsTradeableResponse
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	queryParams[QueryKeyTickerID] = tickerID

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// PlaceOrder places trade (TODO)
func (c *Client) PlaceOrder(accountID int64, input model.PostStockOrderRequest) (*model.PostOrderResponse, error) {
	var (
		u, _       = url.Parse(TradeEndpoint + "/order/" + strconv.FormatInt(accountID, 10) + "/placeStockOrder")
		headersMap = make(map[string]string)
		response   model.PostOrderResponse
	)

	if input.SerialId == nil || len(*input.SerialId) == 0 {
		input.SerialId = &c.UUID
	}

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = c.PostAndDecode(*u, &response, &headersMap, nil, payload)
	if err != nil {
		return &response, err
	}
	if response.OrderId == nil || len(*response.OrderId) == 0 {
		err = fmt.Errorf("Placed order not confirmed")
	}
	return &response, err
}

// CheckOtocoOrder checks OTOCO order (TODO)
func (c *Client) CheckOtocoOrder(accountID int64, input model.PostOtocoOrderRequest) (*interface{}, error) {
	var (
		u, _       = url.Parse(TradeEndpoint + "/v2/corder/stock/place/" + strconv.FormatInt(accountID, 10))
		headersMap = make(map[string]string)
		response   interface{}
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = c.PostAndDecode(*u, &response, &headersMap, nil, payload)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// PlaceOtocoOrder places OTOCO trade (TODO)
func (c *Client) PlaceOtocoOrder(accountID string, input model.PostOtocoOrderRequest) (*interface{}, error) {
	var (
		u, _       = url.Parse(TradeEndpoint + "/v2/corder/stock/place/" + accountID)
		headersMap = make(map[string]string)
		response   interface{}
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = c.PostAndDecode(*u, &response, &headersMap, nil, payload)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// CancelOrder cancels trade
func (c *Client) CancelOrder(accountID, orderID string) (*interface{}, error) {
	var (
		u, _       = url.Parse(TradeEndpoint + "/paper/1/acc/" + accountID + "/orderop/cancel/" + orderID)
		headersMap = make(map[string]string)
	)
	var response interface{}

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	err := c.PostAndDecode(*u, &response, &headersMap, nil, nil)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// ModifyOrder modifies trade (TODO)
func (c *Client) ModifyOrder(accountID string, orderID string, input model.PostStockOrderRequest) (*interface{}, error) {
	var (
		u, _       = url.Parse(TradeEndpoint + "/order/" + accountID + "/modifyStockOrder/" + orderID)
		headersMap = make(map[string]string)
	)
	var response interface{}

	if input.SerialId == nil || len(*input.SerialId) == 0 {
		input.SerialId = &c.UUID
	}

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = c.PostAndDecode(*u, &response, &headersMap, nil, payload)
	if err != nil {
		return &response, err
	}
	return &response, err
}

type GetOrdersRequest struct {
	DateType        string      `json:"dateType"`
	PageSize        int         `json:"pageSize"`
	StartTimeStr    string      `json:"startTimeStr"`
	EndTimeStr      string      `json:"endTimeStr"`
	Action          interface{} `json:"action,omitempty"`
	LastCreateTime0 int64       `json:"lastCreateTime0"`
	SecAccountID    int64       `json:"secAccountId"`
	Status          string      `json:"status"`
}

// GetOrdersV returns orders.
func (c *Client) GetOrdersV5(accountID int64, status model.OrderStatus, stTime time.Time, endTime time.Time, count int32) ([]*model.OrderItemV5, error) {
	var (
		u, _        = url.Parse(UsTradeEndpointV + "/order/list")
		response    []model.OrderItemV5
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	queryParams["secAccountId"] = strconv.FormatInt(accountID, 10)

	input := GetOrdersRequest{
		DateType:     "ORDER",
		PageSize:     int(count),
		SecAccountID: accountID,
		Status:       string(status),
	}

	if stTime.Year() > 2000 {
		input.StartTimeStr = stTime.Format("2006-01-02")
	} else {
		input.StartTimeStr = "2015-01-01"
	}
	if endTime.Year() > 2000 {
		input.LastCreateTime0 = endTime.Unix() * 1000
		input.EndTimeStr = endTime.Format("2006-01-02")
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = c.PostAndDecode(*u, &response, &headersMap, &queryParams, payload)
	if err != nil {
		return nil, err
	}

	rsFiltered := make([]*model.OrderItemV5, 0)
	if len(response) > 0 {
		// asking for all ?
		if strings.ToLower(input.Status) != strings.ToLower(string(model.ALL)) {
			// filter based on status
			for _, o := range response {
				if strings.ToLower(o.GetStatus()) == strings.ToLower(input.Status) {

					ord := &model.OrderItemV5{}
					*ord = o
					rsFiltered = append(rsFiltered, ord)
				}
			}
		} else {
			for _, o := range response {
				ord := &model.OrderItemV5{}
				*ord = o
				rsFiltered = append(rsFiltered, ord)
			}
		}
	}
	return rsFiltered, nil
}

// GetFilledOrdersByTicker returns orders.
func (c *Client) GetFilledOrdersByTicker(accountID int64, tickerId int64, lastFillTimeMs int64, count int32) ([]*model.OrderFill, error) {
	var (
		u, _        = url.Parse(UsTradeEndpointV + "/order/filledOrders")
		response    []model.OrderFill
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	queryParams["secAccountId"] = strconv.FormatInt(accountID, 10)
	queryParams["tickerId"] = strconv.FormatInt(tickerId, 10)
	queryParams["lastFilledTime"] = strconv.FormatInt(lastFillTimeMs, 10)
	queryParams["pageSize"] = strconv.FormatInt(int64(count), 10)

	fills := make([]*model.OrderFill, 0)
	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return fills, nil
	}

	for _, fill := range response {
		of := &model.OrderFill{}
		*of = fill
		fills = append(fills, of)
	}
	return fills, err
}

type CancelStOrderResponse struct {
	Result       bool   `json:"result"`
	OrderId      int64  `json:"orderId"`
	LastSerialId string `json:"lastSerialId"`
}

// Cancel order
func (c *Client) CancelOrderV5(accountID int64, orderId int64) (bool, error) {

	var (
		u, _        = url.Parse(UsTradeEndpointV + "/order/stockOrderCancel")
		response    CancelStOrderResponse
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	queryParams["secAccountId"] = strconv.FormatInt(accountID, 10)
	queryParams["serialId"] = c.UUID
	queryParams["orderId"] = fmt.Sprintf("%d", orderId)

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return false, err
	}

	return response.Result, nil
}

func (c *Client) PlaceOrderV5(accountID int64, input model.PostStockOrderRequest) (*model.PostOrderResponse, error) {
	var (
		u, _        = url.Parse(UsTradeEndpointV + "/order/stockOrderPlace")
		response    model.PostOrderResponse
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	queryParams["secAccountId"] = strconv.FormatInt(accountID, 10)

	if input.SerialId == nil || len(*input.SerialId) == 0 {
		sid := uuid.New().String()
		input.SerialId = model.PtrString(sid)

		rqid := uuid.New().String()
		rqid = strings.ReplaceAll(rqid, "-", "")
		headersMap["reqid"] = rqid
	}

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = c.PostAndDecode(*u, &response, &headersMap, &queryParams, payload)
	if err != nil {
		return &response, err
	}
	if response.OrderId == nil || len(*response.OrderId) == 0 {
		err = fmt.Errorf("Placed order not confirmed")
	}
	return &response, err
}

type PostComboOrderResponse struct {
	ComboId      *string `json:"comboId,omitempty"`
	LastSerialId *string `json:"lastSerialId,omitempty"`
}

type PostComboRequest struct {
	Orders   []model.PostStockOrderRequest `json:"newOrders,omitempty"`
	SerialId *string                       `json:"serialId,omitempty"`
}

func (c *Client) PlaceOrderV5Combo(accountID int64,
	slOrder *model.PostStockOrderRequest,
	tpOrder *model.PostStockOrderRequest) (*PostComboOrderResponse, error) {
	var (
		u, _        = url.Parse(UsTradeEndpointV + "/order/comboOrderPlace")
		response    PostComboOrderResponse
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	queryParams["secAccountId"] = strconv.FormatInt(accountID, 10)
	osid := uuid.New().String()
	pcr := PostComboRequest{
		Orders:   nil,
		SerialId: model.PtrString(osid),
	}

	if slOrder != nil {
		if slOrder.SerialId == nil || len(*slOrder.SerialId) == 0 {
			sid := uuid.New().String()
			slOrder.SerialId = model.PtrString(sid)
		}
		pcr.Orders = append(pcr.Orders, *slOrder)
	}
	if tpOrder != nil {
		if tpOrder.SerialId == nil || len(*tpOrder.SerialId) == 0 {
			sid := uuid.New().String()
			tpOrder.SerialId = model.PtrString(sid)
		}
		pcr.Orders = append(pcr.Orders, *tpOrder)
	}

	rqid := uuid.New().String()
	rqid = strings.ReplaceAll(rqid, "-", "")
	headersMap["reqid"] = rqid

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()
	payload, err := json.Marshal(pcr)
	if err != nil {
		return nil, err
	}

	err = c.PostAndDecode(*u, &response, &headersMap, &queryParams, payload)
	if err != nil {
		return &response, err
	}
	if response.ComboId == nil || len(*response.ComboId) == 0 {
		err = fmt.Errorf("ComboId should not be empty")
	}
	return &response, err
}

package webull

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	model "quantfu.com/webull/openapi"
)

// CancelAllPaperOrders is a wrapper for cancelling a number of WORKING orders.
// Note: no pagination so no guarantee all orders will cancel
func (c *Client) CancelAllPaperOrders(accountID int64) ([]int64, error) {
	if paperOrders, err := c.GetPaperOrders(accountID, model.WORKING, time.Unix(0, 0), 200); err != nil {
		return nil, err
	} else if paperOrders == nil {
		return nil, fmt.Errorf("no orders returned")
	} else {
		cancelledOrders := make([]int64, 0)
		for _, order := range paperOrders {
			cancellation, err := c.CancelPaperOrder(accountID, *order.OrderId)
			if err != nil {
				fmt.Printf("TODO: fix marshalling error\n")
				cancelledOrders = append(cancelledOrders, *order.OrderId)
				//return cancelledOrders, err
			} else {
				cancelledOrders = append(cancelledOrders, *order.OrderId)
			}
			fmt.Printf("cancellation: %v", cancellation)
		}
		return cancelledOrders, nil
	}
}

// PlacePaperOrder places paper trade
func (c *Client) PlacePaperOrder(accountID int64, input model.PostStockOrderRequest) (*model.PostOrderResponse, error) {
	var (
		u, _       = url.Parse(PaperTradeEndpointV + "/paper/1/acc/" + strconv.FormatInt(accountID, 10) + "/orderop/place/" + strconv.FormatInt(*input.TickerId, 10))
		headersMap = make(map[string]string)
		response   model.PostOrderResponse
	)

	if input.SerialId == nil || len(*input.SerialId) == 0 {
		input.SerialId = &c.UUID
	}

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
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

// CancelPaperOrder cancels paper trade
func (c *Client) CancelPaperOrder(accountID int64, oid int64) (*interface{}, error) {
	var (
		u, _       = url.Parse(PaperTradeEndpointV + "/paper/1/acc/" + strconv.FormatInt(accountID, 10) + "/orderop/cancel/" + strconv.FormatInt(oid, 10))
		headersMap = make(map[string]string)
	)
	var response interface{}

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken

	err := c.PostAndDecode(*u, &response, &headersMap, nil, nil)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// ModifyPaperOrder modifies paper trade
func (c *Client) ModifyPaperOrder(accountID int64, orderID string, input model.PostStockOrderRequest) (*interface{}, error) {
	var (
		u, _       = url.Parse(PaperTradeEndpointV + "/paper/1/acc/" + strconv.FormatInt(accountID, 10) + "/orderop/modify/" + orderID)
		headersMap = make(map[string]string)
	)
	var response interface{}

	if input.SerialId == nil || len(*input.SerialId) == 0 {
		input.SerialId = &c.UUID
	}

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
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

// GetPaperOrders gets user paper trades
func (c *Client) GetPaperOrders(paperAccountID int64, orderStatus model.OrderStatus, stTime time.Time, count int32) ([]*model.OrderItemV5, error) {
	var (
		u, _       = url.Parse(PaperTradeEndpointV + "/paper/1/acc/" + strconv.FormatInt(paperAccountID, 10) + "/order")
		headersMap = make(map[string]string)
		urlMap     = make(map[string]string)
		response   []model.PaperOrder
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	if stTime.Year() > 2000 {
		urlMap["startTime"] = stTime.Format("2006-01-02")
	} else {
		urlMap["startTime"] = "2000-1-1"
	}
	urlMap["dateType"] = strings.ToUpper(string(orderStatus))
	urlMap["pageSize"] = strconv.FormatInt(int64(count), 10)
	urlMap["status"] = string(orderStatus)
	err := c.GetAndDecode(*u, &response, &headersMap, &urlMap)

	rsFiltered := make([]*model.OrderItemV5, 0)
	if err != nil {
		return rsFiltered, err
	}

	if len(response) > 0 {
		// asking for all ?
		if strings.ToLower(string(orderStatus)) != strings.ToLower(string(model.ALL)) {
			// filter based on status
			for _, o := range response {
				if strings.ToLower(o.GetStatus()) == strings.ToLower(string(orderStatus)) {
					ord := c.toOrderItemV5(o)
					rsFiltered = append(rsFiltered, ord)
				}
			}
		} else {
			for _, o := range response {
				ord := c.toOrderItemV5(o)
				rsFiltered = append(rsFiltered, ord)
			}
		}
	}
	return rsFiltered, nil
}

func (c *Client) toOrderItemV5(o model.PaperOrder) *model.OrderItemV5 {
	ord := &model.OrderItemV5{
		OrderId:                   o.OrderId,
		OutsideRegularTradingHour: o.OutsideRegularTradingHour,
		Quantity:                  o.FilledQuantity,
		FilledQuantity:            o.FilledQuantity,
		FilledAmount:              nil,
		Action:                    o.Action,
		Status:                    o.Status,
		StatusName:                o.StatusStr,
		TimeInForce:               o.TimeInForce,
		OrderType:                 o.OrderType,
		CanModify:                 o.CanModify,
		CanCancel:                 o.CanCancel,
		LmtPrice:                  o.LmtPrice,
		AuxPrice:                  nil,
		Items:                     nil,
		FilledTotalAmount:         nil,
		TotalAmount:               nil,
	}

	// currently all paper is stock
	ord.ComboTickerType = model.PtrString("stock")

	ord.Items = make([]model.OrderItemV5ItemsInner, 1)
	ord.Items[0] = model.OrderItemV5ItemsInner{
		BrokerId:                  nil,
		OrderId:                   o.OrderId,
		BrokerOrderId:             nil,
		TickerType:                o.Ticker.Template,
		Ticker:                    o.Ticker,
		Action:                    o.Action,
		OrderType:                 o.OrderType,
		TotalQuantity:             o.TotalQuantity,
		TickerId:                  o.Ticker.TickerId,
		TimeInForce:               o.TimeInForce,
		FilledQuantity:            o.FilledQuantity,
		StatusName:                o.StatusStr,
		Symbol:                    o.Ticker.Symbol,
		CreateTime0:               o.CreateTime0,
		CreateTime:                o.CreateTime,
		FilledTime0:               o.FilledTime0,
		FilledTime:                o.FilledTime,
		AvgFilledPrice:            o.AvgFilledPrice,
		CanModify:                 o.CanModify,
		CanCancel:                 o.CanCancel,
		AssetType:                 o.Ticker.Template,
		RemainQuantity:            nil,
		PlaceAmount:               nil,
		FilledAmount:              nil,
		OutsideRegularTradingHour: o.OutsideRegularTradingHour,
		AmOrPm:                    nil,
	}
	return ord
}

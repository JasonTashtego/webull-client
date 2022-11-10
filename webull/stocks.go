package webull

import (
	"fmt"
	"net/url"
	"strconv"

	model "quantfu.com/webull/openapi"
)

// GetTicker gets ticker information for a provided stock symbol
func (c *Client) GetTicker(symbol string) (*model.LookupTickerResponse, error) {
	var (
		u, _        = url.Parse(StockInfoEndpoint + "/search/tickers5")
		response    model.LookupTickerResponse
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	queryParams["keys"] = symbol
	queryParams["queryNumber"] = strconv.Itoa(int(1))

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return &response, err
	}

	return &response, err
}

// GetTickerID is a helper function for getting a ticker ID from a stock symbol
func (c *Client) GetTickerID(symbol string) (int64, error) {
	res, err := c.GetTickerV5(symbol)
	if err != nil {
		return 0, err
	}
	if len(res.Data) < 1 {
		return 0, fmt.Errorf("No ticker found")
	}
	for _, symbolInfo := range res.Data {
		return *symbolInfo.TickerId, nil
	}
	return 0, nil
}

// GetRealtimeStockQuote gets real-time data for ticker `tickerID`
func (c *Client) GetRealtimeStockQuote(tickerID int64) (*model.GetStockQuoteResponse, error) {
	var (
		u, _       = url.Parse(QuotesEndpoint + "/quote/tickerRealTimes/v5/" + strconv.FormatInt(tickerID, 10))
		response   model.GetStockQuoteResponse
		headersMap = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	err := c.GetAndDecode(*u, &response, &headersMap, nil)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// GetStockFundamentals gets stock fundamentals for ticker `tickerID`
func (c *Client) GetStockFundamentals(tickerID string) (*model.GetFundamentalsResponse, error) {
	var (
		u, _       = url.Parse(QuotesEndpoint + "/securities/financial/index/" + tickerID)
		response   model.GetFundamentalsResponse
		headersMap = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	err := c.GetAndDecode(*u, &response, &headersMap, nil)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// GetActiveGainersLosers gets the day's active gainers or losers.
func (c *Client) GetActiveGainersLosers(direction, regionID, userRegionID string) (*[]model.ActiveGainersLosers, error) {
	var (
		u, _        = url.Parse(SecuritiesEndpoint + "/securities/market/v5/card/stockActivityPc." + direction + "/list")
		response    []model.ActiveGainersLosers
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	queryParams["regionId"] = regionID
	queryParams["userRegionId"] = userRegionID

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// GetStockAnalysis gets Webull stock analysis for tickerID `tickerID`
func (c *Client) GetStockAnalysis(tickerID string) (*model.GetStockAnalysisResponse, error) {
	var (
		u, _        = url.Parse(StockInfoEndpoint + "/securities/ticker/v5/analysis/" + tickerID)
		response    model.GetStockAnalysisResponse
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return &response, err
	}

	return &response, err
}

// GetTicker gets ticker information for a provided stock symbol
func (c *Client) GetTickerV5(symbol string) (*model.GetTickerV5Response, error) {
	var (
		u, _        = url.Parse(BrokerQuotesGWEndpointV + "/search/pc/tickers")
		response    model.GetTickerV5Response
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	queryParams["keyword"] = symbol
	queryParams["pageIndex"] = strconv.Itoa(1)
	queryParams["pageSize"] = strconv.Itoa(20)

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return &response, err
	}

	return &response, err
}

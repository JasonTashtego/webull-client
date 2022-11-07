package webull

import (
	"fmt"
	"net/url"
	model "quantfu.com/webull/openapi"
	"strconv"
	"time"
)

// GetAccounts gets all associated accounts
func (c *Client) GetAccounts() (*model.GetSecurityAccountsResponse, error) {
	var (
		u, _       = url.Parse(TradeEndpoint + "/account/getSecAccountList/v4")
		response   model.GetSecurityAccountsResponse
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

// GetAccounts gets all associated accounts
func (c *Client) GetAccountsV5() (*model.GetSecurityAccountsResponseV5, error) {
	var (
		u, _       = url.Parse(TradeEndpointV + "/tradetab/display")
		response   model.GetSecurityAccountsResponseV5
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

// GetAccountID gets an account ID
func (c *Client) GetAccountID() (int64, error) {
	res, err := c.GetAccounts()
	if err != nil {
		return 0, err
	}
	if res == nil {
		return 0, fmt.Errorf("No paper trade account found")
	}
	for _, acc := range res.Data {
		return int64(*acc.SecAccountId), nil
	}
	return 0, err
}

// GetAccountIDs gets all account IDs
func (c *Client) GetAccountIDs() (accountIDs []int64, err error) {
	if res, err := c.GetAccounts(); err != nil {
		return accountIDs, err
	} else if res == nil {
		return accountIDs, fmt.Errorf("No paper trade account found")
	} else {
		accountIDs = make([]int64, len(res.Data))
		for i, acc := range res.Data {
			accountIDs[i] = int64(*acc.SecAccountId)
		}
		return accountIDs, err
	}
}

// GetAccount gets account details for account `accountID`
func (c *Client) GetAccount(accountID int) (*model.GetAccountResponse, error) {
	var (
		path       = TradeEndpoint + "/v3/home/" + strconv.Itoa(accountID)
		u, _       = url.Parse(path)
		headersMap = make(map[string]string)
		response   model.GetAccountResponse
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	err := c.GetAndDecode(*u, &response, &headersMap, nil)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// GetAccountV5 gets account details for account.
func (c *Client) GetAccountV5() (*model.GetAccountsResponseV5, error) {
	var (
		path       = TradeEndpoint + "/v5/home"
		u, _       = url.Parse(path)
		headersMap = make(map[string]string)
		response   model.GetAccountsResponseV5
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	err := c.GetAndDecode(*u, &response, &headersMap, nil)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// GetAccountV5 gets account details for account.
func (c *Client) GetNetLiquidation(accountID int64, stTime time.Time) (*[]model.NetLiqidationTrendInner, error) {
	var (
		path        = UsTradeEndpointV + "/profitloss/account/listNetLiquidationTrend"
		u, _        = url.Parse(path)
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
		response    []model.NetLiqidationTrendInner
	)

	queryParams["secAccountId"] = strconv.FormatInt(accountID, 10)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	if stTime.Year() > 2000 {
		stTimeStr := stTime.Format("2006-01-02")
		queryParams["startDate"] = stTimeStr
	}

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return &response, err
	}
	return &response, err
}

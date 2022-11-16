package webull

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	model "quantfu.com/webull/openapi"
)

// GetPaperTradeAccounts gets information for all paper accounts.
func (c *Client) GetPaperTradeAccounts() (*[]model.PaperAccount, error) {
	var (
		u, _       = url.Parse(PaperTradeEndpointV + "/myaccounts/true")
		headersMap = make(map[string]string)
		response   []model.PaperAccount
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	err := c.GetAndDecode(*u, &response, &headersMap, nil)
	if err != nil {
		return &response, err
	}
	return &response, err
}

// GetPaperTradeAccountID is a a helper function for getting a single paper trading account ID
func (c *Client) GetPaperTradeAccountID() (int64, error) {
	res, err := c.GetPaperTradeAccounts()
	if err != nil {
		return 0, err
	}
	if res == nil {
		return 0, fmt.Errorf("No paper trade account found")
	}
	for _, acc := range *res {
		return int64(*acc.Id), nil
	}
	return 0, err
}

// GetPaperTradeAccountIDs is a a helper function for getting all paper trading account IDs.
func (c *Client) GetPaperTradeAccountIDs() ([]int64, error) {
	if res, err := c.GetPaperTradeAccounts(); err != nil {
		return []int64{}, err
	} else if res == nil {
		return []int64{}, fmt.Errorf("No paper trade account found")
	} else {
		accountIDs := make([]int64, len(*res))
		for i, acc := range *res {
			accountIDs[i] = int64(*acc.Id)
		}
		return accountIDs, err
	}
}

// ResetPaperAccount gets information for all paper accounts.
/*
func (c *Client) ResetPaperAccount(newBalance int32) (*model.ResetPaperAccountResponse, error) {
	var (
		headersMap = make(map[string]string)
		response   model.ResetPaperAccountResponse
	)
	accID, err := c.GetPaperTradeAccountID()
	if err != nil {
		return nil, err
	}
	var u, _       = url.Parse(PaperTradeEndpointV + "/paper/1/acc/reset/" + accID + "/" + fmt.Sprintf("%d", newBalance))
	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	err = c.GetAndDecode(*u, &response, &headersMap, nil)
	if err != nil {
		return &response, err
	}
	return &response, err
}
*/

func (c *Client) GetNetLiquidationPaper(accountID int64, stTime time.Time) (*[]model.NetLiqidationTrendInner, error) {
	var (
		path        = PaperTradeEndpointV + "/paper/1/acc/" + strconv.FormatInt(accountID, 10) + "/accountpl/summary"
		u, _        = url.Parse(path)
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
		response    model.PaperSummary
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	if stTime.Year() > 2000 {
		stTimeStr := stTime.Format("2006-01-02")
		queryParams["startDate"] = stTimeStr
	}

	rsp := make([]model.NetLiqidationTrendInner, 0)

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return &rsp, err
	}

	for _, s := range response.NetLiquidationHistories.Items {

		nl := model.NetLiqidationTrendInner{
			Currency:       nil,
			Date:           s.Date,
			NetLiquidation: s.Value,
		}
		rsp = append(rsp, nl)
	}
	return &rsp, err
}

func (c *Client) GetPaperAccountSummary(accountID int64) (*model.PaperAccountSummary, error) {

	var (
		path        = PaperTradeEndpointV + "/paper/1/acc/" + strconv.FormatInt(accountID, 10)
		u, _        = url.Parse(path)
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
		response    model.PaperAccountSummary
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID
	headersMap[HeaderKeyTradeToken] = c.TradeToken
	headersMap[HeaderKeyTradeTime] = getTimeSeconds()

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return nil, err
	}

	for i, p := range response.Positions {
		if p.GetTickerType() == "EQUITY" {
			response.Positions[i].SetAssetType("stock")
		}
	}

	return &response, err
}

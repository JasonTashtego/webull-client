package webull

import (
	"fmt"
	"net/url"
	"strconv"

	model "gitlab.com/brokerage-api/webull-openapi/openapi"
)

// GetAccounts returns all the accounts associated with a login/client.
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

// GetAccountID returns an account ID
func (c *Client) GetAccountID() (string, error) {
	res, err := c.GetAccounts()
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", fmt.Errorf("No paper trade account found")
	}
	for _, acc := range *&res.Data {
		return fmt.Sprintf("%d", acc.SecAccountId), nil
	}
	return "", err
}

// GetAccount returns an account
func (c *Client) GetAccount(accountID int) (*model.GetAccountResponse, error) {
	var (
		path       = TradeEndpoint + "/v2/home/" + strconv.Itoa(int(accountID))
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

package webull

import (
	"net/url"
	"strconv"

	model "quantfu.com/webull/openapi"
)

// GetAccountDividends gets account `accountID` total dividends.
func (c *Client) GetAccountDividends(accountID int64) (*model.GetDividendsResponse, error) {
	var (
		u, _        = url.Parse(TradeEndpoint + "/v2/account/" + strconv.FormatInt(accountID, 10) + "/dividends")
		response    model.GetDividendsResponse
		headersMap  = make(map[string]string)
		queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	queryParams["direct"] = "in"

	err := c.GetAndDecode(*u, &response, &headersMap, &queryParams)
	if err != nil {
		return &response, err
	}

	return &response, err
}

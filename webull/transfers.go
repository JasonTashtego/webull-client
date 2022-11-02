package webull

import (
	"bytes"
	"encoding/json"
	"strconv"

	// "fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	model "quantfu.com/webull/openapi"
)

// GetTransfers returns Transfers.
func (c *Client) GetTransfers(accountID int64, count uint32) (*model.Transfers, error) {
	var (
		u, _       = url.Parse(TradeEndpoint + "/asset/" + strconv.FormatInt(accountID, 10) + "/getWebullTransferList")
		response   *model.Transfers
		headersMap = make(map[string]string)
		// queryParams = make(map[string]string)
	)

	headersMap[HeaderKeyAccessToken] = c.AccessToken
	headersMap[HeaderKeyDeviceID] = c.DeviceID

	// Login request body
	lrId := "0"
	ct := float32(count)
	request := model.GetTransfersRequest{
		PageSize:     &ct,
		LastRecordId: &lrId,
	}
	requestBody, _ := json.Marshal(request)
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	req.Header.Add(HeaderKeyDeviceID, c.DeviceID)
	req.Header.Add(HeaderKeyAccessToken, c.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}
	// Send and parse request
	err = c.DoAndDecode(req, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

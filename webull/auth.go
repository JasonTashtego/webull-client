package webull

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	// "fmt"
	"io"
	"io/ioutil"

	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	model "quantfu.com/webull/openapi"
)

const (
	// DefaultDeviceID if none is supplied by user.
	DefaultDeviceID = "2292c4714f144eb08ed3edec7f7ce284"
	// PasswordSalt is used for salting your password
	PasswordSalt = "wl_app-a&b@!423^"
	// DefaultDeviceName is a device name
	DefaultDeviceName = "test"
	// DefaultTokenExpiryFormat is used to parse the custom datetime returned by Webull
	DefaultTokenExpiryFormat = "2006-01-02T15:04:05.000+0000"
)

// Credentials implements oauth2 using the webull implementation
type Credentials struct {
	Username    string
	Password    string
	DeviceID    string
	TradePIN    string
	MFA         string
	DeviceName  string
	AccountType model.AccountType
	Creds       oauth2.TokenSource
}

// Token implements TokenSource
func (c *Client) Token() (*oauth2.Token, error) {
	var (
		u, _       = url.Parse(UserBrokerEndpoint + "/passport/login/v5/account")
		response   model.PostLoginResponse
		httpClient = http.Client{Timeout: time.Second * 10}
		cliID      = c.DeviceID
		deviceName = c.DeviceName
	)
	// Client ID
	if cliID == "" {
		cliID = DefaultDeviceID
	}
	// Device Name
	if deviceName == "" {
		deviceName = DefaultDeviceName
	}
	// Login request body
	grade := int32(0)
	rgn := int32(1)
	requestBody, err := json.Marshal(model.PostLoginParametersRequest{
		Account:     &c.Username,
		AccountType: &c.AccountType,
		DeviceId:    &cliID,
		DeviceName:  &deviceName,
		Grade:       &grade,
		Pwd:         &c.HashedPassword,
		RegionId:    &rgn,
		ExtInfo:     &model.PostLoginParametersRequestExtInfo{VerificationCode: &c.MFA},
	})
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(HeaderKeyDeviceID, cliID)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}
	tok := oauth2.Token{}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode/100 != 2 {
		b := &bytes.Buffer{}
		var e model.ErrorBody
		err = json.NewDecoder(io.TeeReader(res.Body, b)).Decode(&e)
		if err != nil {
			return nil, fmt.Errorf("got response %q and could not decode error body %q", res.Status, b.String())
		}
		return nil, fmt.Errorf(*e.Msg)
	}

	if err != nil {
		return nil, fmt.Errorf("Got read error on body: %s", err.Error())
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("Got JSON unmarshal error on body: %s", err.Error())
	}
	if err != nil {
		return nil, err
	}
	tok.Expiry, err = time.Parse(DefaultTokenExpiryFormat, *response.TokenExpireTime)
	c.AccessTokenExpiration = tok.Expiry
	tok.TokenType = "Token"
	tok.AccessToken, c.AccessToken = *response.AccessToken, *response.AccessToken
	tok.RefreshToken, c.RefreshToken = *response.RefreshToken, *response.RefreshToken
	c.UUID = *response.Uuid
	return &tok, nil
}

// Login implements TokenSource
func (c *Client) Login(creds Credentials) (err error) {
	var (
		// u, _     = url.Parse(UserBrokerEndpoint + "/passport/login/v5/account")
		u, _     = url.Parse(UserBrokerEndpoint + "/login/account/v2")
		hasher   = md5.New()
		response model.PostLoginResponse
	)

	// Client ID
	if creds.DeviceID != "" {
		c.DeviceID = creds.DeviceID
	} else {
		if c.DeviceID == "" {
			c.DeviceID = DefaultDeviceID
		}
	}

	// Client Name
	if creds.DeviceName != "" {
		c.DeviceName = creds.DeviceName
	} else {
		if c.DeviceName == "" {
			c.DeviceName = DefaultDeviceName
		}
	}

	// Client Name
	if creds.Username != "" {
		c.Username = creds.Username
	} else {
		if c.Username == "" {
			return fmt.Errorf("Username required")
		}
	}

	// Client Name
	if creds.Password != "" {
		// UTF-8 encoded salted password
		hasher.Write([]byte(PasswordSalt + creds.Password))
		c.HashedPassword = hex.EncodeToString(hasher.Sum(nil))
	} else {
		if c.HashedPassword == "" {
			return fmt.Errorf("Password has not been set")
		}
	}
	c.AccountType = creds.AccountType

	// if we have meta-data passed, then
	// copy and return
	if c.haveMetaData(false) {

		// check that it's good.
		_, err = c.GetAccounts()
		if err == nil {
			return nil
		}
		c.expireTokens()
	}

	grade := int32(1)
	rgn := int32(1)

	// Login request body
	request := model.PostLoginParametersRequest{
		Account:     &c.Username,
		AccountType: &c.AccountType,
		DeviceId:    &c.DeviceID,
		DeviceName:  &c.DeviceName,
		Grade:       &grade,
		Pwd:         &c.HashedPassword,
		RegionId:    &rgn,
		ExtInfo:     model.NewPostLoginParametersRequestExtInfo(),
	}

	request.ExtInfo.VerificationCode = &creds.MFA
	requestBody, _ := json.Marshal(request)
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	req.Header.Add(HeaderKeyDeviceID, c.DeviceID)
	if err != nil {
		return errors.Wrap(err, "could not create request")
	}
	err = c.DoAndDecode(req, &response)
	if err != nil {
		return err
	}
	c.AccessToken = *response.AccessToken
	c.AccessTokenExpiration, err = time.Parse(DefaultTokenExpiryFormat, *response.TokenExpireTime)
	if err != nil {
		// default to +1 week.
		c.AccessTokenExpiration = time.Now().AddDate(0,0,7)
	}
	c.RefreshToken = *response.RefreshToken
	c.UUID = *response.Uuid

	// update acct meta w/tokens
	c.updateMetaData()

	return
}

// TradeLogin implements TokenSource
func (c *Client) TradeLogin(creds Credentials) (err error) {
	var (
		// Login URL
		u, _     = url.Parse(TradeEndpoint + "/login")
		response model.PostTradeTokenResponse
		hasher   = md5.New()
		pwd      string
	)

	// Client ID
	if creds.DeviceID != "" {
		c.DeviceID = creds.DeviceID
	} else {
		if c.DeviceID == "" {
			c.DeviceID = DefaultDeviceID
		}
	}

	// Client Name
	if creds.DeviceName != "" {
		c.DeviceName = creds.DeviceName
	} else {
		if c.DeviceName == "" {
			c.DeviceName = DefaultDeviceName
		}
	}

	// Client Name
	if creds.Username != "" {
		c.Username = creds.Username
	} else {
		if c.Username == "" {
			return fmt.Errorf("Username required")
		}
	}

	// Client Name
	if creds.TradePIN != "" {
		// UTF-8 encoded salted password
		hasher.Write([]byte(PasswordSalt + creds.TradePIN))
		pwd = hex.EncodeToString(hasher.Sum(nil))
		c.HashedPassword = pwd
	} else {
		if c.HashedPassword == "" {
			return fmt.Errorf("Password has not been set")
		}
	}

	c.AccountType = creds.AccountType

	// if we have meta-data passed, then
	// copy. test  and return
	if c.haveMetaData(true) {
		id, err := c.GetAccountID()
		if err == nil {
			_, err := c.GetOrders(strconv.FormatInt(id, 10), "all", 10)
			if err == nil {
				return nil
			}
		}
	}

	grade := int32(0)
	rgn := int32(6)

	devId := DefaultDeviceID
	devNm := DefaultDeviceName

	// Login request body
	request := model.PostLoginParametersRequest{
		Account:     &creds.Username,
		AccountType: &creds.AccountType,
		DeviceId:    &devId,
		DeviceName:  &devNm,
		Grade:       &grade,
		Pwd:         &c.HashedPassword,
		RegionId:    &rgn,
		ExtInfo:     model.NewPostLoginParametersRequestExtInfo(),
	}
	request.ExtInfo.VerificationCode = &creds.MFA
	requestBody, _ := json.Marshal(request)
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	req.Header.Add(HeaderKeyDeviceID, c.DeviceID)
	req.Header.Add(HeaderKeyAccessToken, c.AccessToken)
	if err != nil {
		return errors.Wrap(err, "could not create request")
	}
	// Send and parse request
	err = c.DoAndDecode(req, &response)
	if err != nil {
		return err
	}
	if *response.Success {
		c.TradeToken = *response.Data.TradeToken
		tokenTimeMs := *response.Data.TradeTokenExpireIn
		tmNowUtc := time.Now().UTC()
		c.TradeTokenExpiration = tmNowUtc.Add(time.Duration(tokenTimeMs) * time.Millisecond) // Assuming ms?

		// update acct meta w/tokens
		c.updateMetaData()

	}
	return nil
}

// TradeLogin implements TokenSource
func (c *Client) TradeLoginV5(creds Credentials) (err error) {
	var (
		// Login URL
		u, _     = url.Parse(TradeEndpointV + "/trade/login")
		response model.PostTradeTokenResponseData
		hasher   = md5.New()
		pwd      string
	)

	// Client ID
	if creds.DeviceID != "" {
		c.DeviceID = creds.DeviceID
	} else {
		if c.DeviceID == "" {
			c.DeviceID = DefaultDeviceID
		}
	}

	// Client Name
	if creds.DeviceName != "" {
		c.DeviceName = creds.DeviceName
	} else {
		if c.DeviceName == "" {
			c.DeviceName = DefaultDeviceName
		}
	}

	// Client Name
	if creds.Username != "" {
		c.Username = creds.Username
	} else {
		if c.Username == "" {
			return fmt.Errorf("Username required")
		}
	}

	// Client Name
	if creds.TradePIN != "" {
		// UTF-8 encoded salted password
		hasher.Write([]byte(PasswordSalt + creds.TradePIN))
		pwd = hex.EncodeToString(hasher.Sum(nil))
		c.HashedPassword = pwd
	} else {
		if c.HashedPassword == "" {
			return fmt.Errorf("Password has not been set")
		}
	}

	c.AccountType = creds.AccountType

	// if we have meta-data passed, then
	// copy, test and return
	if c.haveMetaData(true) {
		_, err = c.GetAccountsV5()
		if err == nil {
			return nil
		}
	}

	grade := int32(0)
	rgn := int32(6)

	devId := DefaultDeviceID
	devNm := DefaultDeviceName

	// Login request body
	request := model.PostLoginParametersRequest{
		Account:     &creds.Username,
		AccountType: &creds.AccountType,
		DeviceId:    &devId,
		DeviceName:  &devNm,
		Grade:       &grade,
		Pwd:         &c.HashedPassword,
		RegionId:    &rgn,
		ExtInfo:     model.NewPostLoginParametersRequestExtInfo(),
	}
	request.ExtInfo.VerificationCode = &creds.MFA
	requestBody, _ := json.Marshal(request)
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	req.Header.Add(HeaderKeyDeviceID, c.DeviceID)
	req.Header.Add(HeaderKeyAccessToken, c.AccessToken)
	if err != nil {
		return errors.Wrap(err, "could not create request")
	}
	// Send and parse request
	err = c.DoAndDecode(req, &response)
	if err != nil {
		return err
	}
	c.TradeToken = *response.TradeToken
	tokenTimeMs := *response.TradeTokenExpireIn
	tmNowUtc := time.Now().UTC()
	c.TradeTokenExpiration = tmNowUtc.Add(time.Duration(tokenTimeMs) * time.Millisecond) // Assuming ms?

	// snap tokens
	c.updateMetaData()

	return nil
}

// GetMFA requests for a 2FA code
func (c *Client) GetMFA(creds Credentials) (err error) {
	var (
		// Login URL
		u, _        = url.Parse(UserBrokerEndpoint + "/passport/verificationCode/sendCode")
		response    interface{}
		queryParams = make(map[string]string)
		headersMap  = make(map[string]string)
	)

	// Client ID
	c.DeviceID = creds.DeviceID
	if c.DeviceID == "" {
		c.DeviceID = DefaultDeviceID
	}

	headersMap["did"] = c.DeviceID

	queryParams["deviceId"] = c.DeviceID
	queryParams["accountType"] = fmt.Sprintf("%d", creds.AccountType)
	queryParams["account"] = creds.Username
	queryParams["codeType"] = "5"
	queryParams["regionCode"] = "1"

	if err != nil {
		return errors.Wrap(err, "could not create request")
	}
	// Send and parse request
	err = c.PostAndDecode(*u, response, &headersMap, &queryParams, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) IsTradeTokenValid(window int64) bool {
	if len(c.TradeToken) > 0 {
		tmNow := time.Now().UTC()
		// win early renwal, token time is at least window+ in the future
		if c.TradeTokenExpiration.Unix() > (tmNow.Unix() - window) {
			return true
		}
		return false
	}
	return false
}

func (c *Client) updateMetaData() {

	if c.MdProvider != nil {
		metaMap := c.MdProvider.GetMetaMap()
		metaMap["AccessToken"] = c.AccessToken
		metaMap["AccessTokenExpiration"] = strconv.FormatInt(c.AccessTokenExpiration.Unix(), 10)
		metaMap["RefreshToken"] = c.RefreshToken
		metaMap["Client.Uuid"] = c.UUID
		metaMap["TradeToken"] = c.TradeToken
		metaMap["TradeTokenExpiration"] = strconv.FormatInt(c.TradeTokenExpiration.Unix(), 10)
		c.MdProvider.Save(metaMap)
	}
}

func (c *Client) haveMetaData(withTrade bool) bool {

	if c.MdProvider != nil {
		var ok bool

		metaMap := c.MdProvider.GetMetaMap()

		c.AccessToken, ok = metaMap["AccessToken"]
		if !ok {
			return false
		}
		accessTknTmStr, ok := metaMap["AccessTokenExpiration"]
		if !ok {
			return false
		}
		accessTknTm, err := strconv.ParseInt(accessTknTmStr, 10, 64)
		if err != nil {
			return false
		}
		c.AccessTokenExpiration = time.Unix(accessTknTm, 0)
		c.RefreshToken, ok = metaMap["RefreshToken"]
		if !ok {
			return false
		}
		c.UUID, ok = metaMap["Client.Uuid"]
		if !ok {
			return false
		}

		if accessTknTm < time.Now().Unix() {
			c.AccessTokenExpiration = time.Now().AddDate(-1, 0, 0)
			c.AccessToken = ""
			return false
		}


		if withTrade {
			c.TradeToken, ok = metaMap["TradeToken"]
			if !ok || len(c.TradeToken) == 0 {
				return false
			}
			tradeTknTmStr, ok := metaMap["TradeTokenExpiration"]
			if !ok {
				return false
			}
			tradeTknTm, err := strconv.ParseInt(tradeTknTmStr, 10, 64)
			if err != nil {
				return false
			}
			c.TradeTokenExpiration = time.Unix(tradeTknTm, 0)

			if tradeTknTm < time.Now().Unix() {
				c.TradeTokenExpiration = time.Now().AddDate(-1, 0, 0)
				c.TradeToken = ""
				return false
			}

		}
		return true
	}

	return false
}

func (c *Client) expireTokens() {
	c.AccessTokenExpiration = time.Now().AddDate(-1, 0, 0)
	c.AccessToken = ""
	c.TradeTokenExpiration = time.Now().AddDate(-1, 0, 0)
	c.TradeToken = ""
}

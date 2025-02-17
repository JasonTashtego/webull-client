package webull

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"quantfu.com/webull/client/internal"
	"time"

	model "quantfu.com/webull/openapi"
)

// Endpoints for the Webull API
const (
	QuotesEndpoint         = "https://quoteapi.webull.com/api"
	UserEndpoint           = "https://userapi.webull.com/api"
	BrokerQuotesEndpoint   = "https://quoteapi.webullbroker.com/api"
	BrokerQuotesGWEndpoint = "https://quotes-gw.webullbroker.com/api"
	SecuritiesEndpoint     = "https://securitiesapi.webullbroker.com/api"
	//UserBrokerEndpoint     = "https://userapi.webullbroker.com/api"
	// UserBrokerEndpoint = "https://userapi.webull.com/api"
	//PaperTradeEndpoint     = "https://act.webullbroker.com/webull-paper-center/api"
	PaperTradeEndpointV     = "https://act.webullfintech.com/webull-paper-center/api"
	TradeEndpoint           = "https://tradeapi.webulltrade.com/api/trade"
	TradeEndpointV          = "https://trade.webullfintech.com/api/trading/v1/global"
	UsTradeEndpointV        = "https://ustrade.webullfinance.com/api/trading/v1/webull"
	BrokerQuotesGWEndpointV = "https://quotes-gw.webullfintech.com/api"

	UserBrokerEndpoint = "https://nauser.webullfintech.com/api/user/v1"

	StockInfoEndpoint = "https://infoapi.webull.com/api"
)

// HTTPClient is the context key to use with golang.org/x/net/context's
// WithValue function to associate an *http.Client value with a context.
var HTTPClient internal.ContextKey

// ErrAuthExpired signals the user must retrieve a new token
//var ErrAuthExpired = errors.New("Authentication token expired")

// AuthExpiredError returned when token needs to be refreshed
type AuthExpiredError struct{}

type userCallback func(context.Context, Topic, interface{}) error

func (e *AuthExpiredError) Error() string {
	return fmt.Sprint("Authentication token expired")
}

type MetaDataProvider interface {
	HasData() bool
	GetMetaMap() map[string]string
	Save(map[string]string)
}

// Client is a helpful abstraction around some common metadata required for
// API operations.
type Client struct {
	Username       string
	HashedPassword string
	AccountType    model.AccountType
	MFA            string
	UUID           string
	DeviceName     string

	AccessToken           string
	AccessTokenExpiration time.Time
	RefreshToken          string

	TradeToken           string
	TradeTokenExpiration time.Time

	DeviceID string

	httpClient         *http.Client
	WebsocketCallbacks map[string]userCallback

	sessionHeaders map[string]string

	MdProvider MetaDataProvider
}

// NewClient is a constructor for the Webull-Client client
func NewClient(creds *Credentials) (c *Client, err error) {
	c = &Client{
		httpClient: &http.Client{Timeout: time.Second * 10},
	}
	c.sessionHeaders = make(map[string]string)

	if creds != nil {
		c.DeviceID = creds.DeviceID
		c.Username = creds.Username
		c.AccountType = creds.AccountType
		if creds.DeviceID == "" {
			c.DeviceID = DefaultDeviceID
		}
		hasher := md5.New()
		hasher.Write([]byte(PasswordSalt + creds.Password))
		c.HashedPassword = hex.EncodeToString(hasher.Sum(nil))
		_, err = c.Token()
		if err != nil {
			return nil, err
		}
	}
	return
}

func NewClientWithContext(ctx context.Context, creds *Credentials) (c *Client, err error) {
	httpCln := internal.ContextClient(ctx)
	c = &Client{
		httpClient: httpCln,
	}
	c.sessionHeaders = make(map[string]string)
	if creds != nil {
		c.DeviceID = creds.DeviceID
		c.Username = creds.Username
		c.AccountType = creds.AccountType
		if creds.DeviceID == "" {
			c.DeviceID = DefaultDeviceID
		}
		hasher := md5.New()
		hasher.Write([]byte(PasswordSalt + creds.Password))
		c.HashedPassword = hex.EncodeToString(hasher.Sum(nil))
		_, err = c.Token()
		if err != nil {
			return nil, err
		}
	}
	return

}

func (c *Client) HttpClient() *http.Client {
	return c.httpClient
}

func (c *Client) AddSessionHeader(k string, v string) {
	c.sessionHeaders[k] = v
}

// RegisterCallback registers a callback, overriding an existing callback if one exists
func (c *Client) RegisterCallback(override bool, callback func(context.Context, Topic, interface{}) error, topic ...string) error {
	if c.WebsocketCallbacks == nil {
		c.WebsocketCallbacks = make(map[string]userCallback, 0)
	}
	for _, t := range topic {
		if _, ok := c.WebsocketCallbacks[t]; ok {
			if !override {
				return fmt.Errorf("callback already exists")
			}
		}
		c.WebsocketCallbacks[t] = callback
	}
	return nil
}

// DeregisterCallback de-registers (unsets) a callback for a particular topic number
func (c *Client) DeregisterCallback(topic string) error {
	if c.WebsocketCallbacks == nil {
		c.WebsocketCallbacks = make(map[string]userCallback, 0)
	}
	if _, ok := c.WebsocketCallbacks[topic]; ok {
		delete(c.WebsocketCallbacks, topic)
	} else {
		return fmt.Errorf("callback does not exist")
	}
	return nil
}

// GetAndDecode retrieves from the endpoint and unmarshals resulting json into
// the provided destination interface, which must be a pointer.
func (c *Client) GetAndDecode(URL url.URL, dest interface{}, headers *map[string]string, urlValues *map[string]string) error {
	if time.Now().After(c.AccessTokenExpiration) {
		return &AuthExpiredError{}
	}
	v := url.Values{}
	if urlValues != nil {
		for key, val := range *urlValues {
			v.Add(key, val)
		}
	}
	URL.RawQuery = v.Encode()

	if req, err := http.NewRequest(http.MethodGet, URL.String(), nil); err != nil {
		return err
	} else if req == nil {
		return fmt.Errorf("unable to create request")
	} else {
		if len(c.sessionHeaders) > 0 {
			for key, val := range c.sessionHeaders {
				req.Header.Add(key, val)
			}
		}
		if headers != nil {
			for key, val := range *headers {
				req.Header.Add(key, val)
			}
		}
		return c.DoAndDecode(req, dest)
	}
}

// PostAndDecode retrieves from the endpoint and unmarshals resulting json into
// the provided destination interface, which must be a pointer.
func (c *Client) PostAndDecode(URL url.URL, dest interface{}, headers *map[string]string, urlValues *map[string]string, payload []byte) error {
	if c.AccessToken != "" {
		if time.Now().After(c.AccessTokenExpiration) {
			return &AuthExpiredError{}
		}
	}
	v := url.Values{}
	if urlValues != nil {
		for key, val := range *urlValues {
			v.Set(key, val)
		}
	}
	URL.RawQuery = v.Encode()
	uStr := URL.String()
	if req, err := http.NewRequest(http.MethodPost, uStr, bytes.NewReader(payload)); err != nil {
		return err
	} else if req == nil {
		return fmt.Errorf("unable to create request")
	} else {
		if len(c.sessionHeaders) > 0 {
			for key, val := range c.sessionHeaders {
				req.Header.Add(key, val)
			}
		}
		if headers != nil {
			for key, val := range *headers {
				req.Header.Add(key, val)
			}
		}
		return c.DoAndDecode(req, dest)
	}
}

func parseAnything(data []byte) (output interface{}, err error) {
	if err = json.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("Unable to marshal body as interface")
	}
	return output, nil
}

// ConnectWebsockets connects to a streaming API by Webull
// NOTE: client still unstable
func (c *Client) ConnectWebsockets(ctx context.Context, messageTypes []string, tickerIDs []string) (err error) {
	err = c.ConnectStreamingQuotes(ctx, c.Username, c.HashedPassword, c.DeviceID, c.AccessToken, messageTypes, tickerIDs)
	return err
}

// DoAndDecode provides useful abstractions around common errors and decoding
// issues. Ideally unmarshals into `dest`. On error, it'll use the Webull `ErrorBody` model.
// Last fallback is a plain interface.
func (c *Client) DoAndDecode(req *http.Request, dest interface{}) (err error) {
	var anyBody interface{}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Got read error on body: %s", err.Error())
	}

	var e model.ErrorBody
	if res.StatusCode/100 != 2 {
		b := &bytes.Buffer{}
		err = json.Unmarshal(body, &e)
		if err != nil {
			// anything
			if anyBody, err = parseAnything(body); err != nil {
				return fmt.Errorf("Unable to unmarshal body as interface")
			}
			dest = anyBody
			return fmt.Errorf("got response %q and could not decode error body %q", res.Status, b.String())
		}
		// anything
		if anyBody, err = parseAnything(body); err != nil {
			return fmt.Errorf("Unable to unmarshal body as interface")
		}
		dest = anyBody

		msg := ""
		if e.Msg != nil {
			msg += *e.Msg
		}
		if e.Code != nil {
			msg += " - " + *e.Code
		}
		return fmt.Errorf(msg)
	}
	if err = json.Unmarshal(body, &dest); err != nil {
		// handle 200? w/error body
		err = json.Unmarshal(body, &e)
		if err != nil {
			// anything
			var err2 error
			if anyBody, err2 = parseAnything(body); err2 != nil {
				return fmt.Errorf("Unable to unmarshal body as interface")
			}
			dest = anyBody
		} else {
			return nil
		}
	}
	return err
}

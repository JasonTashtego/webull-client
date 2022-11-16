package webull

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	model "quantfu.com/webull/openapi"
)

func TestGetPaperTradeAccounts(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	c, err := NewClient(&Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)
	asrt.NotNil(c)
	res, err := c.GetPaperTradeAccounts()
	asrt.Empty(err)
	asrt.NotEmpty(res)
}

func TestGetPaperTradeAccountID(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	c, err := NewClient(&Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)
	asrt.NotNil(c)
	paperAccID, err := c.GetPaperTradeAccountID()
	asrt.Empty(err)
	asrt.NotEmpty(paperAccID)
}

/*
func TestResetPaperAccount(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	c, err := NewClient(&Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)
	asrt.NotNil(c)
	paperAccID, err := c.ResetPaperAccount(5000)
	asrt.Empty(err)
	asrt.NotEmpty(paperAccID)
}
*/

func TestGetNetLiquidationPaper(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	c, err := NewClient(nil)
	asrt.Empty(err)
	err = c.Login(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		AccountType: model.AccountType(2),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)

	paperAccID, err := c.GetPaperTradeAccountID()
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		stTime := time.Date(2022, 11, 1, 12, 0, 0, 0, time.UTC)
		lv, err := c.GetNetLiquidationPaper(paperAccID, stTime)
		asrt.Empty(err)
		asrt.NotNil(lv)
	}

}

func TestGetPaperAccountSummary(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	c, err := NewClient(nil)
	asrt.Empty(err)
	err = c.Login(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		AccountType: model.AccountType(2),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)

	paperAccID, err := c.GetPaperTradeAccountID()
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		summary, err := c.GetPaperAccountSummary(paperAccID)
		asrt.Empty(err)
		asrt.NotNil(summary)
	}
}

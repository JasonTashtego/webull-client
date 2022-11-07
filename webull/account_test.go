package webull

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	model "quantfu.com/webull/openapi"
)

func TestGetAccounts(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	c, err := NewClient(&Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt := assert.New(t)
	asrt.Empty(err)
	res, err := c.GetAccounts()
	asrt.Empty(err)
	asrt.True(*res.Success)
}

func TestGetAccount(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	c, err := NewClient(nil)
	err = c.Login(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)
	accs, err := c.GetAccounts()
	asrt.Empty(err)
	asrt.True(*accs.Success)
	if accs.Data == nil {
		t.Errorf("No accounts returned")
		t.FailNow()
	}
	if len(accs.Data) < 1 {
		t.Errorf("No accounts returned")
		t.FailNow()
	}

	acc, err := c.GetAccount(int(*accs.Data[0].SecAccountId))
	asrt.Empty(err)
	asrt.NotNil(acc)
}

func TestGetAccountV5(t *testing.T) {
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

	accs, err := c.GetAccountV5()
	asrt.Empty(err)
	asrt.NotNil(accs)

}

func TestGetNetLiquidation(t *testing.T) {
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

	if accts, err := c.GetAccountsV5(); err != nil {
		t.Log(err)
		t.Fail()
	} else {

		c.AddSessionHeader(HeaderLzone, *accts.AccountList[0].Rzone)

		stTime := time.Date(2022, 11, 1, 12, 0, 0, 0, time.UTC)

		lv, err := c.GetNetLiquidation(accts.AccountList[0].GetSecAccountId(), stTime)
		asrt.Empty(err)
		asrt.NotNil(lv)
	}

}

func TestGetAccountID(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	c, err := NewClient(nil)
	err = c.Login(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)
	res, err := c.GetAccountID()
	asrt.Empty(err)
	asrt.NotEmpty(res)
}

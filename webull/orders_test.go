package webull

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	model "quantfu.com/webull/openapi"
)

func TestGetOrders(t *testing.T) {
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
	if accountIDs, err := c.GetAccountIDs(); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		for _, accID := range accountIDs {
			orders, err := c.GetOrders(strconv.FormatInt(accID, 10), "all", 200)
			asrt.Empty(err)
			asrt.NotEmpty(orders)
		}
	}
}

func TestGetOrdersV5(t *testing.T) {
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
	if accts, err := c.GetAccountsV5(); err != nil {
		t.Log(err)
		t.Fail()
	} else {

		c.AddSessionHeader(HeaderLzone, *accts.AccountList[0].Rzone)

		st := time.Time{}
		et := time.Now()
		for _, acct := range accts.AccountList {
			if acct.SecAccountId != nil {
				orders, err := c.GetOrdersV5(int64(*acct.SecAccountId), model.FILLED, st, et, 200)
				asrt.Empty(err)
				asrt.NotEmpty(orders)
			}
		}
	}
}

func TestIsTradeable(t *testing.T) {
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
	orders, err := c.IsTradeable("913243251")
	asrt.Empty(err)
	asrt.NotEmpty(orders)
}

func TestPlaceTrade(t *testing.T) {
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

	accID, err := c.GetAccountID()
	asrt.Empty(err)
	asrt.NotEmpty(accID)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})

	res, err := c.PlaceOrder(accID, model.PostStockOrderRequest{
		Action:                    model.PtrOrderSide(model.BUY),
		ComboType:                 model.PtrString("NORMAL"),
		LmtPrice:                  model.PtrFloat32(0.05),
		OrderType:                 model.PtrOrderType(model.LMT),
		OutsideRegularTradingHour: model.PtrBool(true),
		Quantity:                  model.PtrInt32(1),
		SerialId:                  model.PtrString(c.UUID),
		TickerId:                  model.PtrInt64(913243251),
		TimeInForce:               model.PtrTif(model.DAY),
	})
	asrt.Empty(err)
	asrt.NotEmpty(res)
}

// func TestPlaceOtocoTrade(t *testing.T) {
// 	if os.Getenv("WEBULL_USERNAME") == "" {
// 		t.Skip("No username set")
// 		return
// 	}
// 	asrt := assert.New(t)
// 	c, err := NewClient(&Credentials{
// 		Username:    os.Getenv("WEBULL_USERNAME"),
// 		Password:    os.Getenv("WEBULL_PASSWORD"),
// 		AccountType: model.AccountType(2),
// 	})
// 	asrt.Empty(err)
// 	asrt.NotNil(c)

// 	accID, err := c.GetAccountID()
// 	asrt.Empty(err)
// 	asrt.NotEmpty(accID)

// 	err = c.TradeLogin(Credentials{
// 		Username:    os.Getenv("WEBULL_USERNAME"),
// 		Password:    os.Getenv("WEBULL_PASSWORD"),
// 		TradePIN:    os.Getenv("WEBULL_PIN"),
// 		AccountType: model.AccountType(2),
// 	})

// 	res, err := c.PlaceOtocoOrder(accID, model.PostStockOrderRequest{
// 		Action:                    model.BUY,
// 		ComboType:                 "NORMAL",
// 		LmtPrice:                  4.69,
// 		OrderType:                 model.MKT,
// 		OutsideRegularTradingHour: true,
// 		Quantity:                  1,
// 		SerialId:                  c.UUID,
// 		TickerId:                  913243251,
// 		TimeInForce:               model.DAY,
// 	})
// 	asrt.Empty(err)
// 	asrt.NotEmpty(res)
// }

func TestCancelTrade(t *testing.T) {
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

	accID, err := c.GetAccountID()
	asrt.Empty(err)
	asrt.NotEmpty(accID)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})

	// Place Trade
	placed, err := c.PlaceOrder(accID, model.PostStockOrderRequest{
		Action:                    model.PtrOrderSide(model.BUY),
		ComboType:                 model.PtrString("NORMAL"),
		LmtPrice:                  model.PtrFloat32(4.69),
		OrderType:                 model.PtrOrderType(model.MKT),
		OutsideRegularTradingHour: model.PtrBool(true),
		Quantity:                  model.PtrInt32(1),
		SerialId:                  model.PtrString(c.UUID),
		TickerId:                  model.PtrInt64(913243251),
		TimeInForce:               model.PtrTif(model.DAY),
	})
	asrt.Empty(err)
	asrt.NotEmpty(placed)

	// Cancel Trade
	cancelled, err := c.CancelPaperOrder(accID, fmt.Sprintf("%d", *placed.OrderId))
	asrt.Empty(err)
	asrt.NotEmpty(cancelled)
}

func TestModifyTrade(t *testing.T) {
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

	accID, err := c.GetAccountID()
	asrt.Empty(err)
	asrt.NotEmpty(accID)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})

	// Place Trade
	placed, err := c.PlacePaperOrder(accID, model.PostStockOrderRequest{
		Action:                    model.PtrOrderSide(model.BUY),
		ComboType:                 model.PtrString("NORMAL"),
		LmtPrice:                  model.PtrFloat32(4.69),
		OrderType:                 model.PtrOrderType(model.MKT),
		OutsideRegularTradingHour: model.PtrBool(true),
		Quantity:                  model.PtrInt32(1),
		SerialId:                  model.PtrString(c.UUID),
		TickerId:                  model.PtrInt64(913243251),
		TimeInForce:               model.PtrTif(model.DAY),
	})
	asrt.Empty(err)
	asrt.NotEmpty(placed)

	// Cancel Trade
	_, err = c.ModifyPaperOrder(accID, fmt.Sprintf("%d", *placed.OrderId), model.PostStockOrderRequest{
		Action:                    model.PtrOrderSide(model.BUY),
		ComboType:                 model.PtrString("NORMAL"),
		LmtPrice:                  model.PtrFloat32(4.69),
		OrderType:                 model.PtrOrderType(model.MKT),
		OutsideRegularTradingHour: model.PtrBool(true),
		Quantity:                  model.PtrInt32(1),
		SerialId:                  model.PtrString(c.UUID),
		TickerId:                  model.PtrInt64(913243251),
		TimeInForce:               model.PtrTif(model.DAY),
	})
	asrt.Empty(err)
}

func TestCheckOtocoOrder(t *testing.T) {
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

	accID, err := c.GetAccountID()
	asrt.Empty(err)
	asrt.NotEmpty(accID)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	input := model.PostOtocoOrderRequest{
		[]model.PostStockOrderRequest{
			{
				Action:                    model.PtrOrderSide(model.BUY),
				ComboType:                 model.PtrString("MASTER"),
				LmtPrice:                  model.PtrFloat32(4.69),
				OrderType:                 model.PtrOrderType(model.LMT),
				OutsideRegularTradingHour: model.PtrBool(true),
				Quantity:                  model.PtrInt32(1),
				SerialId:                  model.PtrString(c.UUID),
				TickerId:                  model.PtrInt64(913243251),
				TimeInForce:               model.PtrTif(model.DAY),
			},
			{
				OrderType:                 model.PtrOrderType(model.STP),
				TimeInForce:               model.PtrTif(model.DAY),
				Quantity:                  model.PtrInt32(1),
				OutsideRegularTradingHour: model.PtrBool(false),
				Action:                    model.PtrOrderSide(model.SELL),
				TickerId:                  model.PtrInt64(913243251),
				LmtPrice:                  model.PtrFloat32(30),
				ComboType:                 model.PtrString("STOP_LOSS"),
			},
			{
				OrderType:                 model.PtrOrderType(model.LMT),
				TimeInForce:               model.PtrTif(model.DAY),
				Quantity:                  model.PtrInt32(1),
				OutsideRegularTradingHour: model.PtrBool(false),
				Action:                    model.PtrOrderSide(model.SELL),
				TickerId:                  model.PtrInt64(913243251),
				LmtPrice:                  model.PtrFloat32(50),
				ComboType:                 model.PtrString("STOP_PROFIT"),
			},
		},
	}

	// Place Trade
	placed, err := c.CheckOtocoOrder(accID, input)
	asrt.Empty(err)
	asrt.NotEmpty(placed)

	// // Cancel Trade
	// _, err = c.ModifyPaperOrder(accID, fmt.Sprintf("%d", *placed.OrderId), model.PostStockOrderRequest{
	// 	Action:                    model.BUY,
	// 	ComboType:                 "NORMAL",
	// 	LmtPrice:                  200,
	// 	OrderType:                 model.MKT,
	// 	OutsideRegularTradingHour: false,
	// 	Quantity:                  1,
	// 	SerialId:                  c.UUID,
	// 	TickerId:                  913243251,
	// 	TimeInForce:               model.DAY,
	// })
	// asrt.Empty(err)
}

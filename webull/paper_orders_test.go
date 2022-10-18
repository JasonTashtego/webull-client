package webull

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	model "quantfu.com/webull/openapi"
)

func TestPlacePaperOrder(t *testing.T) {
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

	tickerID, err := c.GetTickerID("GE")
	asrt.Empty(err)
	asrt.NotEmpty(tickerID)

	tickerIDNumber, err := strconv.Atoi(tickerID)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})

	res, err := c.PlacePaperOrder(paperAccID, model.PostStockOrderRequest{
		Action:                    model.PtrOrderSide(model.BUY),
		ComboType:                 model.PtrString("NORMAL"),
		LmtPrice:                  model.PtrFloat32(68),
		OrderType:                 model.PtrOrderType(model.LMT),
		OutsideRegularTradingHour: model.PtrBool(false),
		Quantity:                  model.PtrInt32(1),
		SerialId:                  model.PtrString(c.UUID),
		TickerId:                  model.PtrInt32(int32(tickerIDNumber)),
		TimeInForce:               model.PtrTif(model.DAY),
	})
	asrt.Empty(err)
	asrt.NotEmpty(res)
}

func TestGetPaperOrders(t *testing.T) {
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
	paperTradeAccID, err := c.GetPaperTradeAccountID()
	asrt.Empty(err)
	asrt.NotEmpty(paperTradeAccID)
	paperTradeOrders, err := c.GetPaperOrders(paperTradeAccID, "", "ORDER", model.WORKING)
	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)
	asrt.NotNil(paperTradeOrders)
}

func TestCancelPaperOrder(t *testing.T) {
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

	tickerID, err := c.GetTickerID("SPY")
	asrt.Empty(err)
	asrt.NotEmpty(tickerID)

	tickerIDNumber, err := strconv.Atoi(tickerID)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})

	// Place Trade
	placed, err := c.PlacePaperOrder(paperAccID, model.PostStockOrderRequest{
		Action:                    model.PtrOrderSide(model.BUY),
		ComboType:                 model.PtrString("NORMAL"),
		LmtPrice:                  model.PtrFloat32(200),
		OrderType:                 model.PtrOrderType(model.MKT),
		OutsideRegularTradingHour: model.PtrBool(false),
		Quantity:                  model.PtrInt32(1),
		SerialId:                  model.PtrString(c.UUID),
		TickerId:                  model.PtrInt32(int32(tickerIDNumber)),
		TimeInForce:               model.PtrTif(model.DAY),
	})
	asrt.Empty(err)
	asrt.NotEmpty(placed)

	// Cancel Trade
	cancelled, err := c.CancelPaperOrder(paperAccID, fmt.Sprintf("%d", *placed.OrderId))
	asrt.Empty(err)
	asrt.NotEmpty(cancelled)
}

func TestModifyPaperOrder(t *testing.T) {
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

	tickerID, err := c.GetTickerID("SPY")
	asrt.Empty(err)
	asrt.NotEmpty(tickerID)

	tickerIDNumber, err := strconv.Atoi(tickerID)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})

	// Place Trade
	placed, err := c.PlacePaperOrder(paperAccID, model.PostStockOrderRequest{
		Action:                    model.PtrOrderSide(model.BUY),
		ComboType:                 model.PtrString("NORMAL"),
		LmtPrice:                  model.PtrFloat32(200),
		OrderType:                 model.PtrOrderType(model.MKT),
		OutsideRegularTradingHour: model.PtrBool(false),
		Quantity:                  model.PtrInt32(1),
		SerialId:                  model.PtrString(c.UUID),
		TickerId:                  model.PtrInt32(int32(tickerIDNumber)),
		TimeInForce:               model.PtrTif(model.DAY),
	})
	asrt.Empty(err)
	asrt.NotEmpty(placed)

	// Cancel Trade
	_, err = c.ModifyPaperOrder(paperAccID, fmt.Sprintf("%d", *placed.OrderId), model.PostStockOrderRequest{
		Action:                    model.PtrOrderSide(model.BUY),
		ComboType:                 model.PtrString("NORMAL"),
		LmtPrice:                  model.PtrFloat32(200),
		OrderType:                 model.PtrOrderType(model.MKT),
		OutsideRegularTradingHour: model.PtrBool(false),
		Quantity:                  model.PtrInt32(1),
		SerialId:                  model.PtrString(c.UUID),
		TickerId:                  model.PtrInt32(int32(tickerIDNumber)),
		TimeInForce:               model.PtrTif(model.DAY),
	})
	asrt.Empty(err)
}

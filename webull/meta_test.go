package webull

import (
	"github.com/stretchr/testify/assert"
	"os"
	model "quantfu.com/webull/openapi"
	"testing"
)

type MockMetaProvider struct {
	localMeta map[string]string
}

func (m *MockMetaProvider) HasData() bool {
	return len(m.localMeta) > 0
}

func (m *MockMetaProvider) GetMetaMap() map[string]string {
	return m.localMeta
}

func (m *MockMetaProvider) Save(v map[string]string) {
	m.localMeta = v
}

func NewMockMetaProvider() *MockMetaProvider {
	return &MockMetaProvider{localMeta: make(map[string]string)}
}

var (
	metaMock *MockMetaProvider
)

func TestMetaProvider(t *testing.T) {

	metaMock = NewMockMetaProvider()

	testStoreMeta(t)
	testLoadMeta(t)
}

func testStoreMeta(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)

	metaProv := metaMock

	c, err := NewClient(nil)
	c.MdProvider = metaProv

	asrt.Empty(err)
	err = c.Login(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceID:    os.Getenv("WEBULL_DEVID"),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)

	err = c.TradeLogin(Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		AccountType: model.AccountType(2),
		TradePIN:    os.Getenv("WEBULL_PIN"),
		DeviceID:    os.Getenv("WEBULL_DEVID"),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)

	accs, err := c.GetAccountV5()
	asrt.Empty(err)
	asrt.NotNil(accs)

}

func testLoadMeta(t *testing.T) {
	asrt := assert.New(t)

	metaProv := metaMock

	c, err := NewClient(nil)
	c.MdProvider = metaProv

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

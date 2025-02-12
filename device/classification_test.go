package device

import (
	"testing"

	deviceMocks "github.com/danilarff86/miio-go/device/mocks"
	"github.com/danilarff86/miio-go/device/product"
	"github.com/stretchr/testify/assert"
)

// Device is set to non-provisional after classify.
func TestClassify1(t *testing.T) {
	baseDev := &deviceMocks.Device{}
	baseDev.On("Provisional").Return(true)
	baseDev.On("GetProduct").Return(product.Unknown, nil)
	baseDev.On("SetProvisional", false).Once()

	Classify(baseDev)
	baseDev.AssertExpectations(t)
}

// Non-provisional devices are not classified
func TestClassify2(t *testing.T) {
	baseDev := &deviceMocks.Device{}
	baseDev.On("Provisional").Return(false)

	dev, err := Classify(baseDev)
	assert.NoError(t, err)
	assert.Equal(t, baseDev, dev)
}

func Classify_SetUp(product product.Product) *deviceMocks.Device {
	dev := &deviceMocks.Device{}
	dev.On("Provisional").Return(true)
	dev.On("GetProduct").Return(product, nil)
	dev.On("SetProvisional", false)
	dev.On("Outbound").Return(nil)
	return dev
}

func TestClassify_PowerPlug(t *testing.T) {
	baseDev := Classify_SetUp(product.PowerPlug)
	baseDev.On("RefreshThrottle").Return(nil)

	dev, err := Classify(baseDev)

	assert.NoError(t, err)
	assert.IsType(t, &PowerPlug{}, dev)
}

func TestClassify_Yeelight(t *testing.T) {
	baseDev := Classify_SetUp(product.Yeelight)
	baseDev.On("RefreshThrottle").Return(nil)

	dev, err := Classify(baseDev)

	assert.NoError(t, err)
	assert.IsType(t, &Yeelight{}, dev)
}

func TestClassify_Unknown(t *testing.T) {
	baseDev := Classify_SetUp(product.Unknown)

	_, err := Classify(baseDev)

	assert.Error(t, err)
}

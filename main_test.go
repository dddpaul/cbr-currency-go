package cbr

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTodayCurrencyRates(t *testing.T) {
	rates := GetCurrencyRates()
	assert.Equal(t, "EUR", rates["EUR"].ISOCode)
	assert.Equal(t, "USD", rates["USD"].ISOCode)
}

func TestPastCurrencyRates(t *testing.T) {
	d, _ := time.Parse(DF, "01/12/2001")
	rates := FetchCurrencyRates(d)
	assert.Equal(t, float64(26.52), rates["EUR"].Value)
	assert.Equal(t, float64(29.9), rates["USD"].Value)
}

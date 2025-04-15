package tests

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestDecimal(t *testing.T) {
	d := decimal.NewFromFloat(1.1)
	f, exact := d.Float64()

	t.Logf("float: %f, exact: %v", f, exact)
}

func TestDecimalRound(t *testing.T) {
	value := decimal.NewFromFloat(12.3456)
	rounded := value.Round(2)

	f, _ := rounded.Float64()

	t.Log(f)
}

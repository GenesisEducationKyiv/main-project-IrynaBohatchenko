package test

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

func TestPackage_ShouldNotDependOn(t *testing.T) {
	archtest.Package(t, "github.com/btc-price/cmd/price_check/handler").
		ShouldNotDependOn("github.com/btc-price/internal/binanceclient",
			"github.com/btc-price/internal/coingeckoclient",
			"github.com/btc-price/internal/rateprovider",
			"github.com/btc-price/internal/subscription/storage")

	archtest.Package(t, "github.com/btc-price/internal/...").
		ShouldNotDependOn("github.com/btc-price/cmd/price_check/handler")
}

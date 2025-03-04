package domain

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/mrehanabbasi/go-logicboxes/core"
	"github.com/stretchr/testify/require"
)

var d = New(core.New(core.Config{
	ResellerID:   os.Getenv("RESELLER_ID"),
	APIKey:       os.Getenv("API_KEY"),
	IsProduction: false,
}, http.DefaultClient))

var (
	domainName = os.Getenv("TEST_DOMAIN_NAME")
	orderID    = os.Getenv("TEST_ORDER_ID")
	_          = os.Getenv("TEST_CNS") // cns
	authCode   = os.Getenv("TEST_AUTH_CODE")
)

func TestSuggestNames(t *testing.T) {
	res, err := d.SuggestNames(context.Background(), "domain", "", false, false)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestGetOrderID(t *testing.T) {
	res, err := d.GetOrderID(context.Background(), domainName)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestGetRegistrationOrderDetails(t *testing.T) {
	res, err := d.GetRegistrationOrderDetails(context.Background(), orderID, []string{"All"})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestModifyNameServers(t *testing.T) {
	res, err := d.ModifyNameServers(context.Background(), orderID, []string{"ns1.domain.asia"})
	require.NoError(t, err)
	require.NotNil(t, res)

	res, err = d.ModifyNameServers(context.Background(), orderID, []string{"ns2.domain.asia"})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestAddChildNameServer(t *testing.T) {
	res, err := d.AddChildNameServer(context.Background(), orderID, "new."+domainName, []string{"0.0.0.0", "1.1.1.1"})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestModifyPrivacyProtectionStatus(t *testing.T) {
	res, err := d.ModifyPrivacyProtectionStatus(context.Background(), orderID, true, "some reason")
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestModifyAuthCode(t *testing.T) {
	res, err := d.ModifyAuthCode(context.Background(), orderID, authCode)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestApplyTheftProtectionLock(t *testing.T) {
	res, err := d.ApplyTheftProtectionLock(context.Background(), orderID)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestGetTheListOfLocksAppliedOnDomainName(t *testing.T) {
	res, err := d.GetTheListOfLocksAppliedOnDomainName(context.Background(), orderID)
	require.NoError(t, err)
	require.NotNil(t, res)
}

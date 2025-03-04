// Package pricing contains API methods for pricing.
package pricing

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mrehanabbasi/go-logicboxes/core"
)

type Pricing interface {
	GettingCustomerPricing(ctx context.Context, customerID string) (CustomerPrice, error)
	GettingResellerPricing(ctx context.Context, resellerID string) (ResellerPrice, error)
	GettingResellerCostPricing(ctx context.Context, resellerID string) (ResellerCostPrice, error)
	GettingPromoPrices(ctx context.Context) (PromoPrice, error)
}

func New(c core.Core) Pricing {
	return &pricing{c}
}

type pricing struct {
	core core.Core
}

func (p *pricing) GettingCustomerPricing(ctx context.Context, customerID string) (CustomerPrice, error) {
	data := make(url.Values)
	data.Add("customer-id", customerID)

	resp, err := p.core.CallAPI(ctx, http.MethodGet, "products", "customer-price", data)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result CustomerPrice
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pricing) GettingResellerPricing(ctx context.Context, resellerID string) (ResellerPrice, error) {
	data := make(url.Values)
	data.Add("reseller-id", resellerID)

	resp, err := p.core.CallAPI(ctx, http.MethodGet, "products", "reseller-price", data)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result ResellerPrice
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pricing) GettingResellerCostPricing(ctx context.Context, resellerID string) (ResellerCostPrice, error) {
	data := make(url.Values)
	data.Add("reseller-id", resellerID)

	resp, err := p.core.CallAPI(ctx, http.MethodGet, "products", "reseller-cost-price", data)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result ResellerCostPrice
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pricing) GettingPromoPrices(ctx context.Context) (PromoPrice, error) {
	data := make(url.Values)

	resp, err := p.core.CallAPI(ctx, http.MethodGet, "products", "promo-details", data)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result PromoPrice
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

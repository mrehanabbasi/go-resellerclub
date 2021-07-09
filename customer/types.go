package customer

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/xpartacvs/go-resellerclub/core"
)

type JSONBool bool
type JSONFloat float64
type JSONTime time.Time

type SignUpForm struct {
	Username              string `validate:"required,email" query:"username"`
	Password              string `validate:"required,min=9,max=16,rcpassword" query:"passwd"`
	Name                  string `validate:"required" query:"name"`
	Company               string `validate:"required" query:"company"`
	Address               string `validate:"required" query:"address-line-1"`
	AddressLine2          string `validate:"omitempty" query:"address-line-2,omitempty"`
	AddressLine3          string `validate:"omitempty" query:"address-line-3,omitempty"`
	City                  string `validate:"required" query:"city"`
	State                 string `validate:"required" query:"state"`
	OtherState            string `validate:"omitempty" query:"other-state,omitempty"`
	Country               string `validate:"required,iso3166_1_alpha2" query:"country"`
	Zipcode               string `validate:"required" query:"zipcode"`
	LanguageCode          string `validate:"required" query:"lang-pref"`
	PhoneCountryCode      string `validate:"required,len=2" query:"phone-cc"`
	Phone                 string `validate:"required,number" query:"phone"`
	AltPhoneCountryCode   string `validate:"omitempty,len=2" query:"alt-phone-cc,omitempty"`
	AltPhone              string `validate:"omitempty,number" query:"alt-phone,omitempty"`
	FaxCountryCode        string `validate:"omitempty,len=2" query:"fax-cc,omitempty"`
	Fax                   string `validate:"omitempty,number" query:"fax,omitempty"`
	MobileCountryCode     string `validate:"omitempty,len=2" query:"Mobile-cc,omitempty"`
	Mobile                string `validate:"omitempty,number" query:"Mobile,omitempty"`
	VatID                 string `validate:"omitempty" query:"vat-id,omitempty"`
	SmsConcent            bool   `validate:"omitempty" query:"sms-consent,omitempty"`
	EmailMarketingConcent bool   `validate:"omitempty" query:"email-marketing-consent,omitempty"`
	AcceptPolicy          bool   `validate:"omitempty" query:"accept-policy,omitempty"`
	CustomerId            string `validate:"-"`
}

type CustomerDetail struct {
	Id                      string    `json:"customerid"`
	Username                string    `json:"username"`
	ResellerId              string    `json:"resellerid"`
	ParentId                string    `json:"parentid"`
	Name                    string    `json:"name"`
	Company                 string    `json:"company"`
	Email                   string    `json:"useremail"`
	PhoneCountryCode        string    `json:"telnocc"`
	Phone                   string    `json:"telno"`
	MobileCountryCode       string    `json:"mobilenocc"`
	Mobile                  string    `json:"mobileno"`
	Address                 string    `json:"address1"`
	AddressLine2            string    `json:"address2"`
	AddressLine3            string    `json:"address3"`
	City                    string    `josn:"city"`
	State                   string    `josn:"state"`
	StateId                 string    `josn:"stateid"`
	CountryCode             string    `josn:"country"`
	Zipcode                 string    `josn:"zip"`
	Pin                     string    `josn:"pin"`
	TimeCreation            JSONTime  `josn:"creationdt"`
	Status                  string    `josn:"customerstatus"`
	SalesContactId          string    `json:"salescontactid"`
	LanguagePreference      string    `json:"langpref"`
	TotalReceipts           JSONFloat `json:"totalreceipts"`
	Is2FA                   JSONBool  `json:"twofactorauth_enabled"`
	Is2FASms                JSONBool  `json:"twofactorsmsauth_enabled"`
	Is2FAGoogle             JSONBool  `json:"twofactorgoogleauth_enabled"`
	IsDominicanTaxConfgired JSONBool  `json:"isDominicanTaxConfiguredByParent"`
}

type CustomerCriteria struct {
	core.Criteria
	Username       string  `validate:"omitempty" query:"username,omitempty"`
	Name           string  `validate:"omitempty" query:"name,omitempty"`
	Company        string  `validate:"omitempty" query:"company,omitempty"`
	City           string  `validate:"omitempty" query:"city,omitempty"`
	State          string  `validate:"omitempty" query:"state,omitempty"`
	ReceiptLowest  float64 `validate:"omitempty" query:"total-receipt-start,omitempty"`
	ReceiptHighest float64 `validate:"omitempty" query:"total-receipt-end,omitempty"`
}

func (c CustomerCriteria) UrlValues() (url.Values, error) {
	if err := validator.New().Struct(c); err != nil {
		return url.Values{}, err
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}

	urlValues := url.Values{}
	valueCriteria := reflect.ValueOf(c)
	typeCriteria := reflect.TypeOf(c)

	for i := 0; i < valueCriteria.NumField(); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			vField := valueCriteria.Field(idx)
			tField := typeCriteria.Field(idx)
			fieldTag := tField.Tag.Get("query")
			if len(fieldTag) > 0 {
				if strings.HasSuffix(fieldTag, "omitempty") && vField.IsZero() {
					return
				}
				queryField := strings.TrimSuffix(fieldTag, ",omitempty")

				switch vField.Kind() {
				case reflect.Float32, reflect.Float64:
					rwMutex.Lock()
					urlValues.Add(queryField, fmt.Sprintf("%.2f", vField.Float()))
					rwMutex.Unlock()
				case reflect.Uint8, reflect.Uint16:
					rwMutex.Lock()
					urlValues.Add(queryField, fmt.Sprintf("%d", vField.Uint()))
					rwMutex.Unlock()
				case reflect.String:
					rwMutex.Lock()
					urlValues.Add(queryField, vField.String())
					rwMutex.Unlock()
				case reflect.Struct:
					if vField.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
						timeField := vField.Interface().(time.Time)
						rwMutex.Lock()
						urlValues.Add(queryField, fmt.Sprintf("%d", timeField.Unix()))
						rwMutex.Unlock()
					}
				case reflect.Slice:
					for j := 0; j < vField.Len(); j++ {
						vSlice := vField.Index(j)
						if vSlice.Type().Kind() == reflect.String {
							rwMutex.Lock()
							urlValues.Add(queryField, vSlice.String())
							rwMutex.Unlock()
						}
					}
				}
			}
		}(i)
	}

	wg.Wait()
	return urlValues, nil
}

func (r SignUpForm) UrlValues() (url.Values, error) {
	valider := validator.New()
	if err := valider.RegisterValidation("rcpassword", validatePassword); err != nil {
		return url.Values{}, err
	}
	if err := valider.Struct(r); err != nil {
		return url.Values{}, err
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}

	urlValues := url.Values{}
	valueForm := reflect.ValueOf(r)
	typeForm := reflect.TypeOf(r)

	for i := 0; i < valueForm.NumField(); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			vField := valueForm.Field(idx)
			tField := typeForm.Field(idx)
			fieldTag := tField.Tag.Get("query")
			if len(fieldTag) > 0 {
				if strings.HasSuffix(fieldTag, "omitempty") && vField.IsZero() {
					return
				}
				queryField := strings.TrimSuffix(fieldTag, ",omitempty")
				switch vField.Kind() {
				case reflect.String:
					rwMutex.Lock()
					urlValues.Add(queryField, vField.String())
					rwMutex.Unlock()
				case reflect.Bool:
					rwMutex.Lock()
					urlValues.Add(queryField, fmt.Sprintf("%t", vField.Bool()))
					rwMutex.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	return urlValues, nil
}

func validatePassword(fl validator.FieldLevel) bool {
	return matchPasswordWithPattern(fl.Field().String(), false)
}

func matchPasswordWithPattern(password string, withRangeOfLength bool) bool {
	if withRangeOfLength && (len(password) < 9 || len(password) > 16) {
		return false
	}
	rgxAlphaLower := regexp.MustCompile(`[a-z]`)
	rgxAlphaUpper := regexp.MustCompile(`[A-Z]`)
	rgxSymbol := regexp.MustCompile(`[\~\*\!\@\$\#\%\_\+\.\?\:\,\{\}]`)
	return rgxAlphaLower.MatchString(password) && rgxAlphaUpper.MatchString(password) && rgxSymbol.MatchString(password)
}

func (j *JSONBool) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	bValue, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*j = JSONBool(bValue)
	return nil
}

func (j JSONBool) ToBool() bool {
	return bool(j)
}

func (j *JSONFloat) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	fValue, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*j = JSONFloat(fValue)
	return nil
}

func (j JSONFloat) ToFloat64() float64 {
	return float64(j)
}

func (j *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	tValue, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*j = JSONTime(time.Unix(tValue, 0))
	return nil
}

func (j JSONTime) ToTime() time.Time {
	return time.Time(j)
}

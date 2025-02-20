package general

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/mrehanabbasi/go-logicboxes/core"
)

type (
	CountryISO string
	countryDB  map[CountryISO]string
)

// Const for ISO countries.
const (
	CountryNetherlandsAntilles              CountryISO = "AN"
	CountryCongo                            CountryISO = "CD"
	CountryCaymanIslands                    CountryISO = "KY"
	CountryEastTimor                        CountryISO = "TP"
	CountryBotswana                         CountryISO = "BW"
	CountryChina                            CountryISO = "CN"
	CountryFalklandIslands                  CountryISO = "FK"
	CountryBrazil                           CountryISO = "BR"
	CountryGuyana                           CountryISO = "GY"
	CountryMarshallIslands                  CountryISO = "MH"
	CountryUnitedArabEmirates               CountryISO = "AE"
	CountryGreenland                        CountryISO = "GL"
	CountryBosnia                           CountryISO = "BA"
	CountryHerzegovina                      CountryISO = "BA"
	CountryMorocco                          CountryISO = "MA"
	CountryUzbekistan                       CountryISO = "UZ"
	CountryComoros                          CountryISO = "KM"
	CountryTogo                             CountryISO = "TG"
	CountryUruguay                          CountryISO = "UY"
	CountryCanada                           CountryISO = "CA"
	CountryCentralAfricanRepublic           CountryISO = "CF"
	CountryBangladesh                       CountryISO = "BD"
	CountryTrinidad                         CountryISO = "TT"
	CountryTobago                           CountryISO = "TT"
	CountryTimorLeste                       CountryISO = "TL"
	CountryEquatorialGuinea                 CountryISO = "GQ"
	CountryItaly                            CountryISO = "IT"
	CountryBermuda                          CountryISO = "BM"
	CountryIreland                          CountryISO = "IE"
	CountryMacedonia                        CountryISO = "MK"
	CountryPapuaNewGuinea                   CountryISO = "PG"
	CountryAntigua                          CountryISO = "AG"
	CountryBarbuda                          CountryISO = "AG"
	CountryPoland                           CountryISO = "PL"
	CountryCzech                            CountryISO = "CZ"
	CountryTajikistan                       CountryISO = "TJ"
	CountryCapeVerde                        CountryISO = "CV"
	CountryThailand                         CountryISO = "TH"
	CountryColombia                         CountryISO = "CO"
	CountryKiribati                         CountryISO = "KI"
	CountryChile                            CountryISO = "CL"
	CountryMontenegro                       CountryISO = "ME"
	CountryLesotho                          CountryISO = "LS"
	CountryAndorra                          CountryISO = "AD"
	CountryDenmark                          CountryISO = "DK"
	CountryBarbados                         CountryISO = "BB"
	CountryBelgium                          CountryISO = "BE"
	CountryBouvetIsland                     CountryISO = "BV"
	CountryHeard                            CountryISO = "HM"
	CountryMcDonaldIslands                  CountryISO = "HM"
	CountryIsrael                           CountryISO = "IL"
	CountryTurks                            CountryISO = "TC"
	CountryCaicosIslands                    CountryISO = "TC"
	CountryRwanda                           CountryISO = "RW"
	CountryBulgaria                         CountryISO = "BG"
	CountryJersey                           CountryISO = "JE"
	CountryMozambique                       CountryISO = "MZ"
	CountryMicronesia                       CountryISO = "FM"
	CountryIceland                          CountryISO = "IS"
	CountryFrenchSouthernTerritories        CountryISO = "TF"
	CountryTokelau                          CountryISO = "TK"
	CountryMacau                            CountryISO = "MO"
	CountryRomania                          CountryISO = "RO"
	CountryAfghanistan                      CountryISO = "AF"
	CountryReunion                          CountryISO = "RE"
	CountryMyanmar                          CountryISO = "MM"
	CountryNorthernMarianaIslands           CountryISO = "MP"
	CountryMali                             CountryISO = "ML"
	CountryKyrgyzstan                       CountryISO = "KG"
	CountrySwitzerland                      CountryISO = "CH"
	CountryArmenia                          CountryISO = "AM"
	CountryTaiwan                           CountryISO = "TW"
	CountrySomalia                          CountryISO = "SO"
	CountryFaroeIslands                     CountryISO = "FO"
	CountryFinland                          CountryISO = "FI"
	CountryDjibouti                         CountryISO = "DJ"
	CountryMongolia                         CountryISO = "MN"
	CountryBelarus                          CountryISO = "BY"
	CountryTanzania                         CountryISO = "TZ"
	CountryIsleOfMan                        CountryISO = "IM"
	CountryPalau                            CountryISO = "PW"
	CountryGermany                          CountryISO = "DE"
	CountrySweden                           CountryISO = "SE"
	CountrySwaziland                        CountryISO = "SZ"
	CountryPitcairnIsland                   CountryISO = "PN"
	CountryMauritius                        CountryISO = "MU"
	CountryNorfolkIsland                    CountryISO = "NF"
	CountryMalaysia                         CountryISO = "MY"
	CountrySierraLeone                      CountryISO = "SL"
	CountryVanuatu                          CountryISO = "VU"
	CountrySouthGeorgia                     CountryISO = "GS"
	CountryTheSouthSandwichIslands          CountryISO = "GS"
	CountrySaintVincent                     CountryISO = "VC"
	CountryTheGrenadines                    CountryISO = "VC"
	CountryNigeria                          CountryISO = "NG"
	CountryUnitedKingdom                    CountryISO = "GB"
	CountryTurkey                           CountryISO = "TR"
	CountryPeru                             CountryISO = "PE"
	CountryLibya                            CountryISO = "LY"
	CountryBelize                           CountryISO = "BZ"
	CountryGuernsey                         CountryISO = "GG"
	CountryGreece                           CountryISO = "GR"
	CountryJordan                           CountryISO = "JO"
	CountryAlgeria                          CountryISO = "DZ"
	CountryNetherlands                      CountryISO = "NL"
	CountryCoteDIvoire                      CountryISO = "CI"
	CountrySvalbard                         CountryISO = "SJ"
	CountryJanMayenIslands                  CountryISO = "SJ"
	CountryAnguilla                         CountryISO = "AI"
	CountryPhilippines                      CountryISO = "PH"
	CountryZimbabwe                         CountryISO = "ZW"
	CountryMartinique                       CountryISO = "MQ"
	CountryGrenada                          CountryISO = "GD"
	CountryKorea                            CountryISO = "KR"
	CountrySuriname                         CountryISO = "SR"
	CountryLebanon                          CountryISO = "LB"
	CountryFijiIslands                      CountryISO = "FJ"
	CountryMayotte                          CountryISO = "YT"
	CountryChad                             CountryISO = "TD"
	CountrySouthAfrica                      CountryISO = "ZA"
	CountrySenegal                          CountryISO = "SN"
	CountrySaudiArabia                      CountryISO = "SA"
	CountryVenezuela                        CountryISO = "VE"
	CountryOman                             CountryISO = "OM"
	CountryBenin                            CountryISO = "BJ"
	CountryFrenchPolynesia                  CountryISO = "PF"
	CountryParaguay                         CountryISO = "PY"
	CountryKazakhstan                       CountryISO = "KZ"
	CountryHonduras                         CountryISO = "HN"
	CountrySaoTome                          CountryISO = "ST"
	CountryPrincipe                         CountryISO = "ST"
	CountrySintMaarten                      CountryISO = "SX"
	CountryNicaragua                        CountryISO = "NI"
	CountryUganda                           CountryISO = "UG"
	CountryMauritania                       CountryISO = "MR"
	CountryLatvia                           CountryISO = "LV"
	CountryBritishIndianOceanTerritory      CountryISO = "IO"
	CountryTurkmenistan                     CountryISO = "TM"
	CountryBurundi                          CountryISO = "BI"
	CountryCongoRepublic                    CountryISO = "CG"
	CountryAmericanSamoa                    CountryISO = "AS"
	CountryCameroon                         CountryISO = "CM"
	CountryUnitedStates                     CountryISO = "US"
	CountryDominica                         CountryISO = "DM"
	CountryPakistan                         CountryISO = "PK"
	CountryKosovo                           CountryISO = "XK"
	CountryBahamas                          CountryISO = "BS"
	CountryFrenchGuiana                     CountryISO = "GF"
	CountryMalawi                           CountryISO = "MW"
	CountryMexico                           CountryISO = "MX"
	CountrySingapore                        CountryISO = "SG"
	CountryIndia                            CountryISO = "IN"
	CountryAzerbaijan                       CountryISO = "AZ"
	CountryAustria                          CountryISO = "AT"
	CountryAngola                           CountryISO = "AO"
	CountryHaiti                            CountryISO = "HT"
	CountryEthiopia                         CountryISO = "ET"
	CountryBolivia                          CountryISO = "BO"
	CountrySeychelles                       CountryISO = "SC"
	CountryPanama                           CountryISO = "PA"
	CountryNamibia                          CountryISO = "NA"
	CountryAruba                            CountryISO = "AW"
	CountryVirginIslandsBirtish             CountryISO = "VI"
	CountryEstonia                          CountryISO = "EE"
	CountryIraq                             CountryISO = "IQ"
	CountryElSalvador                       CountryISO = "SV"
	CountrySlovakia                         CountryISO = "SK"
	CountryCyprus                           CountryISO = "CY"
	CountryGeorgia                          CountryISO = "GE"
	CountryDominicanRepublic                CountryISO = "DO"
	CountryMonaco                           CountryISO = "MC"
	CountryTunisia                          CountryISO = "TN"
	CountrySolomonIslands                   CountryISO = "SB"
	CountryWallis                           CountryISO = "WF"
	CountryFutunaIslands                    CountryISO = "WF"
	CountryVietnam                          CountryISO = "VN"
	CountryMalta                            CountryISO = "MT"
	CountryUnitedStatesMinorOutlyingIslands CountryISO = "UM"
	CountryGuam                             CountryISO = "GU"
	CountryGuinea                           CountryISO = "GN"
	CountryFrance                           CountryISO = "FR"
	CountryKenya                            CountryISO = "KE"
	CountrySpain                            CountryISO = "ES"
	CountryBurkinaFaso                      CountryISO = "BF"
	CountryArgentina                        CountryISO = "AR"
	CountryNiue                             CountryISO = "NU"
	CountryNorway                           CountryISO = "NO"
	CountryHongKong                         CountryISO = "HK"
	CountryGibraltar                        CountryISO = "GI"
	CountryBahrain                          CountryISO = "BH"
	CountryRussia                           CountryISO = "RU"
	CountryNepal                            CountryISO = "NP"
	CountryMadagascar                       CountryISO = "MG"
	CountrySaintHelena                      CountryISO = "SH"
	CountrySanMarino                        CountryISO = "SM"
	CountryMoldova                          CountryISO = "MD"
	CountryAXlandIslands                    CountryISO = "AX"
	CountryLiechtenstein                    CountryISO = "LI"
	CountryTuvalu                           CountryISO = "TV"
	CountryTonga                            CountryISO = "TO"
	CountryMaldives                         CountryISO = "MV"
	CountryCostaRica                        CountryISO = "CR"
	CountryGuineaBissau                     CountryISO = "GW"
	CountryAustralia                        CountryISO = "AU"
	CountryZambia                           CountryISO = "ZM"
	CountrySaintPierre                      CountryISO = "PM"
	CountryMiquelon                         CountryISO = "PM"
	CountryLaos                             CountryISO = "LA"
	CountryGambia                           CountryISO = "GM"
	CountryIndonesia                        CountryISO = "ID"
	CountryChristmasIsland                  CountryISO = "CX"
	CountrySaintKitts                       CountryISO = "KN"
	CountryNevis                            CountryISO = "KN"
	CountrySaintLucia                       CountryISO = "LC"
	CountryLithuania                        CountryISO = "LT"
	CountryBrunei                           CountryISO = "BN"
	CountryGuadeloupe                       CountryISO = "GP"
	CountryGuatemala                        CountryISO = "GT"
	CountryEcuador                          CountryISO = "EC"
	CountryYemen                            CountryISO = "YE"
	CountryGabon                            CountryISO = "GA"
	CountryLiberia                          CountryISO = "LR"
	CountryEritrea                          CountryISO = "ER"
	CountryNiger                            CountryISO = "NE"
	CountryCookIslands                      CountryISO = "CK"
	CountryAntarctica                       CountryISO = "AQ"
	CountryCambodia                         CountryISO = "KH"
	CountryNewCaledonia                     CountryISO = "NC"
	CountrySamoa                            CountryISO = "WS"
	CountryHungary                          CountryISO = "HU"
	CountryBhutan                           CountryISO = "BT"
	CountryEgypt                            CountryISO = "EG"
	CountrySudan                            CountryISO = "SD"
	CountryLuxembourg                       CountryISO = "LU"
	CountryWesternSahara                    CountryISO = "EH"
	CountryVirginIslands                    CountryISO = "VG"
	CountryKuwait                           CountryISO = "KW"
	CountryUkraine                          CountryISO = "UA"
	CountrySriLanka                         CountryISO = "LK"
	CountryNauru                            CountryISO = "NR"
	CountryMontserrat                       CountryISO = "MS"
	CountryJamaica                          CountryISO = "JM"
	CountryCocosIslands                     CountryISO = "CC"
	CountrySerbia                           CountryISO = "RS"
	CountryCuracao                          CountryISO = "CW"
	CountryPuertoRico                       CountryISO = "PR"
	CountryFranceMetropolitan               CountryISO = "FX"
	CountryVaticanCity                      CountryISO = "VA"
	CountryQatar                            CountryISO = "QA"
	CountryGhana                            CountryISO = "GH"
	CountryCroatia                          CountryISO = "HR"
	CountryAlbania                          CountryISO = "AL"
	CountryNewZealand                       CountryISO = "NZ"
	CountryPalestinian                      CountryISO = "PS"
	CountryJapan                            CountryISO = "JP"
	CountrySlovenia                         CountryISO = "SI"
	CountryPortugal                         CountryISO = "PT"
)

func fetchCountryDB(ctx context.Context, c core.Core) (countryDB, error) {
	resp, err := c.CallAPI(ctx, http.MethodGet, "country", "list", url.Values{})
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
		err := json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	keyPairs := map[string]string{}
	if err := json.Unmarshal(bytesResp, &keyPairs); err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}
	mapCountry := countryDB{}
	for key, val := range keyPairs {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			rwMutex.Lock()
			mapCountry[CountryISO(v)] = k
			rwMutex.Unlock()
		}(key, val)
	}
	wg.Wait()

	return mapCountry, nil
}

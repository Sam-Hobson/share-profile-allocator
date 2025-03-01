package grpc

import (
	"encoding/json"
	"net/url"
	pb "share-profile-allocator/internal/grpc/generated/go"
	"share-profile-allocator/internal/utils"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	nullStr = "null"
)

func displayUrl(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		utils.Log("e457deeb").Warn("Failed to parse URL", "rawURL", rawURL)
		return ""
	}

	host := parsedURL.Host
	if strings.HasPrefix(host, "www.") {
		return host
	}

	return "www." + host
}

func displayDollar(s string) string {
	if strings.HasPrefix(s, "$") {
		return s
	}
	return "$" + s
}

func displayPercentage(s string) string {
	if strings.HasSuffix(s, "%") {
		return s
	}
	return s + "%"
}

func displayDouble(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func displayInt64(f int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d\n", f)
}

var ZeroShareData WrappedShareData

type WrappedShareData struct {
	*pb.ShareData
	ParsedCompanyOfficers []CompanyOfficerData
}

func NewWrappedShareData(data *pb.ShareData) *WrappedShareData {
	return &WrappedShareData{ShareData: data}
}

func (sd *WrappedShareData) DisplaySymbol() string {
	symbol := strings.ToUpper(sd.GetSymbol())
	if symbol == "" {
		return nullStr
	}
	if str, ok := strings.CutSuffix(symbol, ".AX"); ok {
		return str
	}
	return symbol
}

func (sd *WrappedShareData) DisplayAsk() string {
	if sd.GetAsk() == 0.0 {
		return nullStr
	}
	return displayDollar(displayDouble(sd.GetAsk()))
}

func (sd *WrappedShareData) DisplayPe() string {
	if sd.GetPe() == 0.0 {
		return nullStr
	}
	return displayDouble(sd.GetPe())
}

func (sd *WrappedShareData) DisplayMarketCap() string {
	if sd.GetMarketCap() == 0 {
		return nullStr
	}
	return displayDollar(displayInt64(sd.GetMarketCap()))
}

func (sd *WrappedShareData) DisplayVolume() string {
	if sd.GetVolume() == 0 {
		return nullStr
	}
	return displayInt64(sd.GetVolume())
}

func (sd *WrappedShareData) DisplayNav() string {
	if sd.GetNav() == 0.0 {
		return nullStr
	}
	return displayDollar(displayDouble(sd.GetNav()))
}

func (sd *WrappedShareData) DisplayDividendYield() string {
	if sd.GetDividendYield() == 0.0 {
		return nullStr
	}
	return displayPercentage(displayDouble(sd.GetDividendYield()))
}

func (sd *WrappedShareData) DisplayWebsite() string {
	return displayUrl(sd.GetWebsite())
}

func (sd *WrappedShareData) DisplayNumberOfEmployees() string {
	return displayInt64(sd.GetNumFullTimeEmployees())
}

type CompanyOfficerData struct {
	MaxAge           int    `json:"maxAge"`
	Name             string `json:"name"`
	Age              int    `json:"age"`
	Title            string `json:"title"`
	YearBorn         int    `json:"yearBorn"`
	FiscalYear       int    `json:"fiscalYear"`
	TotalPay         int    `json:"totalPay"`
	ExercisedValue   int    `json:"exercisedValue"`
	UnexercisedValue int    `json:"unexercisedValue"`
}

func (d CompanyOfficerData) HasName() bool {
	return d.Name != ""
}

func (d CompanyOfficerData) HasAge() bool {
	return d.Age != 0
}

func (d CompanyOfficerData) HasTitle() bool {
	return d.Title != ""
}

func (d CompanyOfficerData) HasTotalPay() bool {
	return d.TotalPay != 0
}

func (d CompanyOfficerData) DisplayTotalPay() string {
	return displayDollar(displayInt64(int64(d.TotalPay)))
}

func (d CompanyOfficerData) HasExercisedValue() bool {
	return d.ExercisedValue != 0
}

func (d CompanyOfficerData) HasUnexercisedValue() bool {
	return d.UnexercisedValue != 0
}

func (sd *WrappedShareData) ParseCompanyOfficers() []CompanyOfficerData {
	if sd.ParsedCompanyOfficers != nil {
		return sd.ParsedCompanyOfficers
	}
	res := []CompanyOfficerData{}

	for _, officerJson := range sd.GetCompanyOfficers() {
		var officer CompanyOfficerData

		if err := json.Unmarshal([]byte(officerJson), &officer); err != nil {
			utils.Log("7ddd8ab2").Warn("Failed to parse company officer", "officerJson", officerJson)
		} else {
			res = append(res, officer)
		}
	}

	sd.ParsedCompanyOfficers = res
	return res
}

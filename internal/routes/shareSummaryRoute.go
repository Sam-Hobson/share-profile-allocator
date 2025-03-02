package routes

import (
	"net/http"
	"share-profile-allocator/internal/grpc"
	"share-profile-allocator/internal/session"
	"share-profile-allocator/internal/state"
	"share-profile-allocator/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

type keyValue[K any, V any] struct {
	Key   K
	Value V
}

type ModalData struct {
	Ticker               string
	LongName             string
	LongBusinessSummary  string
	PrettyWebsiteUrl     string
	WebsiteUrl           string
	ContactInfo          []keyValue[string, string]
	IndustryInfo         []keyValue[string, string]
	GeneralPersonnelInfo []keyValue[string, string]
	CompanyOfficers      []grpc.CompanyOfficerData
	IncomeData           []keyValue[string, string]
	MarginData           []keyValue[string, string]
	DebtData             []keyValue[string, string]
	FinanceOverviewData  []keyValue[string, string]
}

func (d *ModalData) HasLongName() bool {
	return d.LongName != ""
}

func (d *ModalData) HasLongBusinessSummary() bool {
	return d.LongBusinessSummary != ""
}

func (d *ModalData) HasWebsite() bool {
	return d.WebsiteUrl != ""
}

func (d *ModalData) HasContactInfo() bool {
	return d.ContactInfo != nil && len(d.ContactInfo) > 0
}

func (d *ModalData) HasIndustryInfo() bool {
	return d.IndustryInfo != nil && len(d.IndustryInfo) > 0
}

func (d *ModalData) HasGeneralPersonnelInfo() bool {
	return d.GeneralPersonnelInfo != nil && len(d.GeneralPersonnelInfo) > 0
}

func (d *ModalData) HasCompanyOfficers() bool {
	return d.CompanyOfficers != nil && len(d.CompanyOfficers) > 0
}

func (d *ModalData) HasIncomeData() bool {
	return d.IncomeData != nil && len(d.IncomeData) > 0
}

func (d *ModalData) HasMarginData() bool {
	return d.MarginData != nil && len(d.MarginData) > 0
}

func (d *ModalData) HasDebtData() bool {
	return d.DebtData != nil && len(d.DebtData) > 0
}

func (d *ModalData) HasFinanceOverviewData() bool {
	return d.FinanceOverviewData != nil && len(d.FinanceOverviewData) > 0
}

func (d *ModalData) HasAboutSection() bool {
	res := false
	res = res || d.HasLongBusinessSummary()
	res = res || d.HasContactInfo()
	res = res || d.HasIndustryInfo()
	return res
}

func (d *ModalData) HasPersonnelSection() bool {
	res := false
	res = res || d.HasGeneralPersonnelInfo()
	res = res || d.HasCompanyOfficers()
	return res
}

func (d *ModalData) HasFinanceSection() bool {
	res := false
	res = res || d.HasFinanceOverviewData()
	res = res || d.HasIncomeData()
	res = res || d.HasMarginData()
	res = res || d.HasDebtData()
	return res
}

func isDefault[T comparable](value T) bool {
	var zero T
	return zero == value
}

func addEntryToListIfNotDefault[K any, V comparable](l *[]keyValue[K, V], key K, value V) bool {
	if !isDefault(value) {
		*l = append(*l, keyValue[K, V]{
			Key:   key,
			Value: value,
		})
		return true
	}

	return false
}

func GetShareSummaryRoute(sessionManager *session.SessionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := sessionManager.GetSessionFromCtx(c)
		if err != nil {
			// If the session could not be retrieved, it must have expired.
			// This will suggest the user reloads their page, which will assign them a new session
			return c.Redirect(http.StatusFound, c.Request().URL.String())
		}

		ticker := strings.ToUpper(c.FormValue("ticker"))
		if ticker == "" {
			utils.Log("9e2b323a").Warn("Ticker not provided to get summary page")
			return c.String(http.StatusBadRequest, "")
		}

		shareData, err := state.GetShareDataCache().GetShareData(ticker)
		if err != nil {
			utils.Log("08f1a82b").Error("Failed to get share data from cache", "err", err.Error())
			return c.String(http.StatusBadRequest, "")
		}

		res := ModalData{
			Ticker:               ticker,
			LongName:             shareData.GetLongName(),
			LongBusinessSummary:  shareData.GetLongBusinessSummary(),
			PrettyWebsiteUrl:     shareData.DisplayWebsite(),
			WebsiteUrl:           shareData.GetWebsite(),
			ContactInfo:          []keyValue[string, string]{},
			IndustryInfo:         []keyValue[string, string]{},
			GeneralPersonnelInfo: []keyValue[string, string]{},
			CompanyOfficers:      []grpc.CompanyOfficerData{},
			IncomeData:           []keyValue[string, string]{},
			MarginData:           []keyValue[string, string]{},
			DebtData:             []keyValue[string, string]{},
			FinanceOverviewData:  []keyValue[string, string]{},
		}

		addEntryToListIfNotDefault(&res.ContactInfo, "Address 1", shareData.GetAddress1())
		addEntryToListIfNotDefault(&res.ContactInfo, "Address 2", shareData.GetAddress2())
		addEntryToListIfNotDefault(&res.ContactInfo, "City", shareData.GetCity())
		addEntryToListIfNotDefault(&res.ContactInfo, "State", shareData.GetState())
		addEntryToListIfNotDefault(&res.ContactInfo, "Zip code", shareData.GetZip())
		addEntryToListIfNotDefault(&res.ContactInfo, "Country", shareData.GetCountry())
		addEntryToListIfNotDefault(&res.ContactInfo, "Phone number", shareData.GetPhoneNumber())

		addEntryToListIfNotDefault(&res.IndustryInfo, "Industry", shareData.GetIndustry())
		addEntryToListIfNotDefault(&res.IndustryInfo, "Sector", shareData.GetSector())
		addEntryToListIfNotDefault(&res.IndustryInfo, "Exchange name", shareData.GetExchangeName())
		addEntryToListIfNotDefault(&res.IndustryInfo, "Region", shareData.GetRegion())

		if !isDefault(shareData.GetNumFullTimeEmployees()) {
			addEntryToListIfNotDefault(&res.GeneralPersonnelInfo, "Number of full time employees", shareData.DisplayNumberOfEmployees())
		}

		res.CompanyOfficers = shareData.ParseCompanyOfficers()

		if !isDefault(shareData.GetMarketCap()) {
			addEntryToListIfNotDefault(&res.FinanceOverviewData, "Market cap", shareData.DisplayMarketCap())
		}
		if !isDefault(shareData.GetEbitda()) {
			addEntryToListIfNotDefault(&res.FinanceOverviewData, "EBITDA", shareData.DisplayEbitda())
		}
		if !isDefault(shareData.GetTotalCash()) {
			addEntryToListIfNotDefault(&res.FinanceOverviewData, "Total cash", shareData.DisplayTotalCash())
		}

		if !isDefault(shareData.GetTotalRevenue()) {
			addEntryToListIfNotDefault(&res.IncomeData, "Total revenue", shareData.DisplayTotalRevenue())
		}
		if !isDefault(shareData.GetGrossProfits()) {
			addEntryToListIfNotDefault(&res.IncomeData, "Gross profit", shareData.DisplayGrossProfit())
		}
		if !isDefault(shareData.GetFreeCashFlow()) {
			addEntryToListIfNotDefault(&res.IncomeData, "Free cash flow", shareData.DisplayFreeCashFlow())
		}
		if !isDefault(shareData.GetOperatingCashFlow()) {
			addEntryToListIfNotDefault(&res.IncomeData, "Operating cash flow", shareData.DisplayOperatingCashFlow())
		}
		if !isDefault(shareData.GetEarningsGrowth()) {
			addEntryToListIfNotDefault(&res.IncomeData, "Earnings growth", shareData.DisplayEarningsGrowth())
		}
		if !isDefault(shareData.GetRevenueGrowth()) {
			addEntryToListIfNotDefault(&res.IncomeData, "Revenue growth", shareData.DisplayRevenueGrowth())
		}
		if !isDefault(shareData.GetRevenuePerShare()) {
			addEntryToListIfNotDefault(&res.IncomeData, "Revenue per share", shareData.DisplayRevenuePerShare())
		}

		if !isDefault(shareData.GetProfitMargins()) {
			addEntryToListIfNotDefault(&res.MarginData, "Profit margin", shareData.DisplayProfitMargin())
		}
		if !isDefault(shareData.GetOperatingMargin()) {
			addEntryToListIfNotDefault(&res.MarginData, "Operating margin", shareData.DisplayOperatingMargin())
		}
		if !isDefault(shareData.GetGrossMargin()) {
			addEntryToListIfNotDefault(&res.MarginData, "Gross margin", shareData.DisplayGrossMargin())
		}
		if !isDefault(shareData.GetEbitdaMargin()) {
			addEntryToListIfNotDefault(&res.MarginData, "EBITDA margin", shareData.DisplayEbitdaMargin())
		}

		if !isDefault(shareData.GetTotalDebt()) {
			addEntryToListIfNotDefault(&res.DebtData, "Total debt", shareData.DisplayTotalDebt())
		}
		if !isDefault(shareData.GetDebtToEquity()) {
			addEntryToListIfNotDefault(&res.DebtData, "Debt to equity (D/E)", shareData.DisplayDebtToEquity())
		}

		utils.Log("14768645").Info("Rendering summary popup", "res", &res)

		return c.Render(http.StatusOK, "shareSummaryPopup", &res)
	}
}

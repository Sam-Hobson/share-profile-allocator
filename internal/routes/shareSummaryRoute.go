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

func (d *ModalData) HasAboutInfo() bool {
	res := false
	res = res || d.HasLongBusinessSummary()
	res = res || d.HasContactInfo()
	res = res || d.HasIndustryInfo()
	return res
}

func (d *ModalData) HasPersonnelInfo() bool {
	res := false
	res = res || d.HasGeneralPersonnelInfo()
	res = res || d.HasCompanyOfficers()
	return res
}

func (d *ModalData) AddContactInfoEntry(key, value string) {
	if value != "" {
		d.ContactInfo = append(d.ContactInfo, keyValue[string, string]{
			Key:   key,
			Value: value,
		})
	}
}

func (d *ModalData) AddIndustryInfoEntry(key, value string) {
	if value != "" {
		d.IndustryInfo = append(d.IndustryInfo, keyValue[string, string]{
			Key:   key,
			Value: value,
		})
	}
}

func (d *ModalData) AddGeneralPersonalInfo(key, value string) {
	if value != "" {
		d.GeneralPersonnelInfo = append(d.GeneralPersonnelInfo, keyValue[string, string]{
			Key:   key,
			Value: value,
		})
	}
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
		}

		res.AddContactInfoEntry("Address 1", shareData.GetAddress1())
		res.AddContactInfoEntry("Address 2", shareData.GetAddress2())
		res.AddContactInfoEntry("City", shareData.GetCity())
		res.AddContactInfoEntry("State", shareData.GetState())
		res.AddContactInfoEntry("Zip code", shareData.GetZip())
		res.AddContactInfoEntry("Country", shareData.GetCountry())
		res.AddContactInfoEntry("Phone number", shareData.GetPhoneNumber())

		res.AddIndustryInfoEntry("Industry", shareData.GetIndustry())
		res.AddIndustryInfoEntry("Sector", shareData.GetSector())
		res.AddIndustryInfoEntry("Exchange name", shareData.GetExchangeName())
		res.AddIndustryInfoEntry("Region", shareData.GetRegion())

		if shareData.GetNumFullTimeEmployees() != 0 {
			res.AddGeneralPersonalInfo("Number of full time employees", shareData.DisplayNumberOfEmployees())
		}

		res.CompanyOfficers = shareData.ParseCompanyOfficers()

		utils.Log("14768645").Info("Rendering summary popup", "res", &res)

		return c.Render(http.StatusOK, "shareSummaryPopup", &res)
	}
}

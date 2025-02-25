package routes

import (
	"net/http"
	"share-profile-allocator/internal/session"
	"share-profile-allocator/internal/static"
	"share-profile-allocator/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetSearchTickerRoute(sessionManager *session.SessionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		tickers, err := static.GetTickers()
		if err != nil {
			utils.Log("da0f457d").Warn("Could not get tickers, so cannot provide ticker dropdown")
			return c.NoContent(http.StatusInternalServerError)
		}

		input := strings.ToUpper(strings.TrimSpace(c.FormValue("ticker")))
		if input == "" {
			utils.Log("d0d09fad").Warn("Ticker not provided for search for ticker dropdown")
			return c.NoContent(http.StatusOK)
		}

		res := []struct {
			Name        string
			Description string
		}{}

		for ticker, desc := range tickers {
			// Use tickers name as search metric, but change to using description if input is long enough
			if strings.HasPrefix(ticker, input) || (len(input) > 4 && strings.Contains(strings.ToUpper(desc), input)) {
				res = append(res, struct {
					Name        string
					Description string
				}{
					Name:        ticker,
					Description: desc,
				})
			}
		}

		return c.Render(http.StatusOK, "autoCompleteDropdown", res)
	}
}

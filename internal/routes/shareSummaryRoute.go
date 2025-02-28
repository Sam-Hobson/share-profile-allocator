package routes

import (
	"net/http"
	"share-profile-allocator/internal/session"
	"share-profile-allocator/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

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

		return c.Render(http.StatusOK, "shareSummaryPopup", ticker)
	}
}

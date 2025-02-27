package routes

import (
	"net/http"
	"share-profile-allocator/internal/session"
	"share-profile-allocator/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetDeleteShareRoute(sessionManager *session.SessionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionData, err := sessionManager.GetSessionFromCtx(c)
		if err != nil {
			// If the session could not be retrieved, it must have expired.
			// This will suggest the user reloads their page, which will assign them a new session
			utils.Log("8ffb1306").Info("Error retrieving session, requesting client to reload")
			return c.Redirect(http.StatusFound, c.Request().URL.String())
		}

		ticker := strings.ToUpper(c.FormValue("ticker"))
		if ticker == "" {
			utils.Log("27327140").Warn("Ticker not provided for deletion")
			return c.String(http.StatusBadRequest, "")
		}

		sessionData.TrackedShares.DeleteFunc(func(s string) bool { return s == ticker })

		return c.Render(http.StatusOK, "deleteShareTableRow", ticker)
	}
}

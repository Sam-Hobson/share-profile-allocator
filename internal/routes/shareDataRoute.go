package routes

import (
	"net/http"
	"share-profile-allocator/internal/grpc"
	"share-profile-allocator/internal/session"
	"share-profile-allocator/internal/state"
	"share-profile-allocator/internal/utils"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const RequestTickerDataTimeout = 2 * time.Second

func GetShareDataRoute(sessionManager *session.SessionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionData, err := sessionManager.GetSessionFromCtx(c)
		if err != nil {
			// If the session could not be retrieved, it must have expired.
			// This will suggest the user reloads their page, which will assign them a new session
			utils.Log("4e000201").Info("Error retrieving session, requesting client to reload")
			return c.Redirect(http.StatusFound, c.Request().URL.String())
		}

		ticker := strings.ToUpper(c.FormValue("ticker"))
		if ticker == "" {
			utils.Log("66979902").Warn("Ticker not provided")
			return c.String(http.StatusBadRequest, "")
		}
		// Check if the ticker is already being displayed to the user
		if sessionData.TrackedShares.IndexFunc(func(s string) bool { return ticker == s }) != -1 {
			utils.Log("ccb62758").Warn("Rejecting get share data, as ticker is already present in the session", "ticker", ticker, "sessionData", sessionData)
			return c.String(http.StatusBadRequest, "Share profile already present")
		}

		data, err := state.GetShareDataCache().GetShareData(ticker)
		if err != nil {
			utils.Log("87598df9").Error("Failed to get share data from cache", "err", err.Error())
			return c.String(http.StatusBadRequest, "")
		}

		sessionData.TrackedShares.Append(ticker)

		return c.Render(http.StatusOK, "shareTableRow", struct {
			Index int
			Data  *grpc.WrappedShareData
		}{
			Index: sessionData.TrackedShares.Len(),
			Data:  data,
		})
	}
}

package routes

import (
	"context"
	"net/http"
	"share-profile-allocator/internal/grpc"
	"share-profile-allocator/internal/session"
	"share-profile-allocator/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

const RequestTickerDataTimeout = 2 * time.Second

func GetShareDataRoute(sessionManager *session.SessionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := sessionManager.GetSession(c)
		if err != nil {
			// If the session could not be retrieved, it must have expired.
			// This will suggest the user reloads their page, which will assign them a new session
			return c.Redirect(http.StatusFound, c.Request().URL.String())
		}

		ticker := c.FormValue("ticker")
		if ticker == "" {
			utils.Log("66979902").Warn("Ticker not provided")
			return c.String(http.StatusBadRequest, "")
		}

		ctx, cancel := context.WithTimeout(context.Background(), RequestTickerDataTimeout)
		defer cancel()

		data, err := grpc.RequestDataForTicker(ctx, ticker)
		if err != nil {
			return c.String(http.StatusBadRequest, "")
		}

		return c.Render(http.StatusOK, "shareTableRow", data)
		// return c.String(http.StatusOK, fmt.Sprintf("%+v", data))
	}
}

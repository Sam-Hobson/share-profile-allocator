package routes

import (
	"net/http"
	"share-profile-allocator/internal/session"

	"github.com/labstack/echo/v4"
)

func GetRootRoute(sessionManager *session.SessionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionData, err := sessionManager.GetSessionFromCtx(c)
		if err != nil {
			// If the session could not be retrieved, it must have expired.
			// This will suggest the user reloads their page, which will assign them a new session
			return c.Redirect(http.StatusFound, c.Request().URL.String())
		}

		return c.Render(http.StatusOK, "root", sessionData)
	}
}

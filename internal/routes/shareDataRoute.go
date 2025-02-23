package routes

import (
	"context"
	"net/http"
	"share-profile-allocator/internal/grpc"
	"share-profile-allocator/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

const RequestTickerDataTimeout = 2 * time.Second

func GetShareDataRoute() echo.HandlerFunc {
	return func(c echo.Context) error {
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

		data.GetSymbol()

		return c.Render(http.StatusOK, "shareTableRow", data)
		// return c.String(http.StatusOK, fmt.Sprintf("%+v", data))
	}
}

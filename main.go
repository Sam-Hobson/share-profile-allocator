// main.go
package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"share-profile-allocator/internal/grpc"
	"share-profile-allocator/internal/routes"
	"share-profile-allocator/internal/utils"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initialiseLogger() {
	utils.SetSessionId(time.Now().UnixMicro())
	// handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))

	utils.Log("476e9ba3").Info("The logger is now configured")
}

func startFinanceServer() error {
	// Command to check Python version
	checkCmd := exec.Command("sh", "-c", "python3 -c 'import sys; print(sys.version_info.major)' 2>/dev/null || python -c 'import sys; print(sys.version_info.major)'")
	var out bytes.Buffer
	checkCmd.Stdout = &out

	if err := checkCmd.Run(); err != nil {
		utils.Log("5c3d3fc7").Error("No suitable Python interpreter found")
		return err
	}

	pythonCmd := "python3"
	if strings.TrimSpace(out.String()) != "3" {
		utils.Log("138e8ba7").Warn("Could not find `python3` in $PATH, defaulting to `python`")
		pythonCmd = "python"
	}

	cmd := exec.Command(pythonCmd, "./internal/python/yfinance_ambassador.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		utils.Log("72823c58").Error("Error occuered while executing the python finance data server", "Error", err.Error())
		return err
	}

	utils.Log("9aa40b72").Info("Started python finance server")

	return nil
}

func main() {
	initialiseLogger()

	if err := startFinanceServer(); err != nil {
		panic(err)
	}

	utils.Log("cbddd68d").Info("Starting Echo web server")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = utils.Template

	// Expose public files from server
	e.Static("/public/js", "js")
	e.Static("/public/styles", "css")

	// Setup routes
	e.GET("/", routes.GetRootRoute())
	e.POST("/sharedata", routes.GetShareDataRoute())
	// e.GET("/contacts", routes.GetContactsRoute())
	// e.POST("/contacts", routes.PostContactsRoute())
	// e.POST("/delete", routes.PostDeleteReoute())

	data, _ := grpc.RequestDataForTicker(context.Background(), "DHHF")
	fmt.Printf("\n\n%+v\n\n", data)

	e.Logger.Fatal(e.Start(":8080"))

}

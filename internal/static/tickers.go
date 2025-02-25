package static

import (
	"encoding/json"
	"os"
	"share-profile-allocator/internal/utils"
	"sync"
)

const (
	staticTickersFileName = "internal/static/tickers.json"
)

var data map[string]string
var fetchError error
var getData sync.Once

func GetTickers() (map[string]string, error) {
	getData.Do(func() {
		utils.Log("ea5907c1").Info("Reading static tickers file", "staticTickersFileName", staticTickersFileName)

		file, err := os.ReadFile(staticTickersFileName)
		if err != nil {
			fetchError = err
			utils.Log("c372ca75").Error("Failed to read static tickers file", "staticTickersFileName", staticTickersFileName, "Error", fetchError.Error())
			return
		}

		err = json.Unmarshal(file, &data)
		if err != nil {
			fetchError = err
			utils.Log("c372ca75").Error("Failed to unmarshal static tickers file", "staticTickersFileName", staticTickersFileName, "Error", fetchError.Error())
		}
	})

	return data, fetchError
}

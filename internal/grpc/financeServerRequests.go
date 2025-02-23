package grpc

import (
	"context"
	pb "share-profile-allocator/internal/grpc/generated/go"
	"share-profile-allocator/internal/utils"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var zeroShareData WrappedShareData

type WrappedShareData struct {
	*pb.ShareData
}

func NewWrappedShareData(data *pb.ShareData) *WrappedShareData {
	return &WrappedShareData{ShareData: data}
}

func (sd *WrappedShareData) GetPrettySymbol() string {
	symbol := sd.GetSymbol()
	prettySymbol, found := strings.CutSuffix(symbol, ".AX")
	if found {
		return prettySymbol
	}
	return symbol
}

func RequestDataForTicker(ctx context.Context, ticker string) (*WrappedShareData, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.Log("f9402212").Error("did not connect to finance data server", "Error", err.Error())
		return &zeroShareData, err
	}
	defer conn.Close()

	client := pb.NewShareAPIClient(conn)

	r, err := client.GetDataForTicker(ctx, &pb.Ticker{Name: ticker})
	if err != nil {
		utils.Log("e199af4d").Error("Could not retrieve data for ticker", "Ticker", ticker, "Error", err.Error())
		return &zeroShareData, err
	}

	utils.Log("76e5d6f4").Info("Successfully retrieved data for ticker", "Ticker", ticker, "Data", r)

	return NewWrappedShareData(r), nil
}

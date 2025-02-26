package grpc

import (
	"context"
	pb "share-profile-allocator/internal/grpc/generated/go"
	"share-profile-allocator/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RequestDataForTicker(ctx context.Context, ticker string) (*WrappedShareData, error) {
	utils.Log("76e5d6f4").Info("Retrieving data for ticker", "Ticker", ticker)

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.Log("f9402212").Error("did not connect to finance data server", "Error", err.Error())
		return &ZeroShareData, err
	}
	defer conn.Close()

	client := pb.NewShareAPIClient(conn)

	r, err := client.GetDataForTicker(ctx, &pb.Ticker{Name: ticker})
	if err != nil {
		utils.Log("e199af4d").Error("Could not retrieve data for ticker", "Ticker", ticker, "Error", err.Error())
		return &ZeroShareData, err
	}

	return NewWrappedShareData(r), nil
}

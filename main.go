// main.go
package main

import (
    "context"
    "log"
    "time"

    pb "share-profile-allocator/internal/grpc/go"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    c := pb.NewShareAPIClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

	ticker := "VDHG"

	r, err := c.GetDataForTicker(ctx, &pb.Ticker{Name: ticker})
    if err != nil {
        log.Fatalf("Could not get info for ticker %s: %v", ticker, err)
    }

    log.Printf("%+v", r)
}

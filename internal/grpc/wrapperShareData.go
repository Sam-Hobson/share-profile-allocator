package grpc

import (
	pb "share-profile-allocator/internal/grpc/generated/go"
	"strconv"
	"strings"
)

const (
	nullStr = "null"
)

func genericDisplayDouble(f float64) string {
	if f == 0.0 {
		return nullStr
	}
	return strconv.FormatFloat(f, 'f', 3, 64)
}

func genericDisplayInt64(f int64) string {
	if f == 0 {
		return nullStr
	}
	return strconv.FormatInt(f, 10)
}

var ZeroShareData WrappedShareData

type WrappedShareData struct {
	*pb.ShareData
}

func NewWrappedShareData(data *pb.ShareData) *WrappedShareData {
	return &WrappedShareData{ShareData: data}
}

func (sd *WrappedShareData) DisplaySymbol() string {
	symbol := strings.ToUpper(sd.GetSymbol())
	if symbol == "" {
		return nullStr
	}
	if str, ok := strings.CutSuffix(symbol, ".AX"); ok {
		return str
	}
	return symbol
}

func (sd *WrappedShareData) DisplayAsk() string {
	return genericDisplayDouble(sd.GetAsk())
}

func (sd *WrappedShareData) DisplayPe() string {
	return genericDisplayDouble(sd.GetPe())
}

func (sd *WrappedShareData) DisplayMarketCap() string {
	return genericDisplayInt64(sd.GetMarketCap())
}

func (sd *WrappedShareData) DisplayVolume() string {
	return genericDisplayInt64(sd.GetVolume())
}

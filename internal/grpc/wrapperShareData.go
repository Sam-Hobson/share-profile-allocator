package grpc

import (
	pb "share-profile-allocator/internal/grpc/generated/go"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	nullStr = "null"
)

func displayDollar(s string) string {
	return "$" + s
}

func displayPercentage(s string) string {
	return s + "%"
}

func displayDouble(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func displayInt64(f int64) string {
	p := message.NewPrinter(language.English)
    return p.Sprintf("%d\n", f)
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
	if sd.GetAsk() == 0.0 {
		return nullStr
	}
	return displayDollar(displayDouble(sd.GetAsk()))
}

func (sd *WrappedShareData) DisplayPe() string {
	if sd.GetPe() == 0.0 {
		return nullStr
	}
	return displayDouble(sd.GetPe())
}

func (sd *WrappedShareData) DisplayMarketCap() string {
	if sd.GetMarketCap() == 0 {
		return nullStr
	}
	return displayDollar(displayInt64(sd.GetMarketCap()))
}

func (sd *WrappedShareData) DisplayVolume() string {
	if sd.GetVolume() == 0 {
		return nullStr
	}
	return displayInt64(sd.GetVolume())
}

func (sd *WrappedShareData) DisplayNav() string {
	if sd.GetNav() == 0.0 {
		return nullStr
	}
	return displayDollar(displayDouble(sd.GetNav()))
}

func (sd *WrappedShareData) DisplayDividendYield() string {
	if sd.GetDividendYield() == 0.0 {
		return nullStr
	}
	return displayPercentage(displayDouble(sd.GetDividendYield()))
}

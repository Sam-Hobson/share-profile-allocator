package session

import (
	"share-profile-allocator/internal/grpc"
	"share-profile-allocator/internal/state"
	"share-profile-allocator/internal/utils"
)

type SessionState struct {
	TrackedShares *utils.ConcurrentSlice[string]
}

func NewSessionState() *SessionState {
	return &SessionState{
		TrackedShares: utils.NewConcurrentSlice[string](),
	}
}

func (ss *SessionState) GetTrackedSharesAllData() []*grpc.WrappedShareData {
	res := []*grpc.WrappedShareData{}

	tickers := ss.TrackedShares.Items()
	ch := state.GetShareDataCache().BatchGetShareData(tickers...)

	for shareData := range ch {
		res=append(res, shareData)
	}

	return res
}

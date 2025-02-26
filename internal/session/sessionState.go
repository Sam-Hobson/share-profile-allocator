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

func (ss *SessionState) GetTrackedRowsData() []*grpc.WrappedShareData {
	res := []*grpc.WrappedShareData{}

	tickers := ss.TrackedShares.Items()
	data := state.GetShareDataCache().BatchGetShareData(tickers...)

	for _, d := range data {
		if d.Err == nil {
			res = append(res, d.Data)
		}
	}

	return res
}

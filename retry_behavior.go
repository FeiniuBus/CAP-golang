package cap

import (
	"math"
	"math/rand"
)

type RetryInThunk func (retries int32) int32

var (
	DefaultRetryInThunk = func (retries int32) int32 {
		val := math.Pow(float64(retries - 1), float64(4)) + float64(15) + float64(rand.Int31n(30)) * float64(retries)
		if (val - float64(int64(val))) > 0.5 {
			return int32(val + 1)
		} else {
			return int32(val)
		}
	}
	DefaultRetryCount = int32(3)	
	DefaultRetry = NewRetryBehavior(true, DefaultRetryCount, DefaultRetryInThunk)
	NoRetry = NewRetryBehavior(false, DefaultRetryCount, DefaultRetryInThunk)
)

type RetryBehavior struct {
	RetryInThunk	RetryInThunk
	Retry			bool
	RetryCount		int32
}

func NewRetryBehavior(retry bool, count int32, thunk RetryInThunk) *RetryBehavior {
	return &RetryBehavior {
		RetryCount: count,
		RetryInThunk: thunk,
		Retry: retry,
	}
}

func (this *RetryBehavior) RetryIn(reties int32) int32 {
	return this.RetryInThunk(reties)
}
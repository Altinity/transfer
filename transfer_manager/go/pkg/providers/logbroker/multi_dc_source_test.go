package logbroker

import (
	"testing"

	"github.com/doublecloud/tross/kikimr/public/sdk/go/persqueue"
	"github.com/doublecloud/tross/library/go/core/xerrors"
	"github.com/stretchr/testify/require"
)

func mockTemporaryError(code int) *persqueue.Error {
	return &persqueue.Error{
		Code:        code,
		Description: "mock",
	}
}

func TestPersqueueTemporaryError(t *testing.T) {
	tempErrorCode, nonTempErrorCode := 0, 17
	var nonTempError = mockTemporaryError(nonTempErrorCode)
	var tempError = mockTemporaryError(tempErrorCode)
	require.False(t, isPersqueueTemporaryError(xerrors.New("irrelevant")), "irrelevant error")
	require.False(t, isPersqueueTemporaryError(nonTempError), "non-temporary")
	require.False(t, isPersqueueTemporaryError(xerrors.Errorf("oh: %w", nonTempError)), "wrapped non-temporary")
	require.True(t, isPersqueueTemporaryError(tempError), "temporary")
	require.True(t, isPersqueueTemporaryError(xerrors.Errorf("oh: %w", tempError)), "wrapped temporary")
}

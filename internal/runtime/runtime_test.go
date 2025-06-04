package runtime

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServerControl is a mock type for the serverControl interface
type MockServerControl struct {
	mock.Mock
}

// Shutdown is a mock method for Shutdown
func (m *MockServerControl) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	// Simulate context deadline respected by Shutdown
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return args.Error(0)
}

func TestHandleShutdown_Successful(t *testing.T) {
	mockSrv := new(MockServerControl)
	ctx := context.Background()
	shutdownTimeout := 100 * time.Millisecond

	// Expect Shutdown to be called with a context that has a deadline derived from shutdownTimeout
	mockSrv.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()

	err := handleShutdown(ctx, mockSrv, shutdownTimeout)

	assert.NoError(t, err)
	mockSrv.AssertExpectations(t)
}

func TestHandleShutdown_ShutdownError(t *testing.T) {
	mockSrv := new(MockServerControl)
	ctx := context.Background()
	shutdownTimeout := 100 * time.Millisecond
	expectedErr := errors.New("shutdown failed")

	mockSrv.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).Return(expectedErr).Once()

	err := handleShutdown(ctx, mockSrv, shutdownTimeout)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockSrv.AssertExpectations(t)
}

func TestHandleShutdown_AlreadyClosed(t *testing.T) {
	mockSrv := new(MockServerControl)
	ctx := context.Background()
	shutdownTimeout := 100 * time.Millisecond

	mockSrv.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).Return(http.ErrServerClosed).Once()

	err := handleShutdown(ctx, mockSrv, shutdownTimeout)

	assert.NoError(t, err) // ErrServerClosed is treated as a successful shutdown
	mockSrv.AssertExpectations(t)
}

func TestHandleShutdown_Timeout(t *testing.T) {
	mockSrv := new(MockServerControl)
	ctx := context.Background()
	shutdownTimeout := 50 * time.Millisecond // Short timeout for the test

	// Simulate Shutdown taking longer than the timeout
	mockSrv.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).
		Run(func(args mock.Arguments) {
			// Wait for the context passed to Shutdown to be canceled by the timeout in handleShutdown
			ctxArg := args.Get(0).(context.Context)
			<-ctxArg.Done()
		}).
		Return(context.DeadlineExceeded). // http.Server.Shutdown returns ctx.Err() when its context is done.
		Once()

	err := handleShutdown(ctx, mockSrv, shutdownTimeout)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded), "Expected context.DeadlineExceeded, got %v", err)
	mockSrv.AssertExpectations(t)
}

func TestHandleShutdown_ParentContextCancellation(t *testing.T) {
	mockSrv := new(MockServerControl)
	// Parent context that we will cancel
	parentCtx, cancelParent := context.WithCancel(context.Background())
	shutdownTimeout := 200 * time.Millisecond // Ample time, cancellation should be quicker

	// Expect Shutdown to be called. The context it receives (shutdownOpCtx) will be a child of parentCtx.
	mockSrv.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).
		Run(func(args mock.Arguments) {
			ctxArg := args.Get(0).(context.Context) // This is shutdownOpCtx
			// Wait until this context is canceled. In this test, it will be due to parentCtx cancellation.
			<-ctxArg.Done()
		}).
		Return(context.Canceled). // When the context passed to Shutdown is canceled.
		Once()

	// Cancel the parent context shortly after calling handleShutdown
	go func() {
		time.Sleep(50 * time.Millisecond) // Allow handleShutdown to set up its own context
		cancelParent()                    // This cancellation should propagate to shutdownOpCtx
	}()

	err := handleShutdown(parentCtx, mockSrv, shutdownTimeout)

	assert.Error(t, err)
	// The error returned by srv.Shutdown when its context is canceled is context.Canceled or context.DeadlineExceeded.
	// Since we trigger parent cancellation, and our mock returns context.Canceled, we expect that.
	assert.True(t, errors.Is(err, context.Canceled), "Expected context.Canceled, got %v", err)
	mockSrv.AssertExpectations(t)
}

func TestHandleShutdown_ZeroShutdownTimeout(t *testing.T) {
	mockSrv := new(MockServerControl)
	ctx := context.Background()
	shutdownTimeout := 0 * time.Millisecond // Zero timeout

	// With a zero timeout, context.WithTimeout creates an already-canceled context.
	// So, Shutdown should be called with an already-canceled context.
	mockSrv.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).
		Run(func(args mock.Arguments) {
			ctxArg := args.Get(0).(context.Context)
			assert.Error(t, ctxArg.Err(), "Context passed to Shutdown should be already canceled")
			assert.True(t, errors.Is(ctxArg.Err(), context.DeadlineExceeded) || errors.Is(ctxArg.Err(), context.Canceled), "Context error should be DeadlineExceeded or Canceled")
		}).
		Return(context.DeadlineExceeded). // http.Server.Shutdown would return ctx.Err().
		Once()

	err := handleShutdown(ctx, mockSrv, shutdownTimeout)

	assert.Error(t, err)
	// Expect DeadlineExceeded because WithTimeout(parent, 0) results in a context that is immediately done with DeadlineExceeded.
	assert.True(t, errors.Is(err, context.DeadlineExceeded), "Expected context.DeadlineExceeded for zero timeout, got %v", err)
	mockSrv.AssertExpectations(t)
}

package booking

import (
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"cinema-booking/internal/booking/adapters/redis"
	"github.com/google/uuid"
)

func TestConcurrentBooking_ExactlyOneWins(t *testing.T) {
	store := NewRedisStore(redis.NewClientFromEnv())
	svc := NewService(store)

	// Default concurrency is intentionally small for CI/local runs.
	// Override with CONCURRENCY env var when you want to run heavier tests.
	concurrency := 1000
	if v := os.Getenv("CONCURRENCY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			concurrency = n
		}
	}

	var (
		successes atomic.Int64
		failures  atomic.Int64
		wg        sync.WaitGroup
	)

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(userNum int) {
			defer wg.Done()
			_, err := svc.Book(Booking{
				MovieID: "screen-1",
				SeatID:  "A1",
				UserID:  uuid.New().String(),
			})
			if err == nil {
				successes.Add(1)
			} else {
				failures.Add(1)
			}
		}(i)
	}
	wg.Wait()

	if got := successes.Load(); got != 1 {
		t.Errorf("expected exactly 1 success, got %d", got)
	}
	if got := failures.Load(); got != int64(concurrency-1) {
		t.Errorf("expected %d failures, got %d", concurrency-1, got)
	}
}

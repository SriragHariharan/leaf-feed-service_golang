package seed

import (
	"fmt"
	"math/rand/v2"
	"time"
)

const ninetyDays = 90 * 24 * time.Hour

// randomTime returns a timestamp uniformly distributed in [now-ninetyDays, now].
func randomTime(r *rand.Rand, now time.Time) time.Time {
	if r == nil {
		panic("seed: randomTime requires non-nil *rand.Rand")
	}
	delta := time.Duration(r.Int64N(int64(ninetyDays)))
	return now.Add(-delta)
}

// randomUsers picks one user ID from pool. pool must be non-empty.
func randomUsers(pool []string, r *rand.Rand) string {
	if len(pool) == 0 {
		panic("seed: randomUsers requires non-empty pool")
	}
	return pool[r.IntN(len(pool))]
}

// generateMediaURL returns a stable-looking CDN URL suitable for load tests.
func generateMediaURL(r *rand.Rand) string {
	w := 640 + r.IntN(480)
	h := 360 + r.IntN(360)
	seed := r.Uint64()
	return fmt.Sprintf("https://picsum.photos/seed/%d/%d/%d.jpg", seed, w, h)
}

func randomBool(r *rand.Rand, probability float64) bool {
	if probability <= 0 {
		return false
	}
	if probability >= 1 {
		return true
	}
	return r.Float64() < probability
}

package applicant

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRateLimit(t *testing.T) {
	tests := []struct {
		name  string
		burst int
		rate  rate.Limit
	}{
		{
			name:  "test1",
			burst: 300,
			rate:  rate.Limit(float64(1) / float64(20)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := rate.NewLimiter(tt.rate, tt.burst)
			if rl.Burst() != tt.burst {
				t.Errorf("Burst() = %v, want %v", rl.Burst(), tt.burst)
			}
			if rl.Limit() != tt.rate {
				t.Errorf("Limit() = %v, want %v", rl.Limit(), tt.rate)
			}

			t.Log("consume all tokens at once", rl.AllowN(time.Now(), tt.burst))

			t.Log("consume more", rl.Allow())

			time.Sleep(time.Second * 5)
			t.Log("consume after 5 seconds", rl.Allow())

			time.Sleep(time.Second * 20)
			t.Log("consume after 20 seconds", rl.Allow())
		})
	}
}

package cache

import (
	"context"
	"testing"
	"time"
)

func getClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewService(&Config{
		Host:         "0.0.0.0",
		Port:         "6379",
		PoolSize:     20,
		IdleTimeout:  time.Duration(5) * time.Second,
		ReadTimeout:  time.Duration(5) * time.Second,
		WriteTimeout: time.Duration(5) * time.Second,
		DialTimeout:  time.Duration(5) * time.Second,
	})
	if err != nil {
		t.Skip("cannot connect with redis")
	}

	return client
}

func TestCacheMiss(t *testing.T) {
	cache := getClient(t)
	ctx := context.Background()

	if err := cache.Get(ctx, "doesnotexist", nil); err != ErrCacheMiss {
		t.Errorf("got %v; want ErrCacheMiss", err)
	}
}

func TestSetGetDelete(t *testing.T) {
	testCases := []struct {
		key   string
		value string
	}{
		{"keytest-1", "valuetest1"},
		{"keytest-2", "valuetest2"},
		{"keytest-3", "valuetest3"},
	}

	cache := getClient(t)
	ctx := context.Background()

	for _, test := range testCases {
		t.Run("Set Get Delete", func(t *testing.T) {
			item := &Item{Key: test.key, Value: test.value}
			if err := cache.Set(ctx, item); err != nil {
				t.Fatalf("Set: %v", err)
			}

			var value string
			if err := cache.Get(ctx, test.key, &value); err != nil {
				t.Fatalf("Error: %v", err)
			}
			if value != test.value {
				t.Errorf("Get: got %q, want %q", value, test.value)
			}

			if err := cache.Delete(ctx, test.key); err != nil {
				t.Fatalf("Delete: %v", err)
			}
		})
	}
}

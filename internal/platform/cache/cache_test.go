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

	var destination interface{} // to moze byc inny typ
	if err := cache.Get(ctx, "doesnotexist", destination); err != ErrCacheMiss {
		t.Errorf("got %v; want ErrCacheMiss", err)
	}
}

// func TestSetGetDelete(t *testing.T) {
// 	cache := getClient(t)
// 	ctx := context.Background()

// 	key := "testKey123456789"
// 	value := "value"

// 	item := &Item{Key: key, Value: value}
// 	if err := cache.Set(ctx, item); err != nil {
// 		t.Fatalf("Set: %v", err)
// 	}

// 	// jest jakis clean up po testach, zeby smieci nie bylo w cache?

// 	var newValue string
// 	if err := cache.Get(ctx, key, &newValue); err == ErrCacheMiss {
// 		t.Errorf("Get: got %v, want value for key: %v", err, key)
// 	}
// }

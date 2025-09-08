package retry

import (
	"errors"
	"fmt"
	"testing"
)

func TestRetry(t *testing.T) {
	t.Run("succeeds first try", func(t *testing.T) {
		calls := 0
		result, err := Retry(3, 0, func() (string, error) {
			calls++
			return "ok", nil
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result != "ok" {
			t.Fatalf("expected 'ok', got %v", result)
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d", calls)
		}
	})

	t.Run("succeeds after retries", func(t *testing.T) {
		calls := 0
		result, err := Retry(3, 0, func() (int, error) {
			calls++
			if calls < 3 {
				return 0, errors.New("fail")
			}
			return 42, nil
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result != 42 {
			t.Fatalf("expected 42, got %v", result)
		}
		if calls != 3 {
			t.Fatalf("expected 3 calls, got %d", calls)
		}
	})

	t.Run("always fails", func(t *testing.T) {
		calls := 0
		_, err := Retry(4, 0, func() (struct{}, error) {
			calls++
			return struct{}{}, fmt.Errorf("fail %d", calls)
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if calls != 4 {
			t.Fatalf("expected 4 calls, got %d", calls)
		}
	})

	t.Run("returns value type", func(t *testing.T) {
		calls := 0
		type myStruct struct{ X int }
		val, err := Retry(2, 0, func() (myStruct, error) {
			calls++
			return myStruct{X: 99}, nil
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if val.X != 99 {
			t.Fatalf("expected X=99, got %v", val.X)
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d", calls)
		}
	})
}

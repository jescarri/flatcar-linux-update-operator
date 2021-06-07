package agent

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func Test_slitNewlineEnv_retain_map_values_when_given_empty_input(t *testing.T) {
	t.Parallel()

	expected := map[string]string{"foo": "bar"}

	input := map[string]string{}

	splitNewlineEnv(input, "foo=bar")

	splitNewlineEnv(input, "")

	if !reflect.DeepEqual(expected, input) {
		t.Fatalf("Should retain original map when given empty content")
	}
}

func Test_sleepOrDone_returns_when_given(t *testing.T) {
	t.Parallel()

	t.Run("channel_is_closed", func(t *testing.T) {
		t.Parallel()

		stop := make(chan struct{})
		close(stop)

		d := time.Now().Add(1 * time.Second)
		ctxWithDeadline, cancelWithDeadline := context.WithDeadline(context.Background(), d)

		t.Cleanup(cancelWithDeadline)

		out := make(chan struct{})

		go func() {
			sleepOrDone(5*time.Second, stop)

			out <- struct{}{}
		}()

		select {
		case <-out:
		case <-ctxWithDeadline.Done():
			t.Fatalf("Waiting for function to return took too long")
		}
	})

	t.Run("timeout_is_reached", func(t *testing.T) {
		t.Parallel()

		d := time.Now().Add(5 * time.Second)
		ctxWithDeadline, cancelWithDeadline := context.WithDeadline(context.Background(), d)

		t.Cleanup(cancelWithDeadline)

		out := make(chan struct{})

		go func() {
			stop := make(chan struct{})
			sleepOrDone(1*time.Microsecond, stop)

			out <- struct{}{}
		}()

		select {
		case <-out:
		case <-ctxWithDeadline.Done():
			t.Fatalf("Waiting for function to return took too long")
		}
	})
}

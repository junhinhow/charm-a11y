package a11y

import (
	"reflect"
	"testing"
)

func TestFrames(t *testing.T) {
	frames := []string{"a", "b", "c"}

	t.Run("com movimento", func(t *testing.T) {
		if got := (Profile{}).Frames(frames); !reflect.DeepEqual(got, frames) {
			t.Fatalf("Frames() = %#v, quer %#v", got, frames)
		}
	})

	t.Run("movimento reduzido", func(t *testing.T) {
		got := (Profile{ReducedMotion: true}).Frames(frames)
		if !reflect.DeepEqual(got, []string{"a"}) {
			t.Fatalf("Frames() = %#v, quer [a]", got)
		}
	})

	t.Run("nil permanece nil", func(t *testing.T) {
		if got := (Profile{ReducedMotion: true}).Frames(nil); got != nil {
			t.Fatalf("Frames(nil) = %#v, quer nil", got)
		}
		if got := (Profile{}).Frames(nil); got != nil {
			t.Fatalf("Frames(nil) = %#v, quer nil", got)
		}
	})

	t.Run("vazio permanece vazio", func(t *testing.T) {
		empty := make([]string, 0)
		for _, p := range []Profile{{}, {ReducedMotion: true}} {
			got := p.Frames(empty)
			if got == nil || len(got) != 0 {
				t.Fatalf("Frames(vazio) = %#v, quer slice vazio não nil", got)
			}
		}
	})

	t.Run("um frame", func(t *testing.T) {
		one := []string{"só"}
		if got := (Profile{ReducedMotion: true}).Frames(one); !reflect.DeepEqual(got, one) {
			t.Fatalf("Frames() = %#v, quer %#v", got, one)
		}
	})

	t.Run("entrada não é alterada", func(t *testing.T) {
		original := append([]string(nil), frames...)
		_ = (Profile{ReducedMotion: true}).Frames(frames)
		if !reflect.DeepEqual(frames, original) {
			t.Fatalf("Frames alterou a entrada: %#v", frames)
		}
	})
}

func TestFramesDeterministic(t *testing.T) {
	frames := []string{"⠋", "⠙", "⠹", "⠸"}
	p := Profile{ReducedMotion: true}
	want := p.Frames(frames)
	for i := 0; i < 100; i++ {
		if got := p.Frames(frames); !reflect.DeepEqual(got, want) {
			t.Fatalf("iteração %d: Frames() = %#v, quer %#v", i, got, want)
		}
	}
}

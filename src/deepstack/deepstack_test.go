package deepstack

import (
	"fallbackmap"
	"reflect"
	"testing"
)

func Test_Stack(t *testing.T) {
	stack := DeepStack{}
	frameOne := fallbackmap.DeepMap(map[string]interface{}{
		"a": "one:a",
		"b": "one:b",
	})
	frameTwo := fallbackmap.DeepMap(map[string]interface{}{
		"a": "two:a",
		"d": "two:d",
	})

	stack.Push(frameOne)
	stack.Push(frameTwo)

	value, _ := stack.Get([]string{"a"})
	if value != "two:a" {
		t.Fatalf("overridden 'a' value not retrieved (%#v instead of %#v)", value, "two:a")
	}

	value, _ = stack.Get([]string{"d"})
	if value != "two:d" {
		t.Fatalf("appended 'd' value not retrieved (%#v instead of %#v)", value, "two:d")
	}

	value, _ = stack.Get([]string{"b"})
	if value != "one:b" {
		t.Fatalf("original 'b' value not retrieved (%#v instead of %#v)", value, "one:b")
	}

	_, hasKey := stack.Get([]string{"c"})
	if hasKey {
		t.Fatalf("undefined 'c' value claims to be set")
	}

	aFrame := stack.Pop()
	if !reflect.DeepEqual(aFrame, frameTwo) {
		t.Fatalf("Pop did not return the most-recent Frame (%#v instead)", aFrame)
	}

	value, _ = stack.Get([]string{"a"})
	if value != "one:a" {
		t.Fatalf("original 'a' value not retrieved after Pop (%#v instead of %#v)", value, "one:a")
	}

	_, hasKey = stack.Get([]string{"d"})
	if hasKey {
		t.Fatalf("appended 'd' value still present after Pop")
	}

	stack.Push(fallbackmap.DeepMap(frameTwo))

	value, _ = stack.Get([]string{"a"})
	if value != "two:a" {
		t.Fatalf("overridden 'a' value not retrieved [again] (%#v instead of %#v)", value, "two:a")
	}

	stack.PopDiscard()
	value, _ = stack.Get([]string{"a"})
	if value != "one:a" {
		t.Fatalf("original 'a' value not retrieved after PopDiscard (%#v instead of %#v)", value, "one:a")
	}
}

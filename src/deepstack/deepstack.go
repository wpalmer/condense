package deepstack

import (
	"fallbackmap"
)

type DeepStack struct {
	frames []fallbackmap.Deep
}

func (stack *DeepStack) Push(frame fallbackmap.Deep) {
	stack.frames = append(stack.frames, frame)
}

func (stack *DeepStack) Pop() fallbackmap.Deep {
	frame := stack.frames[len(stack.frames)-1]
	stack.frames = stack.frames[:len(stack.frames)-1]
	return frame
}

func (stack *DeepStack) PopDiscard() {
	_ = stack.Pop()
}

func (stack *DeepStack) Get(path []string) (value interface{}, has_key bool) {
	i := len(stack.frames) - 1
	for ; i >= 0; i-- {
		frame := stack.frames[i]

		found_value, did_find := frame.Get(path)
		if did_find {
			return found_value, did_find
		}
	}

	return nil, false
}

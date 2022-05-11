package tree

import "testing"

func TestBaseHook(t *testing.T) {
	hook := NewBaseHook(WithHookName("base-hook"), WithHookPriority(27))

	if hook.Name() != "base-hook" || hook.Priority() != 27 {
		t.Error("Spec error")
	}
}

package helpers

import "testing"

func TestNormalizeLineTime(t *testing.T) {
	t.Run("TestNormalizeLineTime", func(t *testing.T) {
		line := &Line{
			Time: "[25/04/25, 12:15:13]",
		}
		expected := &Line{
			Time: "[25/04/25 12h]",
		}
		result := NormalizeLineTime(line)
		if result.Time != expected.Time {
			t.Errorf("Expected '%s', but got '%s'", expected.Time, result.Time)
		}
	})
}

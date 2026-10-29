// Этот тестовый файл сгенерирован автоматически.
package usecase

import "testing"

func TestAverage(t *testing.T) {
	got := average([]float64{20, 40, 60})
	if got != 40 {
		t.Fatalf("average() = %v; want 40", got)
	}

	if got := average(nil); got != -1 {
		t.Fatalf("average(nil) = %v; want -1", got)
	}
}

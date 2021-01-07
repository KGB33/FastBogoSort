package FastBogo

import (
	"fmt"
	"sort"
	"testing"
)

func TestBogoSortDispatcher(t *testing.T) {
	list := []int{2, 3, 1, 5, 4}
	expected := []int{1, 2, 3, 4, 5}

	result := BogoSortDispatcher(list)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("List not sorted properly, got %v, expected %v", result, expected)
		}
	}

}

func TestBogoSort(t *testing.T) {
	list := []int{1} // using one element lists to ensure its sorted
	expected := []int{1}

	result := make(chan []int, 1)
	BogoSort(list, result)
	for len(result) != 1 {
	} // wait for BogoSort to finish
	for i, v := range <-result {
		if v != expected[i] {
			t.Errorf("Could not sort 1 element list")
		}
	}
}

func TestLongBogoSort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	list := []int{
		-956334, 747183, -277354, -225607, 697460, -632270, 749379, 116429, 620957, 841872,
		690465, -798563, -423627, -758154, 309200, -716938, 751019, -379879, -376340, 574886,
		672551, 584498, -623535, -212267, -937714, -647862, 44922, -455400, 8565, 520999,
	}
	expected := append([]int{}, list...)
	sort.Ints(expected)

	result := BogoSortDispatcher(list)
	fmt.Println(result)
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Result not sorted correctly, got %d at index %d, expected %d", result[i], i, expected[i])
		}
	}
}
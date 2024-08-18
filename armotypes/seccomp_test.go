package armotypes

import (
	"reflect"
	"testing"
)

func TestRemoveItems(t *testing.T) {
	tests := []struct {
		name     string
		list1    []string
		list2    []string
		expected []string
	}{
		{
			name:     "No common elements",
			list1:    []string{"apple", "banana", "cherry"},
			list2:    []string{"date", "fig", "grape"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "Some common elements",
			list1:    []string{"apple", "banana", "cherry"},
			list2:    []string{"banana", "cherry"},
			expected: []string{"apple"},
		},
		{
			name:     "All elements removed",
			list1:    []string{"apple", "banana", "cherry"},
			list2:    []string{"apple", "banana", "cherry"},
			expected: []string{},
		},
		{
			name:     "Empty list1",
			list1:    []string{},
			list2:    []string{"banana", "cherry"},
			expected: []string{},
		},
		{
			name:     "Empty list2",
			list1:    []string{"apple", "banana", "cherry"},
			list2:    []string{},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "Both lists empty",
			list1:    []string{},
			list2:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveItems(tt.list1, tt.list2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAppendUniqueStrings(t *testing.T) {
	tests := []struct {
		name     string
		dst      []string
		src      []string
		expected []string
	}{
		{
			name:     "No duplicates, all unique",
			dst:      []string{"apple", "banana"},
			src:      []string{"cherry", "date"},
			expected: []string{"apple", "banana", "cherry", "date"},
		},
		{
			name:     "Some duplicates",
			dst:      []string{"apple", "banana"},
			src:      []string{"banana", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "All duplicates",
			dst:      []string{"apple", "banana"},
			src:      []string{"apple", "banana"},
			expected: []string{"apple", "banana"},
		},
		{
			name:     "Empty src",
			dst:      []string{"apple", "banana"},
			src:      []string{},
			expected: []string{"apple", "banana"},
		},
		{
			name:     "Empty dst",
			dst:      []string{},
			src:      []string{"apple", "banana"},
			expected: []string{"apple", "banana"},
		},
		{
			name:     "Both empty",
			dst:      []string{},
			src:      []string{},
			expected: []string{},
		},
		{
			name:     "Large input with duplicates",
			dst:      []string{"apple", "banana", "cherry"},
			src:      []string{"cherry", "date", "fig", "apple", "grape"},
			expected: []string{"apple", "banana", "cherry", "date", "fig", "grape"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendUniqueStrings(tt.dst, tt.src)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

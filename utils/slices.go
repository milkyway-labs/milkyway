package utils

import (
	"slices"
)

// Find returns the first element in the slice that satisfies the given predicate, or false if none is found.
func Find[T any](slice []T, searchFunction func(T) bool) (*T, bool) {
	for i, v := range slice {
		if searchFunction(v) {
			return &slice[i], true
		}
	}
	return nil, false
}

// FindDuplicate returns the first duplicate element in the slice.
// If no duplicates are found, it returns nil instead.
func FindDuplicate[T comparable](slice []T) *T {
	return FindDuplicateFunc(slice, func(a, b T) bool {
		return a == b
	})
}

// FindDuplicateFunc returns the first duplicate element in the slice.
// If no duplicates are found, it returns nil instead.
func FindDuplicateFunc[T any](slice []T, compare func(a T, b T) bool) *T {
	for i, a := range slice {
		for j, b := range slice {
			if i != j && compare(a, b) {
				return &a
			}
		}
	}
	return nil
}

// Map applies the given function to each element in the slice and returns a new slice with the results.
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// MapWithErr applies the given function to each element in the slice and returns a new slice with the results.
func MapWithErr[T, U any](slice []T, f func(T) (U, error)) ([]U, error) {
	result := make([]U, len(slice))
	for i, v := range slice {
		res, err := f(v)
		if err != nil {
			return nil, err
		}
		result[i] = res
	}
	return result, nil
}

// RemoveDuplicates removes all duplicate elements from the slice.
func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// RemoveDuplicatesFunc removes all duplicate elements from the slice based
// on the value returned by the provided function.
func RemoveDuplicatesFunc[T any, C comparable](slice []T, compareBy func(T) C) []T {
	seen := make(map[C]bool)
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		compareValue := compareBy(v)
		if _, ok := seen[compareValue]; !ok {
			seen[compareValue] = true
			result = append(result, v)
		}
	}
	return result
}

// Remove removes the first instance of value from the provided slice.
func Remove[T comparable](slice []T, value T) (newSlice []T, removed bool) {
	index := -1
	for i, v := range slice {
		if v == value {
			index = i
		}
	}

	if index == -1 {
		return slice, false
	}

	return append(slice[:index], slice[index+1:]...), true
}

// Intersect returns the intersection of two slices (A ∩ B).
// Examples:
// Intersect([]int{1, 2, 3}, []int{2, 3, 4}) => []int{2, 3}
func Intersect[T comparable](a, b []T) []T {
	var result []T
	for _, v := range a {
		if slices.Contains(b, v) {
			result = append(result, v)
		}
	}

	return result
}

// Difference returns the elements that are in a but not in b (A - B).
// Examples:
// Difference([]int{1, 2, 3}, []int{2, 3, 4}) => []int{1}
func Difference[T comparable](a, b []T) []T {
	var result []T
	for _, v := range RemoveDuplicates(a) {
		if !slices.Contains(b, v) {
			result = append(result, v)
		}
	}
	return result
}

// Union returns the union of two slices (A ∪ B).
// Examples:
// Union([]int{1, 2, 3}, []int{2, 3, 4}) => []int{1, 2, 3, 4}
func Union[T comparable](a, b []T) []T {
	return RemoveDuplicates(append(a, b...))
}

// Filter returns the elements of the slice that satisfy the given predicate.
func Filter[T any](slice []T, f func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

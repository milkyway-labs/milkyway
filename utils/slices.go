package utils

import (
	"slices"
)

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

// Intersect returns the intersection of two slices.
func Intersect[T comparable](a, b []T) []T {
	var result []T
	for _, v := range a {
		if slices.Contains(b, v) {
			result = append(result, v)
		}
	}

	return result
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

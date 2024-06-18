package utils

// FindDuplicate returns the first duplicate element in the slice.
// If no duplicates are found, it returns nil instead.
func FindDuplicate[T any](slice []T, compare func(a T, b T) bool) *T {
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

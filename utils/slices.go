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

package arrays

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Filter[T any](ts []T, f func(T) bool) []T {
	us := []T{}
	for i := range ts {
		if pass := f(ts[i]); pass {
			us = append(us, ts[i])
		}
	}
	return us
}

func Generate[T any](length int) []T {
	return make([]T, length)
}

package repository

func Chunks[T any](s []T, size int) [][]T {
	if size < 1 {
		return nil
	}

	result := make([][]T, 0, len(s)/size+1)
	for i := 0; i < len(s); i += size {
		begin := i
		end := i + size
		if len(s) < end {
			end = len(s)
		}
		result = append(result, s[begin:end])
	}
	return result
}

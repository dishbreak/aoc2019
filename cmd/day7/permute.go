package main

func Permute(values []int) [][]int {
	result := make([][]int, 0)
	result = permute(values, 0, len(values)-1, result)
	return result
}

func permute(values []int, j, n int, acc [][]int) [][]int {
	if j == n {
		result := make([]int, len(values))
		copy(result, values)
		return append(acc, result)

	}

	for i := j; i <= n; i++ {
		values[j], values[i] = values[i], values[j]
		acc = permute(values, j+1, n, acc)
		values[j], values[i] = values[i], values[j]
	}
	return acc
}

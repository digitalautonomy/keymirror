package ssh

func transform[T, U any](l []T, f func(T) U) []U {
	output := []U{}
	for _, l := range l {
		output = append(output, f(l))
	}
	return output
}

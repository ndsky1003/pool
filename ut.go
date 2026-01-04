package pool

func resolveOption[T any](old **T, new *T) {
	if new != nil {
		*old = new
	}
}

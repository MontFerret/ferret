package spec

func Compose(base Spec, extenders ...Spec) Spec {
	out := base

	for _, ext := range extenders {
		out = out.Merge(ext)
	}

	return out
}

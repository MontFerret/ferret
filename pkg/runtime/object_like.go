package runtime

// ObjectLike marks values that should be treated as TypeObject.
// It allows VM-internal object implementations to participate
// in object-specific checks without being in the runtime package.
type ObjectLike interface {
	Map
	ObjectLike()
}

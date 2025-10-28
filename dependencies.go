package grove

// A key that you define to give a name to your dependency while registering.
type DependencyKey string

// Dependencies is a container that provides methods to pull Dependencies at runtime, usually during
// bootstrapping.
type Dependencies struct {
	deps map[DependencyKey]any
}

// Initializes the Dependencies struct.
func NewDependencies() *Dependencies {
	return &Dependencies{
		deps: make(map[DependencyKey]any),
	}
}

// Adds a dependency to the container that can be pulled later.
func (d *Dependencies) Set(key DependencyKey, value any) {
	if d.deps == nil {
		return
	}
	d.deps[key] = value
}

// Pulls a dependency from the Dependencies container. If the dependency is not found the application will panic.
func DependencyMustGet[T any](deps *Dependencies, key DependencyKey) T {
	value, ok := deps.deps[key]
	if !ok {
		panic("dependency not found: " + string(key))
	}
	result, ok := value.(T)
	if !ok {
		panic("dependency type mismatch for key: " + string(key))
	}
	return result
}

// Pulls a dependency from the Dependencies container. If it is not found it will return false and default value
// for the dependency type provided.
func DependencyGet[T any](deps *Dependencies, key DependencyKey) (T, bool) {
	value, ok := deps.deps[key]
	if !ok {
		var zero T
		return zero, false
	}
	result, ok := value.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return result, true
}

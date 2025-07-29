package grove

type DependencyKey string

type Dependencies struct {
	deps map[DependencyKey]interface{}
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		deps: make(map[DependencyKey]interface{}),
	}
}

func (d *Dependencies) Set(key DependencyKey, value interface{}) {
	if d.deps == nil {
		d.deps = make(map[DependencyKey]interface{})
	}
	d.deps[key] = value
}

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

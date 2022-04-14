package module

// modules is a set containing available scan modules,
// either those set by default (below in init), or any
// added through AddModule().
var modules = NewSet()

// Set is a container holding named scan modules.
type Set struct {
	modules      map[string]Scan
	defaultPorts map[string]int
}

// NewSet returns an empty ModuleSet.
func NewSet() Set {
	return Set{
		modules:      make(map[string]Scan),
		defaultPorts: make(map[string]int),
	}
}

// GetModule returns a scan module, along with its default port.
func GetModule(name string) Scan {
	return modules.modules[name]
}

// GetAll returns all the available scan modules and their default ports.
func GetAll() (map[string]Scan, map[string]int) {
	return modules.modules, modules.defaultPorts
}

// AddModule adds m to the ModuleSet, accessible via the given name. If the name
// is already in the ModuleSet, it is overwritten.
func AddModule(name string, defaultPort int, m Scan) {
	modules.modules[name] = m
	modules.defaultPorts[name] = defaultPort
}

// RemoveModule removes the module at the specified name. If the name is not in
// the module set, nothing happens.
func RemoveModule(name string) {
	delete(modules.modules, name)
	delete(modules.defaultPorts, name)
}

// CopyInto copies the modules in s to destination. The sets will be unique, but
// the underlying ScanModule instances will be the same.
func (s Set) CopyInto(destination Set) {
	// In case someone passes an empty module set
	if destination.modules == nil {
		destination.modules = make(map[string]Scan)
	}
	if destination.defaultPorts == nil {
		destination.defaultPorts = make(map[string]int)
	}
	// Copy modules
	for name, m := range s.modules {
		destination.modules[name] = m
		destination.defaultPorts[name] = s.defaultPorts[name]
	}
}

package module

import (
	"os"
	"path/filepath"
	"strings"
)

// All valid module files must have this extension
const moduleExtension = ".hcl"

// ModuleRegistry type contains discovered modules as returned by the
// DiscoverModules() function.
// Keys of the map are the module names and their values are the
// absolute path to the discovered module files
type ModuleRegistry map[string]string

// NewModuleRegistry creates a new empty module registry
func NewModuleRegistry() ModuleRegistry {
	registry := make(map[string]string)

	return registry
}

// Discover is used to discover valid modules in a given module path
func Discover(root string) (ModuleRegistry, error) {
	registry := NewModuleRegistry()

	// Module walker function
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directory entries
		if info.IsDir() {
			return nil
		}

		// Skip files which don't appear to be valid module files
		if filepath.Ext(info.Name()) != moduleExtension {
			return nil
		}

		// Remove the root path portion from the discovered module file,
		// remove the module file extension and register the module
		moduleFileWithExt := strings.TrimPrefix(path, root)
		moduleNameWithExt := strings.TrimSuffix(moduleFileWithExt, moduleExtension)
		moduleName := strings.Trim(moduleNameWithExt, string(os.PathSeparator))
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		registry[moduleName] = absPath

		return nil
	}

	err := filepath.Walk(root, walker)
	if err != nil {
		return registry, err
	}

	return registry, nil
}

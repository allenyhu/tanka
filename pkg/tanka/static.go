package tanka

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// StaticLoader loads an environment from a static file called `spec.json`.
// Jsonnet is evaluated as normal
type StaticLoader struct{}

func (s StaticLoader) Load(path string, opts JsonnetOpts) (envs []*v1alpha1.Environment, err error) {
	root, base, err := jpath.Dirs(path)
	if err != nil {
		return nil, err
	}

	env, err := parseStaticSpec(root, base)
	if err != nil {
		return nil, err
	}

	data, err := EvalJsonnet(path, opts)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), &env.Data); err != nil {
		return nil, err
	}

	return append(envs, env), nil
}

// parseStaticSpec parses the `spec.json` of the environment and returns a
// *kubernetes.Kubernetes from it
func parseStaticSpec(root, base string) (*v1alpha1.Environment, error) {
	// name of the environment: relative path from rootDir
	name, err := filepath.Rel(root, base)
	if err != nil {
		return nil, err
	}

	config, err := spec.ParseDir(base, name)
	if err != nil {
		switch err.(type) {
		// the config includes deprecated fields
		case spec.ErrDeprecated:
			log.Println(err)
		// spec.json missing. we can still work with the default value
		case spec.ErrNoSpec:
			return config, nil
		// some other error
		default:
			return nil, err
		}
	}

	return config, nil
}

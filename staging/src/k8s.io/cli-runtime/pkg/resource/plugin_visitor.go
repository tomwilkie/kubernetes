package resource

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// PluginVisitor is a wrapper around a StreamVisitor, to handler post-processing plugins.
type PluginVisitor struct {
	Path string
	*StreamVisitor
}

// Visit in a KustomizeVisitor gets the output of Kustomize build and save it in the Streamvisitor
func (v *PluginVisitor) Visit(fn VisitorFunc) error {
	buf, err := v.runInputPlugin()
	if err != nil {
		return err
	}
	v.StreamVisitor.Reader = bytes.NewReader(buf)
	return v.StreamVisitor.Visit(fn)
}

func (v *PluginVisitor) runInputPlugin() ([]byte, error) {
	ext := filepath.Ext(v.Path)
	if ext == "" {
		return nil, fmt.Errorf("no extension for %s", v.Path)
	}

	cmd, ok := v.findPlugin(ext[1:])
	if !ok {
		return nil, fmt.Errorf("plugin not found for filetype %s", ext)
	}

	return exec.Command(cmd, v.Path).CombinedOutput()
}

func (v *PluginVisitor) findPlugin(name string) (string, bool) {
	for _, prefix := range []string{"kubectl"} { // plugin.ValidPluginFilenamePrefixes {
		path, err := exec.LookPath(fmt.Sprintf("%s-input-%s", prefix, name))
		if err != nil || len(path) == 0 {
			continue
		}
		return path, true
	}
	return "", false
}

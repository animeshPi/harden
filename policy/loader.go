package policy

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFromFile(path string) (*PolicySet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ps PolicySet
	if err := yaml.Unmarshal(data, &ps); err != nil {
		return nil, err
	}

	if len(ps.Policies) == 0 {
		return nil, fmt.Errorf("no policies found in %s", path)
	}

	return &ps, nil
}

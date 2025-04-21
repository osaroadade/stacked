package stack

import "gopkg.in/yaml.v3"

func UnmarshalStack(data []byte, out *StackData) error {
	return yaml.Unmarshal(data, out)
}

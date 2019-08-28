package kinoko

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Gene struct {
	configs map[interface{}]interface{}
	y       string
}

func (g Gene) get(path string) interface{} {
	keys := strings.Split(path, ".")
	r := g.configs
	for i := range keys {
		if r[keys[i]] == nil {
			return nil
		}
		if i == len(keys)-1 {
			return r[keys[i]]
		}
		r = r[keys[i]].(map[interface{}]interface{})
		if r == nil {
			return nil
		}
	}
	return nil
}

func (g Gene) GetMap(path string) map[interface{}]interface{} {
	if v := g.get(path); v != nil {
		return v.(map[interface{}]interface{})
	} else {
		panic("miss config: " + path)
	}
}

func (g Gene) GetInt(path string) int {
	if v := g.get(path); v != nil {
		return v.(int)
	} else {
		panic("miss config: " + path)
	}
}

func (g Gene) GetFloat(path string) float64 {
	if v := g.get(path); v != nil {
		return v.(float64)
	} else {
		panic("miss config: " + path)
	}
}

func (g Gene) GetBool(path string) bool {
	if v := g.get(path); v != nil {
		return v.(bool)
	} else {
		panic("miss config: " + path)
	}
}

func (g Gene) GetString(path string) string {
	if v := g.get(path); v != nil {
		return v.(string)
	} else {
		panic("miss config: " + path)
	}
}

func (g Gene) GetIntOrDefault(path string, def int) int {
	if v := g.get(path); v != nil {
		return v.(int)
	} else {
		return def
	}
}

func (g Gene) GetFloatOrDefault(path string, def float64) float64 {
	if v := g.get(path); v != nil {
		return v.(float64)
	} else {
		return def
	}
}

func (g Gene) GetBoolOrDefault(path string, def bool) bool {
	if v := g.get(path); v != nil {
		return v.(bool)
	} else {
		return def
	}
}

func (g Gene) GetStringOrDefault(path string, def string) string {
	if v := g.get(path); v != nil {
		return v.(string)
	} else {
		return def
	}
}

func NewGene(y string) Gene {
	config := Gene{configs: map[interface{}]interface{}{}, y: y}
	b, err := ioutil.ReadFile(y)
	if err != nil {

		logger.Error("Error reading config file", err)
		return config
	}
	err = yaml.Unmarshal(b, &config.configs)

	if err != nil {
		logger.Error("Error parsing config file", err)
	}
	return config
}

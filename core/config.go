package core

import (
	"github.com/godcong/wego/cache"
	"github.com/godcong/wego/log"
	"github.com/godcong/wego/util"
	"github.com/pelletier/go-toml"
)

const configPath = "config.toml"

func init() {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Println("no config files loaded")
		log.Error(err)
		return
	}
	cache.Set("config", config)
}

/*Config Config Tree */
type Config struct {
	*toml.Tree
}

//DefaultConfig get the default config from cache
func DefaultConfig() *Config {
	if c := cache.Get("config"); c != nil {
		if v, b := c.(*Config); b {
			return v
		}
	}
	return nil
}

/*LoadConfig get config tree with file name*/
func LoadConfig(f string) (*Config, error) {
	t, e := toml.LoadFile(f)
	if e != nil {
		log.Println("filepath: " + f)
		log.Println(e.Error())
		return nil, e
	}
	return cfg(t), nil
}

//IsNil ...
func (c *Config) IsNil() bool {
	if c == nil || c.Tree == nil {
		return true
	}
	return false
}

/*GetSubConfig get sub config from current config */
func (c *Config) GetSubConfig(s string) *Config {
	if v, b := c.GetTree(s).(*toml.Tree); b {
		return cfg(v)
	}
	return cfg(nil)
}

/*GetTree get config tree */
func (c *Config) GetTree(s string) interface{} {
	if c == nil || c.Tree == nil {
		return nil
	}

	return c.Tree.Get(s)
}

/*DeepGet get an interface from config when was in values */
func (c *Config) DeepGet(s ...string) interface{} {
	size := len(s)
	for i := 0; i < size; i++ {
		if c.Has(s[i]) {
			return c.Get(s[i])
		}
	}
	return nil
}

/*DeepGetD get an interface from config when was in values */
func (c *Config) DeepGetD(def string, s ...string) interface{} {
	size := len(s)
	for i := 0; i < size; i++ {
		if c.Has(s[i]) {
			return c.Get(s[i])
		}
	}
	return def
}

/*Get get an interface from config */
func (c *Config) Get(s string) interface{} {
	return c.GetTree(s)
}

/*GetD get interface with default value */
func (c *Config) GetD(s string, d interface{}) interface{} {
	v := c.GetTree(s)
	if v == nil {
		return d
	}
	return v
}

/*GetString get string with out default value */
func (c *Config) GetString(s string) string {
	v := c.GetTree(s)
	if v, b := v.(string); b {
		return v
	}
	return ""
}

/*GetStringD get string with default value */
func (c *Config) GetStringD(s, d string) string {
	v := c.GetTree(s)
	if v, b := v.(string); b {
		return v
	}
	return d
}

/*Set set value */
func (c *Config) Set(k string, v interface{}) *Config {
	c.Tree.Set(k, v)
	return c
}

/*GetBool get bool value */
func (c *Config) GetBool(s string) bool {
	v := c.GetTree(s)
	if v, b := v.(bool); b {
		return v
	}

	return false
}

//GetBoolD get bool with default value
func (c *Config) GetBoolD(s string, d bool) bool {
	v := c.GetTree(s)
	if v, b := v.(bool); b {
		return v
	}

	return d
}

//GetInt get int value
func (c *Config) GetInt(s string) int64 {
	v := c.GetTree(s)
	v0, b := util.ParseInt(v)
	if !b {
		return 0
	}
	return v0
}

//GetIntD get int with default value
func (c *Config) GetIntD(s string, d int64) int64 {
	v := c.GetTree(s)
	v0, b := util.ParseInt(v)
	if !b {
		return d
	}
	return v0
}

//Check check all input keys
//return -1 if all is exist
//return index when not found
func (c *Config) Check(s ...string) int {
	size := len(s)
	for i := 0; i < size; i++ {
		if !c.Has(s[i]) {
			return i
		}
	}
	return -1
}

//GetStringArray return string array
func (c *Config) GetStringArray(key string) []string {
	arr := c.GetArray(key)
	var strArr []string

	if arr != nil {
		size := len(arr)
		for i := 0; i < size; i++ {
			if vv, b := arr[i].(string); b {
				strArr = append(strArr, vv)
			}
		}
		return strArr
	}

	return nil
}

//GetStringArrayD return string array with default value
func (c *Config) GetStringArrayD(key string, d []string) []string {
	arr := c.GetArray(key)
	var strArr []string

	if arr != nil {
		size := len(arr)
		for i := 0; i < size; i++ {
			if vv, b := arr[i].(string); b {
				strArr = append(strArr, vv)
			}
		}
		return strArr
	}

	return d
}

//GetArray return array
func (c *Config) GetArray(key string) []interface{} {
	v := c.GetTree(key)
	v0, b := (v).([]interface{})
	if b {
		return v0
	}
	return nil
}

//GetArrayD return array with default value
func (c *Config) GetArrayD(key string, d []interface{}) []interface{} {
	v := c.GetTree(key)
	v0, b := (v).([]interface{})
	if b {
		return v0
	}
	return d
}

//cfg create a null config
func cfg(tree *toml.Tree) *Config {
	return &Config{
		Tree: tree,
	}
}

//NewConfig create a new null config
func NewConfig(tree *toml.Tree) *Config {
	return &Config{
		Tree: tree,
	}
}

// NilConfig ...
func NilConfig() *Config {
	cfg, err := toml.TreeFromMap(util.Map{})
	if err != nil {
		return NewConfig(nil)
	}
	return NewConfig(cfg)
}

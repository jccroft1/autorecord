package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

// Config can be used to store your apps configuration data in a file.
// It will read the file once at startup, then flush changes to the file after
// every config change.
type Config struct {
	FileName string
	Data     map[string]string
	sync.Mutex
}

var defaultConfig = New("config.txt")

func New(fileName string) *Config {
	c := Config{
		FileName: fileName,
		Data:     make(map[string]string),
	}
	c.load()
	return &c
}

func Get(key string) string {
	return defaultConfig.Get(key)
}

func (c *Config) Get(key string) string {
	c.Lock()
	defer c.Unlock()
	return c.Data[key]
}

func Set(key, value string) {
	defaultConfig.Set(key, value)
}

func (c *Config) Set(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.Data[key] = value
	c.flush()
}

func (c *Config) load() {
	contents, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return
	}

	json.Unmarshal(contents, &c.Data)
}

func (c *Config) flush() {
	contents, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		return
	}

	ioutil.WriteFile(c.FileName, contents, 0644)
}

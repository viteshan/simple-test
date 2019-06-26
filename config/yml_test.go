package config

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Yaml struct of yaml
type Yaml struct {
	Mysql struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	}
	Cache struct {
		Enable bool     `yaml:"enable"`
		List   []string `yaml:"list,flow"`
	}
}

// Yaml1 struct of yaml
type Yaml1 struct {
	SQLConf   *Mysql `yaml:"mysql"`
	CacheConf *Cache `yaml:"cache"`
}

// Yaml2 struct of yaml
type Yaml2 struct {
	Mysql `yaml:"mysql,inline"`
	Cache `yaml:"cache,inline"`
}

// Mysql struct of mysql conf
type Mysql struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

// Cache struct of cache conf
type Cache struct {
	Enable bool     `yaml:"enable"`
	List   []string `yaml:"list,flow"`
}

func TestYml(t *testing.T) {
	// resultMap := make(map[string]interface{})
	conf := new(Yaml1)
	yamlFile, err := ioutil.ReadFile("test.yaml")

	// conf := new(module.Yaml1)
	// yamlFile, err := ioutil.ReadFile("test.yaml")

	// conf := new(module.Yaml2)
	//  yamlFile, err := ioutil.ReadFile("test1.yaml")

	log.Println("yamlFile:", yamlFile)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	// err = yaml.Unmarshal(yamlFile, &resultMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	yaml1, _ := json.Marshal(conf)
	log.Println("conf", string(yaml1))
	println(conf.SQLConf == nil)
	// log.Println("conf", resultMap)
}

func TestJson(t *testing.T) {
	type Road struct {
		Name   string
		Number int
	}
	roads := []Road{
		{"Diamond Fork", 29},
		{"Sheep Creek", 51},
	}

	b, err := json.Marshal(roads)
	if err != nil {
		log.Fatalln(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")

	if err != nil {
		log.Fatalln(err)
	}

	println(out.String())

}

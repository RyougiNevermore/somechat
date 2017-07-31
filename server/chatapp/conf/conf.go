package conf

import (
	"io/ioutil"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"gopkg.in/yaml.v2"
)

var Conf *conf

type conf struct {
	Postgres 	postgres 	`yaml:"postgres"`
	Redis 		redis 		`yaml:"redis"`
	Web 		web 		`yaml:"web"`
}

func (c conf) String() string {
	if p, err := yaml.Marshal(&c); err != nil {
		log.Log().Println(logs.Errorf("conf string() failed, %v", err))
		panic(err)
		return ""
	} else {
		return string(p)
	}
}

type postgres struct {
	Url 	string	`yaml:"url"`
	MaxIdle int	`yaml:"maxIdle"`
	MaxOpen int	`yaml:"maxOpen"`
}

type redis struct {
	Address 	string	`yaml:"addr"`
	Password 	string	`yaml:"password"`
	DB 		int	`yaml:"db"`
}

type web struct {
	Port 	string `yaml:"port"`
	Static 	string `yaml:"static"`
	Tpl	string `yaml:"tpl"`
	Favicon string `yaml:"favicon"`
}




func Read(path string) error {
	p, readFileErr := ioutil.ReadFile(path)
	if readFileErr != nil {
		log.Log().Println(logs.Errorf("read conf file faild, path : %s, error : %v", path, readFileErr).Trace())
		return readFileErr
	}
	c := new(conf)
	yamlErr := yaml.Unmarshal(p, c)
	if yamlErr != nil {
		log.Log().Println(logs.Errorf("Yaml Unmarshal conf file faild, path : %s, error : %v", path, yamlErr).Trace())
		return yamlErr
	}
	Conf = c
	return nil
}
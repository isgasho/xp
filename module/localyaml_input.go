package module

import (
	"fmt"
	"reflect"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	"github.com/spf13/viper"
)

func init() {
	AddInput("localyaml", reflect.TypeOf(LocalYamlInput{}))
}

type LocalYaml struct {
	data map[string]interface{}
}

func (l *LocalYaml) Get() {
	l.data = viper.AllSettings()
}

type LocalYamlInput struct {
	status     StatusPlugin
	yaml       LocalYaml
	connecheck map[string]string
}

func (l *LocalYamlInput) Receive() *Message {
	l.yaml.Get()
	if l.status != Started {
		fmt.Println("LocalYaml input plugin is not running,input nothing.")
		return nil
	}

	return Builder().WithCheck(l.connecheck).WithItemInterface(l.yaml.data).Build()
}

func (l *LocalYamlInput) Start() {
	l.status = Started
	fmt.Println("LocalYamlInput plugin started.")

	// Check machine connecting
	ip := viper.GetStringSlice("host")

	for _, i := range ip {
		if utils.ScanPort(i, "22") {
			l.connecheck[i] = "success"
		} else {
			l.connecheck[i] = "failed"
		}
	}
}

func (l *LocalYamlInput) Stop() {
	l.status = Stopped
	fmt.Println("LocalYamlInput plugin stopped.")
}

func (l *LocalYamlInput) Status() StatusPlugin {
	return l.status
}

// LocalYamlInput的Init函数实现
func (l *LocalYamlInput) Init() {
	l.yaml.data = make(map[string]interface{})
	l.connecheck = make(map[string]string)
}

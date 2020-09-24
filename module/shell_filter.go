package module

import (
	"fmt"
	"reflect"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
)

func init() {
	// 初始化shell filter插件映射关系表
	AddFilter("shell", reflect.TypeOf(ShellFilter{}))
}

// shell 命令运行filter插件
type ShellFilter struct {
	status StatusPlugin
}

func (s *ShellFilter) Process(msgs *Message) *Message {
	if s.status != Started {
		fmt.Println("Shell filter plugin is not running,filter nothing.")
		return msgs
	}

	// 解析yaml结果
	fmt.Printf("%v\n", msgs.Data.Check)
	// 1. 解析stage步骤
	stage := msgs.Data.Items["stage"].([]interface{})
	fmt.Printf("%v\n", stage)
	config := msgs.Data.Items["config"].([]interface{})
	fmt.Printf("%v\n", config)
	// 2. 根据stage进行解析
	for host, status := range msgs.Data.Check {
		if status == "failed" {
			fmt.Printf("host %s is failed, next.\n", host)
		} else {
			fmt.Printf("执行目标主机： %s\n", host)
			for _, x := range stage {
				for _, y := range config {
					tmp := y.(map[interface{}]interface{})
					if d, ok := tmp[x.(string)]; ok {
						tmp_d := d.(map[interface{}]interface{})
						if items, ok2 := tmp_d["with_items"]; !ok2 {
							// rs, err := utils.ExecCommandString(d.(map[interface{}]interface{})["shell"].(string))
							rs, err := utils.New(host, "root", "", 22).Run(d.(map[interface{}]interface{})["shell"].(string))
							if err != nil {
								fmt.Printf("Stage: %s \nName: %s \nShell: %s \nResult: %s \n", x.(string), d.(map[interface{}]interface{})["name"].(string), d.(map[interface{}]interface{})["shell"].(string), err.Error())
							} else {
								fmt.Printf("Stage: %s \nName: %s \nShell: %s \nResult: %s \n", x.(string), d.(map[interface{}]interface{})["name"].(string), d.(map[interface{}]interface{})["shell"].(string), rs)
							}
						} else {
							for _, xx := range items.([]interface{}) {
								cmd2, err := utils.ApplyTemplate(d.(map[interface{}]interface{})["shell"].(string), map[string]interface{}{"items": []string{xx.(string)}})
								if err != nil {
									fmt.Println("cmd2 ", cmd2)
									panic(err)
								}
								fmt.Println("cmd2 ", cmd2)
								// rs, err := utils.ExecCommandString(cmd2)
								rs, err := utils.New(host, "root", "", 22).Run(cmd2)
								if err != nil {
									fmt.Printf("Stage: %s \nName: %s \nShell: %s With_items: %s \nResult: %s\n", x.(string), d.(map[interface{}]interface{})["name"].(string), d.(map[interface{}]interface{})["shell"].(string), items, err.Error())
								} else {
									fmt.Printf("Stage: %s \nName: %s \nShell: %s With_items: %s \nResult: %s\n", x.(string), d.(map[interface{}]interface{})["name"].(string), d.(map[interface{}]interface{})["shell"].(string), items, rs)
								}
							}
						}
					}
				}
			}
		}
	}

	// 3. TODO: 解析yaml中shell的模块，然后进行匹配
	return msgs
}

func (s *ShellFilter) Start() {
	s.status = Started
	fmt.Println("Shell filter plugin started.")
}

func (s *ShellFilter) Stop() {
	s.status = Stopped
	fmt.Println("Shell filter plugin stopped.")
}

func (s *ShellFilter) Status() StatusPlugin {
	return s.status
}

func (s *ShellFilter) Init() {}

# xp
xp is my name,but the project is same like ansible playbook auto ops

# 介绍

插件接口流程：

* pipeline
    * User 用户管理
    * Host 主机清单
    * Playbook
        * YAML
        * module
    * Plugin
        * start
        * stop
        * status
        * init
    * input
        * host conn check
        * host env
        * yaml module 分析
    * filter
        * 执行各个module
            * 连接
                * ssh
                * docker
                * k8s
                * 网络设备
                * snmp等
            * 执行
                * RPC
                * RESTFULL
                * CLI
        * 管理执行的生命周期
            * prepare
            * before
            * runtime
            * after
    * output
        * 输出结果
        * 返回状态

# 功能模块

- [x] yaml解析(cobra viper支持)
- [ ] 动态环境变量(cobra支持)
- [x] with_items迭代器
- [x] CLI命令行工具(cobra)
- [ ] 功能文件夹，提供：files、hosts、env等特殊目录模块
- [ ] roles ansible模块
- [ ] module man模块说明文档
- [x] module plugin插件机制
- [ ] ssh [连接功能](https://github.com/mojocn/felix)

# Useage

> go build
> ./xp test --config ./.devopsxp.yaml

# Test

> make

## 测试执行流程

cli -> main.go -> root.go -> test.go -> pipeline -> init -> start -> check(ssh) -> input(localyaml) -> filter(shell) -> output(console) -> stop

# Module

## Check Plugin

- [x] ssh

## Input Plugin

- [x] localyaml

## Filter Plugin

- [x] shell

## Output Plugin

- [x] console
## 项目简介

<p align="center">
  <a href="https://golang.google.cn/">
    <img src="https://img.shields.io/badge/Golang-1.17-green.svg" alt="golang">
  </a>
  <a href="https://gin-gonic.com/">
    <img src="https://img.shields.io/badge/Gin-1.7.4-red.svg" alt="gin">
  </a>
  <a href="https://gorm.io/">
    <img src="https://img.shields.io/badge/Gorm-1.21-orange.svg" alt="gorm">
  </a>
  <a href="https://redis.io/">
    <img src="https://img.shields.io/badge/redis-3.2.100-brightgreen.svg" alt="redis">
  </a>
  <a href="https://vuejs.org/">
    <img src="https://img.shields.io/badge/Vue-3.0.0-orange.svg" alt="vue">
  </a>
  <a href="https://antdv.com/docs/vue/introduce-cn/">
    <img src="https://img.shields.io/badge/Ant%20Design-2.2.x-blue.svg" alt="Ant Design">
  </a>
</p>

> KubeSpace运维平台项目使用Go1.15.x、 Gin、Gorm开发， 前端使用的是Vue3+Ant Design2.2.x框架。


#### 项目源码
|     | 后端源码             | 前端源码          |
|---  |------------------|---------------|
|  github   | kubespace/server | kubespace/web |


## 使用说明
1. 安装编译
```shell script
# 拉取代码
git clone git@github.com:openstack-test/kubespace.git
```

2. 启动服务前先创建server/etc/config.yaml, 数据库部分配置如下
```shell script
# 数据库配置
mysql:
  path: '127.0.0.1:3306'
  db-name: 'kubespace'
  username: 'root'
  password: '123456'
```

3. 初始化数据, 将`server/sql/init_data.sql`文件导入到mysql数据库kubespace中。


4. 启动服务
```shell
# 启动后端
cd kubespace/server
go run main.go

# 启动前端
cd kubespace/web
# 安装依赖
npm install --registry=https://registry.npm.taobao.org
npm run dev
```

5. 访问。默认用户名admin@123.com，密码123456

http://localhost:8080  

备注：先生成自定义密码的hash值, 然后替换掉users表中的password值。
``` go
package main
import (
	"golang.org/x/crypto/bcrypt"
)

func main() {
	pwd := []byte("123456")
	hash, _ := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	println("hash加密密码: ", string(hash))
}
```

#### 目前已经实现的功能
* 用户登录
  * LDAP/Email
  * 钉钉扫码登录(开源版暂未开放)
* 权限管理
* 用户注册登录
  * [如何配置LDAP](.)
  * [配置钉钉扫码](.)
- K8S多集群管理
  * [集群管理](.)
  * [节点管理](.)
  * [工作负载](.)
  * [存储管理](.)
  * [网络管理](.)
  * [配置管理](.)

- 资产管理
  * [远程连接](.)




## 首页
![avatar](./docs/img/login.png)

## 仪表盘
![avatar](./docs/img/仪表盘.jpg)

## 资产管理
![avatar](./docs/img/资产管理.jpg)

## 远程连接
![avatar](./docs/img/远程登录.png)

## 集群管理
![avatar](./docs/img/集群管理.jpg)

## 集群详情
![avatar](./docs/img/集群详情.png)

## 节点列表
![avatar](./docs/img/节点.jpg)

## 网络管理
![avatar](./docs/img/network.png)

## 存储管理
![avatar](./docs/img/storage.png)

## Roadmap

> [历史版本详情](./docs/version/README.md)
> 

- 应用发布(分批发布、版本回退)

- 作业系统
- 审批流

## License
Everything is Apache License 2.0.


## fiewall-api

firewall api是基于firewalld来远程管理iptables规则的rest-api，无需部署agent

## Features

- 指定一个主机ip，让这个主机上的iptables增加一个规则
- 处理单个IP或CIDR范围（xx.xx.xx.xx/mask，mac，interface）
- 永久生效或定时生效
- 支持服务名，端口，协议
- 支持自定义服务，模板
- 快速切换活动区
- 全dbus远程请求
- rest api
- 不会干扰其他iptables规则，如docker
- 误操作后一键恢复
- 支持nat功能
- 可选，模板入库，快速批量切换远程主机使用的模板(当开启mysql时接口才生效，不开启mysql没有此功能)

## 后期
- 增加接口鉴权
- 整个逻辑的变化
- 图形化
- 高可用
- 监控指标
- 配套的waf


## deploy

```bash
git clone ..
make
```

## use 

接口文档地址： https://documenter.getpostman.com/view/12796679/UV5agGNr

## FAQ

### 为什么不使用ssh ansible等类似工具

dbus为linux基础服务，每个主机必须存在的服务，包括用户登录等都会用到，这里基于dbus的远程功能可以不用部署agent，也无需ssh信任秘钥就可以直接使用

### 和其他的类似项目有什么区别

其他的一般是拼接命令，或者是部署agent方式进行管理，这里可以无需部署agent方式，仅一个控制端就可以无限管理大批量主机。

并且增加了reload和误操作清空功能，类似手动误操作导致不能登陆的情况可以瞬间恢复。

### 服务管理安全吗
批扣
虽说开启了dbus 远程访问，任何人都可以调用dbus去操作firewalld，可以增加iptables规则管理对应端口，以及polkit进行权限的控制，但这并不适用于firewalld


goim v2.0
==============
`Terry-Mao/goim` 是一个支持集群的im及实时推送服务。

---------------------------------------
  * [特性](#特性)
  * [安装](#安装)
  * [配置](#配置)
  * [例子](#例子)
  * [文档](#文档)
  * [集群](#集群)
  * [更多](#更多)

---------------------------------------

## 特性
 * 轻量级
 * 高性能
 * 纯Golang实现
 * 支持单个、多个、单房间以及广播消息推送
 * 支持单个Key多个订阅者（可限制订阅者最大人数）
 * 心跳支持（应用心跳和tcp、keepalive）
 * 支持安全验证（未授权用户不能订阅）
 * 多协议支持（websocket，tcp）
 * 可拓扑的架构（job、logic模块可动态无限扩展）
 * 基于Kafka做异步消息推送

 ## 架构
主要采用订阅与消费的模式进行开发，通过redis和kafka来订阅消息。其中每一个通信都是建立在room的基础上的。
使用discovery来管理各个节点的状态，也使用discovery来实现负载均衡。当客户端需要建立新的连接时，可用通过
访问discovery来获取当前负载最轻的comet，从而与comet进行连接。所以每一个需要监控的节点都需要像discovery
进行注册（在对应的main.go中的register方法中进行注册），通过对应的api就可以获取节点信息，如：
```
/discovery/fetch?zone=sh001&env=dev&status=1&appid=goim.comet

{
    "code": 0,
    "message": "0",
    "ttl": 1,
    "data": {
        "instances": {
            "sh001": [
                {
                    "region": "sh",
                    "zone": "sh001",
                    "env": "dev",
                    "appid": "goim.comet",
                    "hostname": "unitylabss-MacBook-Air.local",
                    "addrs": [
                        "grpc://192.168.12.108:3109,tcp://192.168.12.108:3101,ws://192.168.12.108:3102"
                    ],
                    "version": "",
                    "metadata": {
                        "addrs": "127.0.0.1",
                        "conn_count": "0",
                        "ip_count": "0",
                        "offline": "false",
                        "weight": "10"
                    },
                    "status": 1,
                    "reg_timestamp": 1587614849272206000,
                    "up_timestamp": 1587615139405518000,
                    "renew_timestamp": 1587615119282758000,
                    "dirty_timestamp": 1587615139405518000,
                    "latest_timestamp": 1587615139405518000
                }
            ]
        },
        "latest_timestamp": 1587614849272206000
    }
}
```
所以使用这个框架集合discovery可以轻松的二次开发，对应的discovery中提供大量有效的api接口用于监控使用。

### comet
主要用于与Client直接交互，接受来自Client的消息，以及将消息推送到Client。


### logic
主要用于消息的产生，将消息发送到kafka中。同时这里还进行权限验证。


### job
主要是通过消费kafka中的消息，然后将消息根据订阅者机制发送到comet中，由comet将消息同步到Client中。

## 安装
### 一、安装依赖
```sh
$ yum -y install java-1.7.0-openjdk
```

### 二、安装Kafka消息队列服务

kafka在官网已经描述的非常详细，在这里就不过多说明，安装、启动请查看[这里](http://kafka.apache.org/documentation.html#quickstart).

### 三、搭建golang环境
1.下载源码(根据自己的系统下载对应的[安装包](http://golang.org/dl/))
```sh
$ cd /data/programfiles
$ wget -c --no-check-certificate https://storage.googleapis.com/golang/go1.5.2.linux-amd64.tar.gz
$ tar -xvf go1.5.2.linux-amd64.tar.gz -C /usr/local
```
2.配置GO环境变量
(这里我加在/etc/profile.d/golang.sh)
```sh
$ vi /etc/profile.d/golang.sh
# 将以下环境变量添加到profile最后面
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/data/apps/go
$ source /etc/profile
```

### 四、部署goim
1.下载goim及依赖包
```sh
$ yum install hg
$ go get -u github.com/Terry-Mao/goim
$ mv $GOPATH/src/github.com/Terry-Mao/goim $GOPATH/src/goim
$ cd $GOPATH/src/goim
$ go get ./...
```

2.安装router、logic、comet、job模块(配置文件请依据实际机器环境配置)
```sh
$ cd $GOPATH/src/goim/router
$ go install
$ cp router-example.conf $GOPATH/bin/router.conf
$ cp router-log.xml $GOPATH/bin/
$ cd ../logic/
$ go install
$ cp logic-example.conf $GOPATH/bin/logic.conf
$ cp logic-log.xml $GOPATH/bin/
$ cd ../comet/
$ go install
$ cp comet-example.conf $GOPATH/bin/comet.conf
$ cp comet-log.xml $GOPATH/bin/
$ cd ../logic/job/
$ go install
$ cp job-example.conf $GOPATH/bin/job.conf
$ cp job-log.xml $GOPATH/bin/
```
到此所有的环境都搭建完成！

### 五、启动goim
```sh
$ cd /$GOPATH/bin
$ nohup $GOPATH/bin/router -c $GOPATH/bin/router.conf 2>&1 > /data/logs/goim/panic-router.log &
$ nohup $GOPATH/bin/logic -c $GOPATH/bin/logic.conf 2>&1 > /data/logs/goim/panic-logic.log &
$ nohup $GOPATH/bin/comet -c $GOPATH/bin/comet.conf 2>&1 > /data/logs/goim/panic-comet.log &
$ nohup $GOPATH/bin/job -c $GOPATH/bin/job.conf 2>&1 > /data/logs/goim/panic-job.log &
```
如果启动失败，默认配置可通过查看panic-xxx.log日志文件来排查各个模块问题.

### 六、测试

推送协议可查看[push http协议文档](./docs/push.md)

## 配置

TODO

## 例子

Websocket: [Websocket Client Demo](https://github.com/Terry-Mao/goim/tree/master/examples/javascript)

Android: [Android](https://github.com/roamdy/goim-sdk)

iOS: [iOS](https://github.com/roamdy/goim-oc-sdk)

## 文档
[push http协议文档](./docs/push.md)推送接口

## 集群

### comet

comet 属于接入层，非常容易扩展，直接开启多个comet节点，修改配置文件中的base节点下的server.id修改成不同值（注意一定要保证不同的comet进程值唯一），前端接入可以使用LVS 或者 DNS来转发

### logic

logic 属于无状态的逻辑层，可以随意增加节点，使用nginx upstream来扩展http接口，内部rpc部分，可以使用LVS四层转发

### kafka

kafka 可以使用多broker，或者多partition来扩展队列

### router

router 属于有状态节点，logic可以使用一致性hash配置节点，增加多个router节点（目前还不支持动态扩容），提前预估好在线和压力情况

### job

job 根据kafka的partition来扩展多job工作方式，具体可以参考下kafka的partition负载

##更多
TODO

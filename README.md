# 从零实现Yedis
[![MyWebSite](https://img.shields.io/badge/我的站点-whoiszxl-blue.svg)](https://whoiszxl.github.io)
[![docs](https://img.shields.io/badge/docs-reference-green.svg)](https://whoiszxl.github.io)
[![teach](https://img.shields.io/badge/教程-Monica-orange.svg)](https://github.com/whoiszxl/Monica)
[![email](https://img.shields.io/badge/email-whoiszxl@gmail.com-red.svg)](whoiszxl@gmail.com)

### 前言
Yedis是一款通过Golang语言实现的一个复刻Redis数据库，通过简单的Golang代码来更加简单地学习Redis的设计与实现，其不仅能够由浅入深地学习Redis这款优秀的键值对数据库，而且还能够学习到如何从零读一款开源软件的代码。在代码实现后，会同步出文档教程与视频教程。

### Yedis文档地址
地址：http://monica.whoiszxl.com/

### Yedis介绍
Yedis通过Golang语言实现一个简略版的Redis单机键值对数据库，其中实现了诸如SDS、LinkedList等基础数据结构，string、list、set、zset和hash五种对象类型，订阅发布、密码校验、慢查询日志等功能，RDB和AOF的数据层持久化等功能。因为Redis是使用C语言开发的，像分配SDS简单动态字符串内存释放其内存等操作Yedis都简化成使用Golang原生string，但是有些功能处不会做省略，旨在学习原生Redis的结构，比如sds的len属性，其设计是用来减少strlen函数的内存消耗，golang的string自带len属性，完全可以复用其属性，此处便不省略，为了多一点Redis的汁味。

考虑在代码开发完成后写一份图文手册，从0-1构建Yedis源码，附带Redis3.0的源码分析与不同之处，从而达到更好地理解Redis的设计思路，也能够更好地在工作中使用Redis。


### 项目环境与参考
| 技术       | 版本                  | 地址                                   |
| ---------- | --------------------- | -------------------------------------- |
| Golang     | 1.14.2                | https://golang.org/                    |
| Redis源码   | 3.2                  | https://github.com/antirez/redis/tree/3.2|
| Redis设计与实现 |          | https://book.douban.com/subject/25900156/                |

### Yedis项目演示
xx todo


### 项目构建与运行

#### Windows、Linux构建与运行
1. 将代码`git clone`到本地，并保证机器安装好了Golang运行环境
2. 进入`go-yedis`目录并运行`go build -o yedis-cli .\yedis-cli.go`构建客户端可执行程序
3. 进入`go-yedis`目录并运行`go build -o yedis-server .\yedis-server.go`构建服务端可执行程序
4. 直接像Redis一样执行`yedis-server`和`yedis-cli`便能够直接使用

#### Docker构建与运行
xx todo

##### 架构图
![架构图](https://oss.whoiszxl.com/Go-Yedis.png)

##### 公众号
![公众号](https://oss.whoiszxl.com/qrcode_for_whoisc137_258.jpg)

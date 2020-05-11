# 造轮子手册
[![MyWebSite](https://img.shields.io/badge/我的站点-whoiszxl-blue.svg)](https://whoiszxl.github.io)
[![docs](https://img.shields.io/badge/docs-reference-green.svg)](https://whoiszxl.github.io)
[![teach](https://img.shields.io/badge/教程-BohemianRhapsody-orange.svg)](https://github.com/whoiszxl/Monica)
[![email](https://img.shields.io/badge/email-whoiszxl@gmail.com-red.svg)](whoiszxl@gmail.com)


### 手册进度
▓▓░░░░░░░░░░░░░░░░░░░ 8% Go-Yedis开发中，已完成最基础的string，hash，list，set指令，简单的aof持久化，控制台交互。

### 手册介绍
记录当前学习造轮子的笔记手册，用简单的代码实现一些经典的软件框架等，如Redis,MySQL，Rpc，MyBatis，Spring等。

### 手册地址
文档地址：http://monica.whoiszxl.com/

### 手册目录
#### 1. go-yedis
Go-Yedis通过Golang语言实现一个简略版的Redis单机键值对数据库，其中实现了诸如SDS、LinkedList等基础数据结构，string、list、set、zset和hash五种对象类型，订阅发布、密码校验、慢查询日志等功能，RDB和AOF的数据层持久化等功能。因为Redis是使用C语言开发的，像分配SDS简单动态字符串内存释放其内存等操作Yedis都简化成使用Golang原生string，但是有些功能处不会做省略，旨在学习原生Redis的结构，比如sds的len属性，其设计是用来减少strlen函数的内存消耗，golang的string自带len属性，完全可以复用其属性，此处便不省略，为了多一点Redis的汁味。

考虑在代码开发完成后写一份图文手册，从0-1构建Yedis源码，附带Redis3.0的源码分析与不同之处，从而达到更好地理解Redis的设计思路，也能够更好地在工作中使用Redis。


#### 2. go-wysql
#### 3. go-xrpc
#### 4. java-zyBatis
#### 5. java-summer
#### 6. xx-ui

##### 公众号
![公众号](https://oss.whoiszxl.com/qrcode_for_whoisc137_258.jpg)

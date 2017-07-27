# somechat
A simple about webchat

## 数据模型
* 用户：email，name，password
* 通讯录：owner，userId， userName， userEmail
* 通讯录添加请求：from， to
* 私信数据：room（当前是双方id的md5值，为后续多人模式做铺垫），user（消息发送者），content，index（排序用，目前用全局递增值，改进为一个room一个递增值，防止createtime下不固定顺序），createtime

## 程序设计
* 采用iris作为http框架，理由是它功能足够丰富，且性能不错。
* web 只提供显示页面和基本请求（页面是刷新方式，api是ajax方式），web socket 负责私信和通知（未读信息、通讯录添加请求）
* 全局time使用utc，数据库存时间带zone，网页使用moment.js来处理utc和local。
* 使用cobra作为命令辅助
* 数据存储采用postgres，redis作为内存缓存（当前未加入）

### 目录结构
* cmd：cobra的cmd
* core：核心公用组件（数据）
* log：日志（当前为标准log，加上结构化内容输出，改进方案是可自定义输出源）
* server：暴露出的服务，web app和chat app，各自服务不共享代码，代码自治。

### chat app 组件
> agent

一个用户的web socket conn代理，id为用户id。
功能是发消息。

> agent bus

agent的总线，注册agent，推出agent，发现agent。
当前只有local agent bus，只能发现一个服务内的agent。
其实是提供distributed agent bus功能，使用consul或etcd实现分布式的总线。也是改进方案。

> agent gateway

功能为向某个agent发信息，自身是绑定一个bus和一个message store。
当agent在bus上，往agent的ws conn 中发完信息后，通过message store保存。
当agent不在bus上，直接message store保存信息，同时保存一个未读信息。

> message store

数据存储的代理。
功能为保存通信内容，未读信息。

> message

分 head和body两个部分组成
head：一个KV，通常有kind（chat|notify）name，，room，fromUser，datetime等。
body：内容本身。字符串。


### 逻辑设计
> 登录注册

简陋的方式，客户端维护cookie，保持服务无状态

> 通讯录

* 打开页面时，web直接输出联系人列表，每个联系人的未读信息。
* 联系人添加请求：先创建一个添加请求，再通过该页面已向chat app 建立的ws conn 发 通知。
* 联系人添加请求通知：通过ws conn监听kind 为 notify，name为contact add req的message，通过jquery进行页面绘制。
* 联系人添加请求接收：API的方式往自己的通讯录里和对方的通讯录里加联系人，删除请求通知，jquery本地加联系人节点，通过ws向chat app 发送一个接收消息，如果对方在线，则收到消息后，自动在页面上加联系人。
* chat信息数通知，对方在线向我发信息，本页面的ws conn 监听 该notify，name为 message unread，然后增加页面上对应联系人的未读信息数。
* chat对话：打开一个新的页面，新建ws conn来chat。并调用api删除该联系人的未读信息，成功后删页面上的。
* 陌生人chat：向未在通讯录中的联系人发信息，打开一个chat页面，调用 web app的api 建立一个room（chat history），并向对方发起添加请求，在得道正确的返回结果后，通过ws conn发对应notify。在对方未同意前，是无法看到所发信息，只有同意后，再点开联系人后才能看到。但整个room的内容一直在，不与添加请求相关。

> chat

一个room一个页面，一个页面一个ws conn。
只是通过 ws conn，走chat app 服务，注册一个agent，然后进行对话。
具体过程见chat app 组件。

## 改进方案

* 实现distributed agent bus。
* log自定义输出源。
* redis 缓存加入。
* Chat app 和 web app 不再通过core共享一些公用组件，采用eda进行数据交互，满足分布式微服务的数据自治。

# Example

somechat web --conf={config_file_path}

somechat chat --conf={config_file_path}

config.example  
```
postgres:
  url: host=127.0.0.1 port=15433 dbname=pharosnet user=pharosnet password='pharosnet@db'
    connect_timeout=60 sslmode=disable
  maxIdle: 4
  maxOpen: 128
redis:
  addr: 
  password: 
  db: 0
web:
  port: :8080
  static: /Users/doaman/workspace/projects/liulishuo/res/somechat-web/static
  tpl: /Users/doaman/workspace/projects/liulishuo/res/somechat-web/tpl
  favicon: 
```

 
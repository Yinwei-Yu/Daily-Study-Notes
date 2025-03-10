
[[计算机网络：自顶向下方法（原书第8版）.pdf]]

# 介绍

- 计算机网络和互联网
- 应用层
- 传输层
- 网络层：数据平面
- 网络层：控制平面
- 数据链路层和局域网

不出错 不重复 不丢失 不失序
可靠要求高tcp，速度慢
不高用udp

---

网络层
传统：IP，路由-查路由表转发
现代sdn：
数据平面-交换机
控制平面-网络操作系统

- 网络安全
- 无线和移动网络
- 多媒体网络
- 网络管理

# 概述
## 什么是internet

网络
计算机网络
互联网

---
计算机网络

>节点和边

主机节点,数据源,数据目标
数据交换节点,路由器,交换机,中继器,转发数据,既不是源也不是目标.

>链路

链接主机,交换节点
接入网链路:主机连接到互联网
主干链路:路由器之间的链路

>协议

标准
按照层次不同
每层的协议可以分为不同种

---

互联网
由tcp/ip协议为主的一簇协议为主的网络

---

主机=端系统 end sysytem / host
运行网络应用程序

通信电路
光纤,同轴电缆,无线电,卫星
传输速率 = ==带宽==(bps)

---

> 协议

对等层的实体,运行的模块在通信过程中因该遵循的规则的集合

==报文格式,次序,动作==
 
应用层协议
网络层协议
链路层
物理层

遵守相同协议的设备才能够通信

互联网中所有通信行为都受到协议规范

---

互联网是网络的网络

互联网标准由RFC文档制定,IETF组织

---

> 从服务的角度看互联网

从应用进程角度看:分布式应用进程,基础设施向应用进程提供服务

基础设施:应用层之下所有的应用实体

分布式应用是网络存在的理由:web,分布式游戏,电子商务,社交网络等

网络为应用进程提供socket套接字接口

---

==网络结构===：
网络边缘，网络核心，接入网、物理媒体
## 网络边缘

>边缘系统

主机，应用程序（客户端和服务器）
基础设施

通信模式:
- cs模式:客户端/服务器模式
客户端向服务器请求接受服务
服务器占主要地位
客户端后运行
服务器农场
负载大,若服务器宕机则损失太大

- p2p模式 peer-peer 对等模式
很少(或者没有)服务器
每个节点既是客户端又是服务器
例如迅雷等

采用网络设施的面向连接服务:
在端系统之间传输数据
先握手(在通信主机间建立连接状态)
底层协议栈做好准备
tcp面向连接服务:
1. 可靠,保序
2. 流量控制,控制流量发送速度
3. 拥塞控制
应用:http,ftp(文件传送),远程登陆,邮件

udp:无连接服务
速度快,但是不可靠
应用:流媒体,远程会议,DNS,Interney电话

## 网络核心

>网络核心

起到数据交换作用

由很多数据交换节点形成的分布式系统
路由器交换机

数据怎么通过网络进行传输？

 1. 电路交换

通过信令系统为两个端建立单独的线路(piece)
独享资源,不同享->性能保障
没有数据发送,被分配的资源会浪费
传统电话网络采用
需要建立呼叫连接

网络资源(如带宽)被分为片(piece)
- 频分-FDM
- 时分-TDM
- 波分-WDM

>FDM
>交换节点之间带宽较宽,将有效频率范围划分为多片,多路复用

>TDM
>节点之间分解为T的时间周期(时隙slot),每个周期被分为多个片,每个片供一个用户使用

>WDM
>光纤通信,可用波段可分为若干小波段,每个用户使用一个小波段

因为:
- 建立时间长(秒级)
- 计算机之间的通信有突发性,浪费的片较多
- 可靠性不高

 2. 分组交换

>以==分组==为单位存储-转发方式
>网络带宽资源不再分为片,传输时使用全部带宽
>主机间传输的数据被分为一个个单位,称为分组(packet)
>在每个交换节点之间采用以分组为单位,存储转发的方式

转发前节点需收到整个分组
延迟比线路交换要长->共享性
排队时间 

速率为R bits的链路传输L bits的数据,需L/R的时间
==发送和接受是一体两面==

>排队和延迟

如果到达速率 > 链路的输出速率
- 分组会排队等待传输
- 若路由器的缓存用完了,则分组会被抛弃

>网络核心的关键功能

1. 转发
将分组从路由器的输入链路转移到输出链路

2. 路由
决定分组采用的源到目标的路径->路由算法

>统计多路复用
>随机划分时间片

- 适合突发式数据传输
	- 数据共享
	- 简单,不用建立呼叫
- 过度使用会造成网络拥塞
	- 分组延时和丢失
	- 对可靠地数据传输需要协议来维护
- 怎样实现类似电路连接的功能
	- 保证音视频应用需要的带宽
	- 一个尚未解决的问题

>分类-按照有无网络连接

1. 数据报网络datagram
- 分组的目标地址决定下一跳
- 分组携带了目标的完整地址
- 通信之前不需要建立连接
- 每个分组传输是独立的->可能走不同路径(路由表会变),可能会失序
-  不需维护连接状态

2. 虚电路网络->靠信令传递
- 通信前握手,建立虚拟线路
- 分组携带虚电路标识(VC ID),不携带完整地址
![[Pasted image 20240911221442.png]]

## 接入网,物理媒体

> 接入网 物理媒体

有线或者无线通信链路

住宅接入
单位接入
无线接入

- 接入带宽
- 共享还是独享

## 住宅接入 modem 调制解调器

将上网数据==调制==加载音频信号上，在电话线上传输，在局端将其中的数据==解调==出来

猫：modem

带宽通常低：56kbps之下

不能同时上网和打电话

## 接入网 DSL

语言数据在专享线路的不同频段传输
仍是调制解调
只不过在4khz之上
不对称分布，上行下行

## 接入网：线缆网络

有线电视信号线缆双向改造

FDM：在不同频段传输不同信道的数据，数字电视和上网数据

用户共享到线缆头端的接入网络

HFC：非对称传输速度，下行最高30Mbps，上行最高2Mbps

## 企业接入网络

交换机级联

顶层交换机通过路由器接入到各种网络

## 无线接入网络

> 广域无线接入

基站

> 无线LANs:

建筑物内部

## 物理媒体

第零层

媒体:发送和接受两个节点,在节点之间传送bit的介质

1. 导引型媒体

信号在导体内部

- 双绞线:两根绝缘铜导线拧合
- 同轴电缆;两根同轴铜导线
- 光纤:每个脉冲表示一个bit,在玻璃纤维中传播,误码率低,不受电磁波信号影响
	- 单模光纤->传输距离远
	- 多模光纤

2. 非导引型媒体

开放空间传播
双向
无需线缆
影响:
	反射
	吸收
	干扰

- 微波
- LAN
- 广域网
- 卫星
	- 延迟高
	- 同步静止卫星和低轨卫星
	- 每个信道Kbps到45Mkps
## 互联网和ISP

端系统通过接入ISP访问互联网

接入的ISPs必须是互联的

- 互联--不可行，不可扩展
- 每个接入ISP都连接到全局ISP(global ISP)

ISP竞争->不同的运营商
合作->不同global ISP对等连接(peering link),IXP,不涉及费用结算

业务细分->全球接入和区域接入

ICP(互联网服务提供商)会部署机房(DC,data center)->服务,价格

1. 全球/全国覆盖:点少,速率高,直接与其他一层ISP相连,与大量第二层ISP连接
>POP
>上层ISP和下层ISP通过POP相连
2. 区域性ISP
3. 第三层ISP和localI SP

## 分组延时 丢失 吞吐量

### 分组丢失和延时的发生

分组到达链路的速率超出了链路的输出能力
分组等待排队头才传输

若队列已满,则分组被丢弃
太长的队列没有意义

### 四种分组延时

1. 节点处理延时

检查bit级差错
检查到向何处

2. 排队延时

取决于路由器的拥塞程度

3. 传输延时

R=链路带宽
L=分组的长度
时间=L/R

4. 传播延时
d=物理链路的长度
s=媒体上的传播速度
传播延时=d/s

两个节点之间的传输称为一个hop(段,跳)

节点延时=处理延时(微秒)+排队延时(随机)+传输延时(微秒到毫秒)+传播延时(几微秒到几百毫秒)

### 排队延时

取决于流量强度I

I=La/R

R为链路带宽,L为分组长度,a为单位时间内希望通过链路转发的分组数量

流量强度在0-1之间

I越接近1,排队延时越接近无穷大 

linux traceroute
windows tracert
测试两个地址间的传输延时

利用了ICMP协议

Round trip time 往返延迟 RTT

### 分组丢失

队列缓冲区有限

分组到达满的队列时会丢失

丢失的分组可能会被前一个节点或者源主机重传,UDP不重传

### 吞吐量

数据量/单位时间

瞬间吞吐量
平均吞吐量

吞吐量=min{Rs,Rc}

瓶颈链路:端到端路径上限制端到端吞吐的链路

==实际情况==
一个链路可能被其他的通信对使用,则带宽需要除以通信对数量
瓶颈带宽相应减小

## 协议层次和服务模型

计算机网络是一个复杂的系统

层次化方式实现复杂网络功能

1. 将复杂功能分解为明确的层次，每层有功能，每个模块向上层提供接口服务
2. 本层协议实体相互交互执行本层的协议动作
3. 实现本层协议利用下层提供的服务

## 服务和服务访问点

服务用户:tcp实体上有很多应用
tcp向应用提供服务,应用即用户

服务访问点:
	SAP
	层间接口处有一些点来区分不同的信息提供给上层用户
	下来标注,上去区分

原语:上层使用下层服务的形式,例如api函数

面向连接服务:进程通信之前要先握手

## 服务的类型

无连接的服务
	UDP
面向连接的服务
	TCP

## 服务与协议的区别

服务是垂直
协议水平

服务在系统内部相邻两层之间，服务用户使用下层提供的元语服务

关系：
	本层协议的实现需要下层提供的服务
	实现协议的目的是向上层提供更好的服务

## 数据单元 DU

上层要求下层传递的信息为SDU,接口控制信息为ICI,
两者在一起称为PDU,协议数据单元

SDU过大可能被分解

SDU过大可能被合并

## 分层处理和实现复杂系统的好处

结构化:模块化易于升级和维护,改变某一层的服务不会改变其他层的服务,易于采用新的技术

概念化:结构清晰

>问题
>子系统需要交换信息,效率较低

## internet协议栈

==物理层==:

将上层的帧的每个bit转化为物理信号发送
接收物理信号并转化为数字信号

==链路层==:

在物理层提供服务基础上在相邻两点之间传输以帧为单位的数据
PPP,WIFI

==网络层==:

传输以分组为单位的从源主机到目标主机的端到端的服务
不可靠
IP,路由协议

==传输层==:

1. 进程到进程的区分
2. 将网络层提供的不可靠的服务转化为可靠的服务
TCP,UDP

==应用层==:网络应用
	提供网络应用服务
	FTP,SMTP,HTTP,DNS

## ISO/OSI参考模型

应用层

==表示层==
==会话层==

传输层
网络层
链路层
物理层

表示层:允许应用解释传输的数据,加密,压缩,机器相关转换

会话层:数据交换的同步,检查点,恢复

由应用实现

## 各层次的协议数据单元

应用层:报文message

传输层:报文段(segment)

网络层:分组packet(无连接方式:数据报)

链路层:帧

物理层:位

# 应用层

## 概述

原理
应用及其协议
DNS 域名解析
TCP 编程
UDP 编程

## 应用层原理

### 体系结构

> 客户-服务器模式(c/s)

- 服务器
	- 一直运行
	- 固定的ip地址和周知的端口上
- 客户端
	- 请求服务器资源
	- 运行在动态ip地址上
	- 与互联网有间歇性连接
	- 不与其他客户端通信

服务器是中心,资源在服务器
服务器随着用户增长,超过一个阈值之后,性能急剧下降
可扩展性差
可靠性差

> p2p

任意端系统之间可以通信

每个节点既是服务器也是客户端

缺点是难以管理

> 混合体

==napster==

文件搜索集中

文件分发p2p

==即时通信==

### 进程通信

> 进程

在主机上运行的应用程序

在不同端系统上发送报文message

客户端进程:发起通信
服务器进程:等待连接的进程
==p2p也有客户端和服务器之分==

> 需要解决的问题

1. 进程标识和寻址问题
2. 传输层-应用层服务如何
	位置:层间界面SAP(TCP/IP:socket)
	形式:应用程序接口API(TCP/IP:socket api)
3. 定义报文的格式,时序和相应解释,协议

==问题的解决==:

> 标识进程

1. 在哪个终端设备上-主机有唯一的32位ip地址
2. TCP还是UDP
3. 在哪个端口(port)上

一个进程:IP+port标示 端节点
本质上一堆主机进程之间的通信由2个端节点构成

> 传输层提供的服务

层间接口携带的信息:
	1. message
	2. 发送端的应用进程的端口号,ip
	3. 发给谁

传输层实体根据以上信息对以上信息进行封装
段头和数据部分

每次都传输三种信息比较繁琐

==socket==:代号标示通信的双方或单方
操作系统管理的一个整数,代表双方的ip和端口号
tcp四元组
udp二元组
==本地==标识,对方不知道

TCP
穿过层间接口的信息比较少
本地操作系统维护一张socket表,由socket值映射到源ip,port,对方ip,port

UDP
每个报文都是独立传输的
前后报文可能给不同的分布式进程
只用一个整数表示本地ip.port,不代表一个会话关系
UDP socket:本ip,本端口->一个端节点 
但是在传输报文时必须提供对方的ip,port
即:向传输层提供:数据,socket,对方的ip和port

> 如何使用传输层提供的服务实现应用

应用层协议:定义了运行在不同端系统上的应用进程如何交互报文
- 格式
- 语法
- 字段的语义
- 发送报文和对报文相应的规则

与网络层协议相关的运行中的模块称为实体

应用协议只是应用的一部分

公开协议:HTTP
私有协议:Skype

传输层提供的服务的评测标准:

==数据丢失率==

==吞吐量==

==延时==

==安全性==
	机密
	完整性
	可认证性

 > 安全tcp
 
 SSL:安全套接字层 secure socket layer(应用层)
 1. 在tcp上实现,提供加密的tcp连接
 2. 私密性
 3. 数据完整性
 4. 端到端鉴别
https=http+ssl

ssl socket api
- 应用通过api将明文交给socket,ssl将其加密后在互联网上传输

## Web与HTTP

### 一些术语

- Web:网络应用
- Web页：一些==对象==组成
- 对象可以是HTML文件，JPEG图像，Java小程序等
- Web页有基本HTML文件，这个文件中又包含若干个对象的引用（连接）
- 通过==URL==引用对象
- URL:==访问协议,用户名,口令字,端口等==
- URL格式:协议名 用户:口令 主机名 路径名 端口
- 支持匿名访问,即不提供用户名
- HTTP端口默认为80号端口
- www -> world wide web

### HTTP概况

==超文本传输协议==

- Web的应用层协议
- 客户/服务器模式
	- 客户请求
	- 服务器相应,发送对象的Web服务器
- 使用TCP协议
	- 客户发起一个与服务器的==TCP连接==(建立套接字,端口号为80)
	- 服务器接收客户的TCP连接
	- 在浏览器与Web服务器交换HTTP报文
		- 如果是HTML文件就画出来,其中的链接进行解析并与它的URL进行通信,拉过来画出来
	- 浏览器从Web服务器得到文件后关闭tcp连接
	- HTTP协议是==无状态==的
		- 简单,不需维护状态
		- 支持更多客户端

### HTTP连接

> 非持久HTTP HTTP/1.0

- 一个连接只能传输一个对象
- 传输对象之后连接关闭
- 传输多个对象需要多个连接

![[Pasted image 20241011154449.png]]
> 持久HTTP HTTP/1.1

- 一次往返建立连接
- 客户端接收到对象后连接不关闭
- 再有对象请求在此连接上发送

> 响应时间 往返时间(RTT-round-trip-time)

一个小的分组从客户端到服务器再回到客户端的时间(传输时间可以忽略)

响应时间:
- 一个RTT用于建立TCP连接
- 一个RTT用于HTTP请求和等待HTTP相应
- 文件传输时间
==共两个RTT时间==

> 非持久的缺点

需两个RTT,时间长

> 持久 | 非流水线(pipeline)

一次请求和返回的对象只有一个,在上一个对象接收之后再请求下一个对象

> 持久 | 流水线

在上一个对象还没有接收到之前就发送这一个对象的请求
HTTP/1.1的默认方式

### HTTP请求报文

- 分为:请求,响应

#### 请求报文 ascll码形式

![[Pasted image 20241011155712.png]]
![[Pasted image 20241011155856.png]]
方法:GET,POST(上传),HEAD(只要HTML的头)
HOST: 域名
User-agent: 用户代理
post有实体部分

> 提交表单

1. POST方式
- 网页包括表单输入
- 包含在实体体中的输入提交到服务器
2. URL方式
- 方法:GET
- 输入通过请求行的URL字段上载,URL中提供参数信息

> 方法类型

HTTP/1.0
- GET
- POST
- HEAD
HTTP/1.1
- 上述三个
- PUT
- DELETE

#### 响应报文

![[Pasted image 20241011160834.png]]

![[Pasted image 20241011160847.png]]
状态码:
![[Pasted image 20241011161636.png]]

### 用户-服务器状态:cookies

HTTP是无状态的协议
服务器不维护客户端的状态
- 好处:简单,支持的客户端多

因为有的服务器需要保存用户的状态,比如购物网站

cookie有四个组成部分:
先看原理和过程:
1. 客户端第一次与服务器建立连接后,服务器给客户端分配一个cookie,HTTP报文中有一个cookie首部行,同时服务器端存一个cookie,客户端接收到的cookie保存到本地
2. 隔段时间再次访问时客户端带着cookie访问,服务器通过cookie查找用户状态

组成部分:
1. HTTP响应报文中有一个cookie首部行
2. HTTP请求报文有一个cookie首部行
3. 客户端系统保留一个cookie文件,由用户浏览器管理
4. Web站点有一个后端数据库

### Web 缓存(代理服务器)

==目标==:不访问原始服务器就满足客户的请求

两种访问服务器的方式:
1. 访问服务器

2. 访问服务器的代理服务器

浏览器将HTTP请求发给缓存
- 在缓存中直接返回该对象
- 不在缓存中,缓存再向服务器请求
==好处==
- 快
- 服务器端压力减轻

> 条件GET方法

访问的对象可能在服务器端更新,而在缓存没有更新

使用条件获取
在请求报文头部加上if get,若在last modified后服务器改变了,服务器返回新对象,否则只返回头部字段
304 Not Modified

## FTP

文件传输协议

- 向远程主机上传输文件或从远程主机上接收文件
- 客户端将文件传送到FTP服务器，从服务器下载文件
- FTP应用包括FTP客户端和用户接口，本地文件接口
- RFC 959
- 端口：21

### 控制连接与数据连接分开

- 使用tcp传输协议
- 客户端通过控制连接获得身份确认
- 客户端通过==控制连接==发送命令浏览元辰目录
- 服务器收到文件传输命令后，服务器打开一个客户端20号端口的==数据连接==（服务器主动）
- 文件传输完成后关闭服务器连接
- 服务器打开第二个tcp数据连接用于传输另一个文件
- 控制连接：带外（out of bound）->控制命令和数据传输在不同的连接上
- FTP服务器维护用户状态,当前路径,用户账户与控制连接对应

### 命令和相应

ascii码命令

```
USER username
PASS password
LIST
RETR filename
STOR filename

返回码
331 Uername OK,password required
...
```

## Email

==三个主要组成部分==
1. 用户代理
2. 邮件服务器
3. 简单邮件传输协议:SMTP

==用户代理==
1. 电子邮件软件
2. 撰写,编辑,转发邮件

### 邮件服务器

- ==邮箱==中管理和维护发送给用户的邮件
- 输出报文队列保持待发送邮件报文
- 邮件服务器之间的SMTP协议:发送EMAIL报文

用户代理发给邮件服务器,邮件服务器放在队列中发送给接收方的邮件服务器中,接收方邮件服务器放到用户邮箱中,接收方使用用户代理从用户邮箱中拉取邮件

### SMTP

- 服务器25号端口
- TCP
- 直接传输
- 三个阶段
	- 握手
	- 传输报文
	- 关闭
- 命令/响应 都是ascii码
	- 命令 ASCII码
	- 响应包括状态和响应信息
	- 传输字节都以ascii码形式进行

eg.

![[Pasted image 20241015204633.png]]
1. 持久连接
2. 报文要求为7位ascii码
3. 服务器使用CRLF.CRLF决定报文尾部
4. 客户端对服务器的请求为push
5. HTTP每个响应报文为单独的报文,SMTP将所有对象包括在一个message中

### 邮件报文格式

- 首部行
	- To
	- From
	- Subject
	- ==这些与SMTP指令不同==
- 主体
### 传输中文和其他文件怎么办?

使用MIME:多媒体邮件扩展
在报文首部用额外的行申明MIME内容类型

Base64编码,将若干个不在ASCII码范围内的字节映射到ASCII码范围内
接收方使用Base64解码并用MIME解释

![[Pasted image 20241015205604.png]]

### 邮件访问协议

SMTP传送到接收方邮件服务器
邮件访问协议:接收方从邮件服务器获得邮件
- POP3: 邮局访问协议
- IMAP: Internet邮件访问协议
	- 更多特性:远程目录维护
	- 在服务器上处理存储的报文
- HTTP
	- 方便

>pop3

- 用户确认阶段
	- 客户端用户名和密码
	- 服务器响应
- 事物处理阶段
	- list:报文号列表
	- retr:根据报文号检索报文
	- dele:删除
	- quit
- 下载删除模式/下载保留模式
- 具体某个邮箱无状态的
- 本地管理文件夹

>IMAP

- 每个报文和一个文件夹联系
- 允许用户用目录组织报文
- 允许用户读取报文组件
- 在会话中保留用户状态
- 远程管理文件夹

## DNS

 DNS是给应用用的应用

### 必要性

- IP地址标识主机和路由器
- IP地址不好记忆
- 人类倾向于使用有意义的字符串表示设备（域名）
- 存在字符串-IP地址转换的必要性
- 由DNS将字符串转化为二进制的网络地址

### 需要解决的问题

1. 如何命名设备？
	1. 有意义的字符串
	2. 层次化命名
2. 如何转化IP地址和命名
	1. 分布式数据库维护和相应名字查询
3. 如何维护？
	1. 增加或者删除一个域名，需要做什么工作

> 主要思路

1. 分层基于域的命名机制
2. 分布式数据库完成转换
3. 运行在UDP 53 号端口
4. 是核心的互联网功能,以应用层协议实现,位于网络边缘

> DNS的主要目的

1. 实现主机名和IP地址的转化
2. 主机别名到规范名字的转换
3. 邮件服务器别名到邮件服务器的正规名字的转换
4. 负载均衡(load distribution)

### 命名

- 一个层面的命名设备会有很多重名
- DNS采用树状结构命名
- Internet根划分为几百个顶级域
	- 通用:.com,.edu,.gov
	- 国家:.cn,.us,.nl,.jp
- 每个域下面可划分为若干个子域
- 树叶上主机

> 根名字服务器

共有13个根名字服务器
根名字服务器和顶级域名不同

> 域名管理

一个域管理其下的子域

> 域与物理网络无关

域遵从组织界限,不是物理网络
一个域的主机可以不在一个网络
一个网络的主机不一定在一个域
域的划分是逻辑的

### 域名到IP地址的转换

一个名字服务器的问题:
1. 可靠性:单点故障
2. 扩展性:通信容量
3. 维护问题:远距离集中式数据库

区域:
1. 区域的划分由区域管理者自己决定
2. 将DNS名字空间划分为互不相交的区域,每个区域都是树的一部分
3. 名字服务器:
	1. 每个区有有一个名字服务器,维护其下权威信息
	2. 名字服务器允许放置在区域之外,保障可靠性

TLD服务器:
负责顶级域名和所有国家级的顶级域名

区域名字服务器维护的资源记录:
1. 资源记录
	1. 作用:维护域名-IP地址(其他)的映射关系
	2. 位置:Name Server的分布式数据库中
2. RR格式
	1. Domain_name:域名
	2. Ttl:time to live:生存时间,长期/短期
	3. Class类别,互联网内的值为IN
	4. Value值:数字,域名,ascii串
	5. Type:类别:资源记录的类型(别名转换等)
		1. A:name为主机,value为ip
		2. CNAME:name为别名,value为规范名字
		3. MX:name为别名,value为name对应的邮件服务器名字
		4. NS:name为域名(子域),value为域名的权威服务器的域名(子域的名字服务器的名字)

### 工作过程

应用调用解析器(resolver)
解析器作为客户向Name Server发送查询报文(封装到UDP中)
Name Server返回响应报文(name/ip)
![](附件/Pasted%20image%2020241104155432.png)

> 本地名字服务器(local name server)

- 并不严格属于层次结构
- 每个ISP有一个本地DNS服务器(也称默认名字服务器)
- 当一个主机发起一个DNS查询时,查询被送到本地的DNS服务器(起代理作用)

> 名字服务器解析过程

1. 名字在local name server中
	1. 在区域内
	2. 缓存
2. 不在
	1. 联系根名字服务器顺着根-TLD一直找到权威名字服务器
	2. 递归查询:名字解析负担都放在当前联络的名字的服务器上![](附件/Pasted%20image%2020241104160324.png)
	3. 迭代查询:各级域名服务器返回的不是查询结果,而是下一个NS的地址,最终由权威名字服务器给出解析结果![](附件/Pasted%20image%2020241104160436.png)

### DNS协议,报文

查询和响应报文格式相同
![](附件/Pasted%20image%2020241104160716.png)
![](附件/Pasted%20image%2020241104160605.png)

ID号实现流水线工作,使得一个Name Server可以相应多个查询

### 提高性能:缓存

- 一旦一个服务器学到了一个映射,就将这个映射缓存下来
- 默认TTL为两天

### 维护问题:新增一个域

- 在上级域的名字服务器中新增两条记录,指向这个新增子域的域名和域名服务器的地址
- 在新增子域的名字服务器上运行名字服务器,负责本域的名字解析:名字->IP地址
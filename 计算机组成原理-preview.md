[[计算机组成与设计：硬件软件接口 RISC-V版.pdf]]

课程目标:单处理器计算机系统中各部分的内部工作原理，组成结构以及相互连接方式，具有完整的计算机系统的整机概念

周四下午两点半-五点半 计院c区 office hour A515

# 绪论

![[Chapter_01.pdf]]

摩尔定律:每18月芯片集成度翻番
后摩尔时代

> 分类

- 个人pc

通用处理器
==性价比高== -tradeoff
1. intel-x86
amd

2. arm架构
苹果,华为

3. risc-v架构

4. mips架构
龙芯
loong Arch

- 服务器

基于网络
高容量,性能,可靠
服务范围广

- 超级计算机

高端科学与工程计算
高容量但是市值小

- 嵌入式计算机
主要是arm和risv
隐藏在硬件中
功耗,性能,成本要求高

> 后pc时代

- 个人移动设备PMD

- 云计算
数据仓库
软件及服务

> 性能

1. 算法
2. 语言,编译器
3. 指令架构
4. 处理器
5. 输入输出系统(包括os)

> 八大思想

[[Chapter_01.pdf#page=15&selection=0,17,0,17|Chapter_01, 页面 15]]

1. 为摩尔定律设计
2. 使用抽象简化设计，忽略底层环境，内存大小等
3. 使通用操作更快
4. 并行
5. 流水线
6. 预测
7. 内存层次结构
8. 冗余提高可靠性

 > Below Your Program

[[Chapter_01.pdf#page=16&selection=10,0,10,18|Chapter_01, 页面 16]]

汇编语言-助记符

> 冯诺依曼结构

[[Chapter_01.pdf#page=18&selection=10,0,10,24|Chapter_01, 页面 18]]

特点:
1. 二进制
2. 五大部分:运算器,控制器)->中央处理器,存储器,输入设备,输出设备
3. 存储程序原理->快速,只读

> 显示器

像素,位深度
光栅扫描

图像分辨率\*位深度 存储在帧缓存中(显卡中),即显存

独立显卡有自己的显存
集成显卡用内存作为显存

NPU 神经网络计算单元
TPU -> 谷歌

> Abstractions

[[Chapter_01.pdf#page=24&selection=10,0,10,12|Chapter_01, 页面 24]]

指令集架构 ISA
硬件软件接口

> 半导体科技

[[Chapter_01.pdf#page=28&selection=0,0,0,24|Chapter_01, 页面 28]]

## 性能

 > 定义性能
 
 [[Chapter_01.pdf#page=32&selection=10,0,10,20|Chapter_01, 页面 32]]

 1. 响应时间

处理时间,等待时间(空闲时间),IO时间,系统调用时间

性能反应->processing,os overhead->cpu时间

 2. 吞吐量

单位时间完成的任务量
针对多处理器

==cpu时间==
处理任务所需时间,不包括io和其他任务share的

> cpu工作

[[Chapter_01.pdf#page=36&selection=10,0,10,12|Chapter_01, 页面 36]]

周期
[[Chapter_01.pdf#page=37&selection=10,0,10,8|Chapter_01, 页面 37]]

cpu时间=cpu周期数\*周期时间=周期数/cpu频率

每条机器指令需要的时间和cpu周期数都是固定的

cpu周期数=求和(每条指令执行周期数)

提升性能
[[Chapter_01.pdf#page=37&selection=14,0,14,23|Chapter_01, 页面 37]]

==权衡周期数和频率==

**CPI** 平均每条机器指令所需周期数
**ISA** 指令集架构

cpu周期数=CPI * 指令数

eg1.
[[Chapter_01.pdf#page=38&selection=10,0,10,16|Chapter_01, 页面 38]]
==cpu主频提高和周期数增加不成正比==

eg2.
[[Chapter_01.pdf#page=40&selection=0,9,10,11|Chapter_01, 页面 40]]
同样指令集架构的CPI也不一定相同
使用比值来表示性能强弱

>CPI注

==CISC--精简指令集==

当指令集复杂度太高时,即方差太大时

1. 应用分层平均数

2. CPI加权平均数
[[Chapter_01.pdf#page=41&selection=10,0,10,18|Chapter_01, 页面 41]]

eg.[[Chapter_01.pdf#page=42&selection=10,0,10,11|Chapter_01, 页面 42]]

> cpu性能相关

算法--IC,CPI
语言--CPI,IC
编译器--IC,CPI
指令集架构--IC,CPI,T$_c$

>其他一些指标

[[Chapter_01.pdf#page=44&selection=0,0,0,5|Chapter_01, 页面 44]]

## 功耗

功耗->影响摩尔定律

> 功耗墙

- 无法再降低电压
- 无法散更多的热


# 指令:计算机的语言 ISA

![[Chapter_02-RISC-V.pdf]]

cpu设计基于指令集架构
## 指令架构

- cpu所有指令的集合
- 不同计算机有不同的指令集
	- 但是有很多共同点
- 早期计算机有很简单的指令集
- 许多现代计算机也有简单指令集架构

## 指令的构成

`|操作码(OPcode)  | 操作数(data)|`

操作数:
	源操作数    目的操作数

## 指令类型

- 运算指令
	- 算术运算
	- 逻辑运算--与或非,移位运算

- 数据传输
	- load-从内存读取
	- store-保存到内存
	- 可以读取外设

> 为什么risc-v的读取指令可以读取外设而intel则需要额外的IN OUT指令?

- 控制指令
	- 条件分支
	- 无条件分支

## 指令规整

固定长度
	读取指令位数固定,4个字节

>intel怎么实现变长指令读取?

操作码长度固定为7位
	低7位,送至译码器

操作数长度25位->仅三种
	寄存器
		32个32位寄存器(或64位)
		寄存器有多少位且运算器能完成运算的位数称为机器的位数
		32位为一个字,64位为双字
		编号x0-x31
		5位寄存器号->识别寄存器
		==x0:const value 0==->好处?
		简单且快速
	内存->memory
		内存地址
		数组,结构,联合,动态数据
		risc-v不要求对齐(访问的第一个内存位置可以被4整除)
		操作数长度不够存储地址->寄存器内容加上常数
		==空间大==
	立即数
		常规操作更快

## 指令的表示

[[Chapter_02-RISC-V.pdf#page=18&selection=0,55,2,25|Chapter_02-RISC-V, 页面 18]]

### R型指令->操作数都是寄存器

![[Pasted image 20240920151101.png]]
==运算指令==
==逻辑指令==

寄存器占了15位

指令扩展技术

0110011作为R型指令操作码

3,7位funct作为功能码

==速度快==

### I型指令->有一个操作数是立即数

[[Chapter_02-RISC-V.pdf#page=22&selection=2,0,2,28|Chapter_02-RISC-V, 页面 22]]

![[Pasted image 20240920152324.png]]

12位用于存放立即数

范围为 $-2^{11} -- 2^{11}-1$  

### S型指令 store指令

lw指令为I型指令

sw指令为S型指令

![[Pasted image 20240920153828.png]]

### 逻辑运算

[[Chapter_02-RISC-V.pdf#page=25&selection=2,0,2,18|Chapter_02-RISC-V, 页面 25]]

### 移位指令

[[Chapter_02-RISC-V.pdf#page=26&selection=2,0,2,15|Chapter_02-RISC-V, 页面 26]]

immed=5
fun7=7

移位最多32位
### 内存指令

`lw rd imm(rs1)`
	load 4 bytes at address at rs1+imm into rd
	
`sw rs2 imm(rs1)`
	store rs2 starting at address at rs1+imm

`lb sb` 读/写低8位
	放在低8位,高位符号扩展

`lh sh` 读/写16位

`lbu lhu`无符号扩展

### U型指令 lui

==32-bit constants==
	`lui rd,constant`
	20位立即数放在高20位,低12位为0

==auipc== rd imm
	imm<<12+pc->rd
`auipc x5,Label`
[[Chapter_02-RISC-V.pdf#page=87&selection=12,0,12,17|Chapter_02-RISC-V, 页面 87]]

[[riscv基本指令集.pdf#page=1&selection=91,1,92,6|riscv基本指令集, 页面 1]]

### 条件指令-SB型

`beq rs1,rs2,L1`

`bne rs1,rs2,L1`

L1为一个地址

因为指令低2位永远为00,但又由于支持16位指令,它的低位永远为0,所以不用存立即数的最低位,所以寻址范围 $2^{-12}$  ~  $2^{12}-1$

`beq x0,x0,L1` 可用作无条件跳转指令

### SB型指令

[[Chapter_02-RISC-V.pdf#page=52&selection=21,1,22,10|Chapter_02-RISC-V, 页面 52]]

11,12位的存储位置的设置是为了和S型指令保持一致

12位放在最高位保持了最高位为符号位的习惯

目标地址=pc+imm  *  2

作业:用汇编写一个冒泡排序

### 分支控制指令-过程实现

[[Chapter_02-RISC-V.pdf#page=33&selection=2,0,2,12|Chapter_02-RISC-V, 页面 33]]

[[Chapter_02-RISC-V.pdf#page=37&selection=2,0,2,27|Chapter_02-RISC-V, 页面 37]]

`jal x1,ProcedureLabel`

将返回地址PC+4给x1寄存器
跳转至`pc+ProcedureLabel

![[Pasted image 20240927143550.png]]
`jalr x0,0(x1)` ->i型指令

> 过程调用

[[Chapter_02-RISC-V.pdf#page=36&selection=0,55,2,17|Chapter_02-RISC-V, 页面 36]]

寄存器:
[[Chapter_02-RISC-V.pdf#page=7&selection=0,0,0,16|Chapter_02-RISC-V, 页面 7]]

1. 传参,x10 - x17,x10和x11保存返回值
2. 控制转移->jal
3. 保存调用者状态
4. 进行函数内指令
5. 返回值放到x10或x11寄存器
6. 控制转移给调用者->jalr

> 嵌入式过程调用->non-leaf过程

[[Chapter_02-RISC-V.pdf#page=42&selection=2,0,2,19|Chapter_02-RISC-V, 页面 42]]

- 栈保存调用者状态


## 字符串

[[Chapter_02-RISC-V.pdf#page=47&selection=2,0,2,14|Chapter_02-RISC-V, 页面 47]]

ascll 码
unicode 码-兼容ascll码
变长编码

字符按照字节存储
字符串占用连续的字节空间

### 按字节/半字/字 读写操作

[[Chapter_02-RISC-V.pdf#page=48&selection=2,0,2,29|Chapter_02-RISC-V, 页面 48]]

按照这些指令进行字符的

## 伪指令

```risc-v
mv rs1,rs2
j Label->jal x0,Label
jr x1->jalr x0,0(x1)
li x5 0xDEADBEEF
la x5,Label ->auipc
```
## 寻址模式
 
[Chapter_02-RISC-V.pdf#page=54&selection=0,0,0,25|Chapter_02-RISC-V,  页面 54]]

1. 立即数寻址
2. 寄存器寻址
3. 基址寻址
4. PC相对寻址

## 多任务

[[Chapter_02-RISC-V.pdf#page=56&selection=2,0,2,15|Chapter_02-RISC-V, 页面 56]]

原子操作指令(atomic)->无法被打断

> Synchronization in RISC-V

[[Chapter_02-RISC-V.pdf#page=57&selection=2,0,2,25|Chapter_02-RISC-V, 页面 57]]

`lr.w rd,(rs1)`
	load from address in rs1 to rd
	cpu中一个特殊的寄存器存放rs1中的地址,flag标志位置0
	此后所有访存操作的地址都需要和此地址比较,若相同,则flag置1

`sc.w rd,(rs1),rs2`
	先判断rs1地址是否被访问(用flag)
	store rs2 to address in rs1
	写成功rd置0,更改rs1地址值,否则rd置非零值,不写

> example

[[Chapter_02-RISC-V.pdf#page=58&selection=4,0,4,7|Chapter_02-RISC-V, 页面 58]]

  1表示锁,0表示解锁

# 计算机的算术运算

![[Chapter_03.pdf]]

## 算术运算

整数 ALU ->也能进行浮点计算,但是较慢
小数 协处理器

a-b=a+b的补码=a+~b+1

多媒体运算->矩阵运算->SSE7(单指令多数据流SIMD)->[饱和运算](https://zh.wikipedia.org/zh-cn/%E9%A5%B1%E5%92%8C%E8%BF%90%E7%AE%97)(最大数替换上溢,最小数替换下溢 )

## 乘法运算

| 部分积  | 乘数   | 被乘数  |
| ---- | ---- | ---- |
| 0000 | 1001 | 1000 |
| 1000 |      |      |
| 0100 | 100  |      |
| 0000 |      |      |
| 0010 | 10   |      |
| 0000 |      |      |
| 0001 |      |      |
| 1000 | 1    |      |
| 1001 |      |      |
移位,相加
用到的逻辑部件:四个全加器构成的加法器,移位器,存储部分积、乘数和被乘数的寄存器

上面是课上老师讲的，书上的内容是：
> 在每个流程中,先判断乘数的最低位是1还是0,如果是1,就把被乘数加到积寄存器中,如果是0,就不加
> 接着将被乘数(竖式里写在上面的)左移一位,乘数右移一位,因为在竖式中中间积左移一位,相当于变大了,又因为二进制乘法里面被乘数只有他自己和全0之分,所以直接将被乘数左移一位就相当于竖式里往左写一位.被乘数右移一位是为了看最低位是0还是1
> 重复以上步骤
![[Pasted image 20241009164156.png]]

eg.2 * 3 / 0010 * 0011:
![[Pasted image 20241009164236.png]]


[[Chapter_03.pdf#page=8&selection=10,0,10,20|Chapter_03, 页面 8]]
Product低位存乘数高位存部分积,乘数移位正好和部分积组合

乘法运算用基本指令集实现,需要位数轮循环

### 乘法器改进

[[Chapter_03.pdf#page=9&selection=0,0,0,17|Chapter_03, 页面 9]]
阵列乘法器

### 乘法指令

[[Chapter_03.pdf#page=11&selection=10,0,12,16|Chapter_03, 页面 11]]
```risc-v
mul # 获取低32位结果
mulh # 获取64位乘积中的高32位结果,两个操作数都是有符号
mulhu # 同上,但两个操作数都是无符号
mulhsu # 同上,但是一个是无符号一个是有符号
```


快速除法 SRT

[除法指令]([[Chapter_03.pdf#page=16&selection=10,0,12,10|Chapter_03, 页面 16]])
rem->取余数

## ALU


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

ascii 码
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

[[Chapter_03.pdf#page=17&selection=0,0,0,3|Chapter_03, 页面 17]]

详见附录A

数据选择通路

支持slt和bne指令：

Set信号判断是否小于

溢出检测模块四个输入,a,b,Binvert,Carryout,a+b

## 浮点数

[[IEEE-754]]

> 硬件设计

![[Pasted image 20241011093542.png]]

浮点计算硬件:
1. 加减乘除
2. 倒数(reciprocal)
3. 开方
4. 太慢了,即使用流水线也比整数运算器慢

>浮点数寄存器

32个浮点数寄存器
f0~f31
有64位,单精度浮点数存在低32位中
全部给浮点数用

## 浮点数运算指令

1. fadd.s fsub.s fmul.s fdiv.s fsqrt.s
2. fadd.d fsub.d fmul.d fdiv.d fsqrt.d
3. feq.s flt.s fle.s 等于 小于 小于等于
4. feq.d flt.d fle.d
	1. 比较结果0或1放在整数目的寄存器中
	2. 再使用beq和bne来跳转
5. flw fld 取字,取双字
6. fsw fsd
7. B.cond 条件分支指令

eg. 
[[Chapter_03.pdf#page=45&selection=10,0,16,1|摄氏度华氏度转化]]

# 第四章 处理器

![[Chapter_04.pdf]]

单周期处理器
流水线

支持RV32I,基本指令集

>单周期处理器

单发射

一个cpu周期内只执行一条指令

指令执行:
1. 取指
	1. pc地址发送给mem,mem将对应地址的指令发送给cpu
	2. pc+4 -> pc
2. 译码opcode func3 func7->由control进行
	读操作数->datapath进行
3. 执行指令

![[Pasted image 20241016104611.png]]

![[Pasted image 20241016104737.png]]

## 逻辑部件

1. 组合逻辑单元->执行指令
2. state(sequential) elements 状态单元(时序单元)->D触发器构成->保存信息

## 时序单元

D触发器
边缘触发器

组合单元的输入来源于状态单元->保证运算的准确性和稳定性
组合单元的输出输出到状态单元->保存运算结果

## 设计数据通路

需要有:寄存器,ALUs,MUXs,mem

1. 取指令
	1. pc->RF(寄存器),instruction memory(指令mem)
	2. pc+4->adder
2. 译码和读操作数
	1. control
	2. RF->读出寄存器中内容
		1. ![[Pasted image 20241016111941.png]]
		2. 上图解释:rs1,rs2,rd双端口读,单端口写;RegWrite控制写信号
	3. 执行
		1. R型指令:ALU,RF
			1. 计算->ALU
			2. 结果给上图中的Write Data
		2. load:
			1. ![[5953f37a00bee34da93b9827907a8820_720.jpg]]立即数产生单元
			2. use ALU to caculate address,use 立即数产生单元(ImmGen)从指令中产生立即数
			3. 读数据->使用Data Mem(数据存储器),memRead控制读取,memWrite控制写入
			4. 写寄存器->使用RF(寄存器组)
			5. 步骤最多,花费时间最长,决定了cpu周期的最小值
		3. sw:
			1. ALU,ImmGen
			2. write data -> Data Memrory
		4. branch instruction(以beq为例)
			1. 是否相等->相减(ALU)
			2. 根据Zero是否等于1
				1. if zero\==1,pc=pc+Imm * 2 ->use another adder and shifter to caculate address
				2. if zero\==0,pc=pc+4

![[Pasted image 20241018150350.png]]

## 控制通路设计

![[Pasted image 20241018153301.png]]

根据需要的元件设置不同的控制信号来选择具体单元

> 根据真值表设计门电路：
1. 组合逻辑->快速,但是设计复杂
2. 存储型控制器->只读存储器ROM,读入指令的opcode和fun3和fun7,转化为rom地址,根据地址中指令类型,输出控制信号

>指令所需时间控制cpu周期长短

[[Chapter_04.pdf#page=40&selection=0,26,0,26|Chapter_04, 页面 40]]

cpu周期=800ps

## 多周期处理器

cpu周期:每一步中最长的时间(如内存读)
then lw 需要5步,需要5个周期,beq 3,sw 4,R 4

不同的指令需要不同的周期数
then 功能单元发生变化

为了防止lw sw等读取数据冲掉指令,需要一个指令寄存器,为了保存数据,需要增加一个数据寄存器

同单周期处理器,从寄存器中读出的数据存到额外两个数据寄存器,保证了ALU输入的数据不变

## pipe line

homework:1 6 7 13

> 流水线级数

做完一个事需要几个步骤

例如5级流水线

> 指令执行阶段

1. IF:instruction fetch 取指
2. ID:instruction decode & register read 译码和寄存器读写
3. EX:execute operation or caculate 计算 
4. MEM:access memory operand 访存
5. WB:write back to register 写回

cpu周期取上述步骤中所需时间最长的一步
这样的流水线称为均匀流水线
也有非均匀流水线

==每条指令的时间并没有减少==

> 加速比

加速比用于衡量性能
加速比取决于吞吐量(throughput->单位时间执行的指令的条数)

指令延时未减少

## 流水线设计

1. IF->pc,IM,adder
2. ID->RF
3. EX->ALU
4. MEM->DM
5. WB->RF

==处理RF重复使用问题:==
前半周期RF写,后半周期RF读->前后之分解决后一条指令需要用到前一条指令计算结果的情况such as:`add x7,x8,x9` `sub x10,x7,x11`

为了保存本级流水线执行的信息,在每个状态之间需要加上一个寄存器

![[Pasted image 20241025144544.png]]
### 流水线设计中的问题

1. 结构冒险:一级中可能需要同时执行两个操作->增加功能单元
2. 数据冒险:指令间的数据相关联(以上述为例)->[[Chapter_04.pdf#page=51&selection=0,9,10,12|Chapter_04, 页面 51]]
	c1 c2 c3
add IF  ID  EX
sub      IF   ID(此时x7的数据还未更新)
	how to tackle?
	1. 阻塞(stall)
		1. 上述sub指令ID后持续译码,在第五个周期读出来的是正确的数据(前半周期写后半周期读),阻塞了两个周期
		2. 什么时候阻塞,就持续进行IF或者ID
	2. "旁路"(前推)[[Chapter_04.pdf#page=52&selection=0,9,10,26|Bypassing(forwarding)]]
		1. ALU输出送回ALU输入(rs1和rs2都要)->EX旁路
		2. `lw x1 0(x2) sub x4,x1x5` ->MEM旁路
			       c1   c2   c3   c4        c5
		lw     IF     ID   EX   MEM   WB
		sub           IF    ID    ID       EX(接收上条指令从mem到EX的数据)  仍然阻塞了一个周期
		==R型的数据冒险可以完全解决,load指令仍需阻塞一个周期==
	3. [[Chapter_04.pdf#page=54&selection=0,9,10,31|编译器优化]]
3. 控制冒险:分支指令要在第四个周期才能计算出pc+imm * 2,但是此时beq指令下已经进来了三条pc+4指令,即存在错误指令
	1. 阻塞(stull)
		1. 后取出来的指令不执行,等着
	2. 预测
		1. 静态预测:risc-v总预测不发生分支
		2. 动态预测:软硬件结合

### 具体设计

流水线cpu:
	1. 数据通路
	2. 控制通路

数据通路上面已有
控制通路:

流水线分析图:单周期分析图

> 单周期分析图->某个周期的情况

1. load
	1. IF:read IMM and pc+4->pc
		1. IF/ID中有指令和pc
	2. ID:decode instruction,read rs1 rs2,immGen
		1. ID/EX中有rs1,rs2,imm,pc
	3. EX:计算
		1. EX/MEM有结果,pc+imm * 2 ,rs2
	4. MEM:读内存
		1. MEM/WB中存读取到的数据
	5. WB:
		1. 把读取的数据存回rd->有错误
		2. so we need to store rd!->[[Chapter_04.pdf#page=68&selection=10,0,10,10|Chapter_04, 页面 68]]


> 多周期分析图->多个周期的情况

![[Pasted image 20241025150204.png]]

### 控制信号 

ALU控制器不变

IF 和 ID不需要控制信号

EX需要ALU control AluSrc

MEM需要 MemWrite MemRead Branch PCsrc

WB需要MemToReg

> 流水线寄存器存储

控制器放在ID阶段,控制信号存储在流水线寄存器中[[Chapter_04.pdf#page=76&selection=10,0,10,17|Chapter_04, 页面 76]]

### 实现数据冒险解决方案->旁路

> EX旁路

ALU输入部分多加一个多路选择器,多加一个控制信号选择数据来源
控制信号怎么产生?
==条件:上一条指令的rd是下一条指令的rs1或rs2==
==在上条指令的第四个周期时,rd在EX/MEM寄存器中==
==下一条指令的rs1和rs2在ID/EXE中==
==EXE/MEM中的rd!=0==
==且上条指令一定有EX/MEM中RedWrite=1==

即:
ex/mem rd=id/ex rs1
ex/mem rd=id/ex rs2
且
ex/mem rd != 0
ex/mem RegWrite  =  1

> 为了实现MEM旁路

第一条指令的rd是第三条指令的rs1或rs2

mem/wb rd= id/ex rs1
mem/wb rd =id/ex rs2
且
mem/wb rd != 0
mem/wb RegWrite = 1

### load 数据冒险

[[Chapter_04.pdf#page=89&selection=10,1,12,15|Chapter_04, 页面 89]]

需要阻塞一个周期，用到MEM旁路

条件：
ID/EXE rd=IF/ID rs1 or ID/EXE rd=IF/ID rs2
且ID/EXE MemRead=1
且rd!=0

阻塞的方法:

1. 写进ID/EXE的控制信号清零
2. 第三个周期结束时不写IF/ID寄存器
3. PC不写

### branch控制冒险 静态预测

[[Chapter_04.pdf#page=93&selection=10,0,10,14|Chapter_04, 页面 93]]

因为MEM阶段结束后,pc才拿到pc+imm * 2
所以跳转到的指令第五个周期才能取到正确的指令地址
有三条指令执行错误

>静态预测

只需要在发现分支后把EX/MEM,ID/EX,IF/ID清空即可(控制信号置0)
浪费三个周期

> EX

将pc+imm * 2操作在EX阶段执行,可以少浪费一个周期

> ID

把branch判断移到ID阶段,两个ReadData进行异或来判断zero信号

再节省一个周期

但是!如果有了数据相关和数据冒险的话!
[[Chapter_04.pdf#page=97&selection=10,1,14,14|Chapter_04, 页面 97]]
需要额外增加到ID的旁路选择
[[Chapter_04.pdf#page=98&selection=13,1,14,14|Chapter_04, 页面 98]]
额外阻塞一个周期
[[Chapter_04.pdf#page=99&selection=10,0,10,9|Chapter_04, 页面 99]]
额外阻塞两个周期

### branch动态预测

挪到IF阶段,进行动态预测
此时没法判断指令类型

> 增加一个缓冲区(branch history table,BHT)

表中存已经发生分支的指令

| PC   | 指令  | 状态位   |
| ---- | --- | ----- |
| addr | beq | 0(或1) |
|      |     |       |
IF阶段在表格中寻找PC,若存在->状态为1->分支;状态为0->不分支
若不存在,pc+4

> 一位预测器

[Chapter_04, 页面 101](附件/Chapter_04.pdf#page=101&selection=10,0,12,26)

> 两位预测器

![](附件/Pasted%20image%2020241101144824.png)

| pc  | 指令  | 状态位(两位)00 01 10 11 | 分支地址 |
| --- | --- | ------------------ | ---- |
|     |     |                    |      |
00 01 不发生分支
10 11 发生分支

如果预测00,且实际没有分支,则保持00,若实际分支了,则变为01
在01状态下若发生分支,变为10;若没有发生分支,变为00
10预测发生,若分支,变为11;若没有分支,变为01
11若分支,保持11;没有分支,变为10
保证了二次循环即以上只有最后一次内循环发生错误

### 异常和中断

1. 软件引起:异常 exceptions
	1. syscall
	2. undefined opcode
2. 硬件引起:中断 interrupts cpu中断当前程序,执行别的程序,再回来当前的程序
	1. 内存
	2. 网卡
	3. 等
实现:硬件+软件


只要发生中断,risc-v固定转到`0000 0000 1C09 0000hex`
硬件增加一个SCAUSE寄存器(异常原因寄存器),存储异常种类,64位,每一位存储一种异常种类,发生对应异常则对应位为1,异常处理程序先读SCAUSE

>intel处理器有256个异常,存储在系统初始化时在内存中固定位置创建的表中,异常号叫做中断向量号,表中存储每个中断向量号对应的错误处理程序的入口地址.这种方式叫做向量中断.

SPEC中断程序计数器,保存中断程序的返回地址,用于执行完异常处理程序后返回

如果异常处理完毕,则可以返回源程序继续执行
若异常无法处理,则强制结束程序
死机

### 异常的硬件处理

[Chapter_04, 页面 110](附件/Chapter_04.pdf#page=110&selection=10,0,10,24)

在数据通路中额外增加寄存器

```risc-v
40 sub x11, x2, x4 
44 and x12, x2, x5 
48 or x13, x2, x6 
4c add x1, x2, x1 
50 sub x15, x6, x7 
54 ld x16, 100(x7)
```

![](附件/Pasted%20image%2020241101152140.png)

假设add出问题
则并行执行:
1. add指令的pc即0x4c存到SPEC;设置SCAUSE相应位为1
2. 0x1c090000->pc
3. EX/MEM WBcontrol=0,EX/MEM Mcontrol=0 相当于把add变成nop
4. 后续指令清零: ID/EX \[WB,M,EX\] control=0 相当于把sub变成nop
5. IF/ID 寄存器flush(置零) 相当于ld指令变成nop

> 异常信号INT传给cpu相应引脚

# 第五章 存储层次

存储容量:K($2^{10}$),M,G,T,P,E,Z,Y
基本单位B

## 存储层次

cpu:寄存器
cpu cache:高速缓冲存储器(多级cache)
物理内存:RAM(主存) flash ssd
辅存:机械硬盘

![](附件/Pasted%20image%2020241106100101.png)

由于处理器和存储器gap越来越大,计算机设计由以cpu为中心->存储器为中心

访存过程:
1. cpu给出内存地址
2. 给出读或写信号
3. 写或读数据

若mem中无所需数据,则需访问输入输出接口

## 局部性原理

冯诺依曼机采用存储程序原理,相邻指令和数据存储在相邻存储单元中,一段程序在内存中占有的是一段连续的存储单元

> 时间局部性

最近访问过的数据可能再次被访问
循环控制变量

> 空间局部性

临近最近访问的数据的数据可能马上被访问
连续内存访问,比如数组

充分利用局部性:
1. 内存层次
2. 所有内容都在硬盘里
3. 数据从硬盘到mem->cache

## 内存层次工作原理

每两层之间的数据交换以块(block)为单位

cache->SRAM(static random access memory)
主存->DRAM(dynamic ...)
固态盘->flash
机械盘->disk

## RAM

RAM:随机,易失性存储器
1. SRAM:以MOSS管的导通和截止分别表示1和0,每个单元以矩阵形式构成存储芯片,一个cell需要6个moss管,集成度低,单位价格高,因为moss管状态转换速度快,所以sram读取速度快
2. DRAM:单个moss管,以电容上有无电荷表示1/0,读后会导致电容放电,破坏原有状态,所以叫做破坏性读,因此读过程中包含着写.因为电容充电较慢,所以速度较慢.因为只有一个moss管,所以集成度高,价格便宜.动态:因为电容有一定电阻,所以会漏电荷,每2ms一个单元的电荷就会漏完,所以每2ms要补充电荷,即刷新(读的过程).若刷新时与cpu读产生冲突,刷新的优先级高.DRAM分为多个bank,提高性能.

随机:读写时间和位置无关

DRAM高级技术:
1. 突发模式:读一个单元后,后续单元读不再发送地址和读信号,直接读
2. DDR(double data rate):一个cpu周期内上升和下降沿都使用
3. QDR(quad data rate):分开DDR input和output
4. ROW buffer
5. 同步:SDRAM:采用突发模式读取连续数据
6. banking:多体交叉存储器
	1. 4体
	2. ![](附件/e5c03b0fd6a158390ea0ebe9a4869deb_720.jpg)

## flash storage 固态盘

非易失性半导体存储器

NOR flash:按位读写,用于嵌入式系统
NAND flash:按块存储

每个存储单元只能写1000次
为了保证每个单元写的次数较为接近,采用均衡写算法,延长flash存储器的寿命

## 机械硬盘

使用磁性材料
需防尘防潮
采用文件系统存储

采用磁极的方向来表示0或1

同心圆->磁道->扇区(sector)->

写:
磁头是一个U型磁铁,通过写入不同方向的电流,改变磁性材料的磁化方向

读:
因为磁场不均匀,磁头在移动过程中会产生感应电流

读取:
1. 寻道
2. 旋转(扇区旋转到磁头下)
3. 读取

在使用之前要先格式化,分配文件系统等

每个扇区有一个id,512字节(等)的数据,纠错码,gaps

磁盘平均读写时间:
1. 队列延时如果有其他请求
2. 寻道时间(随机,取平均)
3. 旋转时间(随机,取平均)
4. 数据传送时间
5. 控制器时间
eg.[Chapter_05-RISC-V, 页面 26](附件/Chapter_05-RISC-V.pdf#page=26&selection=2,0,2,18)

操作系统调度算法可以做到减少寻道
智能硬盘控制器可以预测要读的扇区(SCSI,SATA等接口)
硬盘里加入cache,要写回硬盘的数据不仅放入硬盘,还放入cache

## Cache Memory

cpu现在cache中找要求的数据,若找到即为命中(hit),若没有即为缺失(miss),此时cache从mem中读取要求的数据,cpu再从cache中读取数据.
cache是mem部分内容的映像
mem和cache以块(block)为单位传输数据,块也叫cache行.块的大小不固定

块的地址按照每块中字节地址的共同位决定
内存的每个地址被分为两个部分,块号和块内地址,比如32字节的内存,每个块有4个字节,自己画一下就知道了

| 块号  | 块内地址 |
| --- | ---- |

==cache没有地址==

主存的块向cache中调度时有一个调度方式:
1. 直接映射:主存块号%cache块数

因为主存有多个号映射到cache的同一位置,为了标记cache存储的是主存的哪一块,所以给cache增加一个tag位,tag中存储主存块号高位(根据映射算法定)

增加一个有效位,来表示cache的数据是否有效

字节地址:根据地址找到一个字节
字地址:根据地址找到一个字

增大块大小:提高命中率->空间局部性
但是!块数变少,所以主存更多的块会映射到同一块中->频繁调度->miss rate upup->块污染
需要更多时间传输地址,流水线阻塞时间长->缺失代价增大

软件调度算法:
早启动:取到所需数据,流水线就流动(IM中效果好)

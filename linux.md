[[Unix  Linux大学教程 (哈恩) (Z-Library).pdf]]

# 零碎知识&命令

unix是一个多用户操作系统

unix=内核+实用工具

修改密码：passwd

last [username]:

​	显示用户上次或最近几次的登陆时间

---

运行级别:

| 0    | 关机               |
| ---- | ------------------ |
| 1    | 单用户模式(命令行) |
| 2    | 非标准化           |
| 3    | 多用户模式:命令行  |
| 4    | 非标准化           |
| 5    | 多用户模式:GUI     |
| 6    | 重新启动           |

---



> 使用 su  口令并输入密码成为超级用户



---



关机：shutdown [now] [time]

重启:reboot

init命令用于重新设置启动级别

> 需要sudo权限



---



dmesg

显示系统启动与关闭时的信息



---



^M 和 回车发送CR信号 ---返回信号

^J 发送LF信号---换行信号

- 可以认为是相同的

总结:

1. 返回字符=^M
2. 换行字符=新行字符=^J
3. 一般而言,每行文本必须以一个新行字符结束
4. 按下回车,发送一个返回字符,unix自动将返回字符变为新行字符
5. 终端上显示数据,每行必须是"返回+新行"结束



---



<^J>**stty sane**<^J> 

<^J> **reset** <^J> 

 复原终端

# 可立即使用的程序

## which type
查找程序位于什么位置，type的输出信息更加详细，type为终端内置程序

## 停止程序的三种方法

CTRL+D：告诉程序当前输入已经结束

CTRL+C：发送intr信号

quit

## date 

date查看当地时间

date -u 显示格林威治时间

## cal
显示日历

cal显示当月

cal 2024 显示某一年

cal 12 2024 显示某年某月

- 显示某月日期必须加上年份

cal -j 12 2024 ： 显示某天是这一年的第几天

## 系统信息 uptime hostname uname

- uptime 显示系统已经运行了多久
users：用户个数
load average 前1 5 15 分钟等待执行的程序的数量

- hostname
计算机名称
- uname
操作系统名称
uname -a 显示详细信息

## 自己信息 whoami quota
whoami显示当前用户

## 其他用户信息 user who w

users 登录系统的用户的名称

who：谁在哪个终端上什么时间从哪里登录

w：用户，终端，来源 登陆时间 自上次键入空闲时间 登陆后所有进程使用的cpu时间 当前进程使用的cpu时间 在做什么

w mirai 指定名称查看某个用户的具体情况

## 离开提醒

leave 提醒还有多久该离开

leave 1344 =》 13：44

leave +15 =》15分钟后离开

- 若输入的时间小于等于12，则默认时间为接下来12小时之内的时间

## 内置计算器 bc

bc -l 使用bc自带的数学库启动



支持正常优先级运算



小数点精度scale=

- -l启动scale默认为20



bc中可使用变量，只能用小写字母，且不能组合



ibase ：输入基

obase：输出基

使用ABCDEF永远代表对应十进制数字



## dc计算器——基于逆波兰表达式



# 文档查询

## man

less使用

backspace b g G / ？ n N



## 手册组织方式



| 1    | 命令     |
| ---- | -------- |
| 2    | 系统调用 |
| 3    | 库函数   |
| 4    | 特殊文件 |
| 5    | 文件格式 |
| 6    | 游戏     |
| 7    | 杂项信息 |
| 8    | 系统管理 |

手册分为8节，每节对应内容如上



## man中指定节号

```
man 1 kill

```

## whatis

显示命令的描述



## apropos

在name节中搜索包含特定关键字的命令



## info

树结构

快捷命令



# 命令行语法



## 一次输入多个命令

使用';'隔开命令即可



## 命令语法



```
命令名称 选项 参数
```

---

选项:

选项可以结合在一起

单个连字符跟单个字母

双连字符跟全拼

---

注意:

选项须在参数之前

一个或多个

零个或多个(留意默认值)

---

==语法==

学好一个命令的三个问题

- 命令是做什么的
- 如何使用选项
- 如何使用参数

**命令语法规则**

1. 方括号中的项是可选的
2. 不在方括号中的项是必选项
3. 黑体字必须按原样准确输入
4. 斜体字可以用适当值替代
5. 后面接省略号代表参数可重复多次

---

更复杂的规则

6. 如果一个单独的选项和一个参数组合在一起,那么该选项和参数必须同时使用

```bash
man [-P pager] [-S sectionlist] name...
```



7. 由竖线(|)分开的两个或多个项表示可以从列表中选择一个项

eg.

```
who ... [file  | arg1 arg2]
```

表示可以提供一个文件或者两个名为arg1 arg2 的参数

## man手册 

通常把选项用options代替



# shell



## shell介绍

bourne shell

c-shell

临时改变shell

永久改变shell



## 变量和选项



### 交互式和非交互式shell

交互式即由人输入和显示给人看

非交互式即执行脚本的shell

### 环境 进程 变量

- 环境:

一组用来存放信息的变量

- 变量:

存储数据的实体,变量名和值

- - TERM变量存储所使用的终端类型

由shell启动的程序,父进程和子进程

子进程继承父进程的一切环境



### 环境变量和shell变量

**局部变量 全局变量**

- 局部(shell)变量

1. 存放对shell本身有意义的信息,如ignoreeof选项确定是否忽略eof信号
2. 在shell脚本中以普通程序中局部变量的方式使用

---

变量存在的问题:

有的变量同时有局部变量和全局变量的含义

- bash

将变量只定义为局部变量或者同时定义为局部变量和全局变量

*不能有变量只是全局而不是局部*

创建变量时变量自动被设置为shell变量

可以使用==export==命令将变量变成环境变量

```bash
1.
	SPORT=surfing # 注意等号左右两边没有空格
	export SPORT
2. 
	export SPORT=surfing
```



### 显示环境变量 env printenv

env显示默认变量

printenv同样



可以使用管道和sort进行按字母表排序



### 显示shell变量 set

set显示shell变量以及他们的值



### 显示以及使用变量的值 echo print

```bash
echo i love you

echo $TERM

echo 'my name is <$USER>'

echo $USER $TERM $PATH

echo "My favorite sport is ${ACTIVITY}ing" #花括号保证和周围不分开


```



### Bourne shell export unset

```bash
export NAME[=value]...
```

- 删除变量

unset

**称为复位变量**



### shell选项 set -o set +o

```bash
#设置一个选项
set -o option
# 复位一个选项
set +o option

```



### 显示shell选项

```bash
set -o#人类阅读友好 
set +o#脚本友好
```



## 命令和定制

### 元字符

![image-20240805170058015](C:\Users\32284\AppData\Roaming\Typora\typora-user-images\image-20240805170058015.png)



### 引用和转义



按照字面解释元字符:

反斜线,单引号,双引号

---

\

转义字符

改变的是元字符的模式



---

''

单引号

强引用

单引号中的所有字符都按照字面解释

---

""

双引号

弱引用

保留$ ` \ 的特殊含义



---



### 强引用和弱引用

- 强引用

\ 的引用是最强的

可以引用新行字符

\return 将会开启新的一行,但不表示结束(即不具有新行的含义)

eg.

~~~bash
echo this is a long \
long sentence that is \
need to type a long time.

~~~

---

- 单引号

需要一对儿来标记结束

bash家族shell会等待输入第二个引号

c-shell家族会输出错误



### shell内置命令 type

使用

``` type commandname``` 查看命令是否是内部命令



---

搜索shell手册中的builtin来学习内部命令



使用```help -s command``` 来查看命令的语法



### 外部命令以及搜索路径

shell在PATH环境变量中查找外部命令



### 修改搜索路径

使用export命令使PATH成为环境变量

```bash
export PATH="/bin:/usr/bin:usr/ucb:/usr/local/bin"
```

各个名称之间使用冒号隔开,等号两边无空格

该命令应置于登陆文件中

----

也可以添加自己的目录

```shell
export PATH="$PATH:$HOME/bin"
```

- 尽量将自己设置的目录放在系统原有设置目录之后,否则容易覆盖系统目录



### shell提示

使用```export PS1="$ "``` 来设置提示符

---

可综合使用变量的值

eg.

$USDER 

$PWD

下面列出了常用环境变量:

![image-20240805172306568](C:\Users\32284\AppData\Roaming\Typora\typora-user-images\image-20240805172306568.png)

### 引用变量时使用的引号

变量有变化使用单引号

变量不变化使用双引号

eg.

```bash
export PS1='Your lucky number is ${RANDOM} $ '
export PS1="${USER}"
```

第一条里$ 符号保留为\$符号,在使用这个变量时才会解释为元字符

第二条里 直接被解释为元字符,对变量进行展开



### 转义字符的特殊码

![image-20240805172910768](C:\Users\32284\AppData\Roaming\Typora\typora-user-images\image-20240805172910768.png)

### 命令替换

允许在一条命令中嵌入另一个命令

```bash
echo "The time and date are 'date'."
```

先运行date命令,在将date命令的结果嵌入到echo命令中传递给echo命令



---

basename可以抽取任何路径名的最后一个部分

---



### 键盘输入删除快捷操作

^W删除刚才输入的一个单词

\^U或\^X删除整行

使用stty显示系统上所有的键映射

^D删除光标后的一个字符



### 历史列表fc history

使用上下箭头查看一条历史命令

每一个命令在历史列表中称为事件

---

fc -l

history

都可以显示历史列表

---

fc -s num

! num

使用编号为num的命令

---

fc -s

!!

使用上一条命令

---

fc -s pattern=replacement number

pattern和replacement都是字符串 number是事件编号

eg.

**25 vi temp**

**fc -s temp=data 25**

---

也可以修改某一部分

eg.

```bash
datq
fc -s q=e
```



### 设置历史记录大小

```bash
export HISTSIZE=50
```



- 养成先查看文件再删除文件的习惯

> 使用ls和rm结合,用历史替换进行快速操作



### 自动补全

bash支持:

1. 文件名和目录名
2. 命令
3. 变量
4. 用户标识
5. 主机名(联网)

- 连续按tab键会显示出所有匹配的选项



### 别名 alias

```bash
alias [name=commands]
```

等号两边不能有空格

commands包含一个或多个命令,看情况使用引号

使用```unalias``` 移除别名

---

使用```\ls``` 临时挂起别名

### 别名示例,避免删除文件

```bash
ls temp*
fc -s ls=rm
alias del 'rm \!ls:*'

```

### 从历史列表中重用命令

```bash
alias h=history

alias r="fc -s"
# this is a example
vi tempfile
r tempfile=data
 
#gei a alias to alias it self
alias a=alias

```

 ## 初始化文件



初始化文件:

登陆文件

环境文件



---



### 点文件和rc文件

点文件是隐藏文件

使用ls -a查看所有文件

环境文件以rc结尾,意为:running command



### 执行初始化文件的时机

登录shell执行登陆文件和环境文件

非登录shell执行环境文件

- bash的登录shell只执行登录文件而不执行环境文件

### 脚本注释

shell使用#作为注释标志










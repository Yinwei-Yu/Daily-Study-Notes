# 基础
## 数据类型和变量
### 数据类型

- 整数
支持使用下划线分隔
- 浮点数
科学计数法$1.23*10^{23}=1.23e23$
- 字符串
单引号或双引号括起来的文本
文本内需要引号请使用转义字符==\\==
r' '表示''内部的字符串默认不转义
‘’‘...'''表示多行文本
- 布尔数
- 空值
None表示，不等于0
### 变量
变量
常量——python中没有常量，使用全大写字母来表示常量
## 字符串和编码
python采用Unicode编码
- ord()——获取字符的正数表示
- chr()——把编码转换成对应的字符
- 
python字符保存到磁盘上需要从str转化为bytes，使用带b的前缀表示
~~~python
x=b'ABC'
#注意与'ABC'的区别
~~~
以Unicode表示的str可以使用encode（）方法编码为指定的bytes
decode（）方法从bytes转化为str

len()计算字符串长度，若转化为字节则计算字节数

在python开头写下下面两行
~~~python
```
#!/usr/bin/env python3
# -*- coding: utf-8 -*-
```
~~~
第一行注释是为了告诉Linux/OS X系统，这是一个Python可执行程序，Windows系统会忽略这个注释；

第二行注释是为了告诉Python解释器，按照UTF-8编码读取源代码，否则，你在源代码中写的中文输出可能会有乱码。

python中格式化输出的方式与c语言相同，采用%
~~~python
'Hi, %s, you have $%d.' % ('Michael', 1000000)
~~~
常见的占位符如下
%d——整数
%f——浮点数
%s——字符串
%x——十六进制整数

format()
使用传入的参数依次替换字符串内的占位符{0}、{1}等

f-string
以f开头的字符串，字符串内{x}的内容会被x的内容替换

## 列表和元组
### list
~~~python
name=[a,b,c]
>>>len(name)
3
~~~
-1获取倒数第一个元素，-2倒数第二个，以此类推
insert（position）
pop（）删除末尾元素
pop（position）删除指定位置元素
name\[i]=1，直接赋值
list元素类型可以不同
list元素可以是另一个list
### tuple
元组一旦初始化就不能修改，指的是元组指向的变量不能修改
~~~python
classmates=("a","b","c")
~~~
可以使用索引访问
空元组：t=（）；
一个元素t=（1,） **逗号表元组**

## 条件
### if
~~~python
age = 20
if age >= 6:
    print('teenager')
elif age >= 18:
    print('adult')
else:
    print('kid')
~~~
int（）函数将字符串转化为整数
float（）将字符串转化为浮点数

### match
xxxxxxxxxx7 1​2ON DELETE SET NULL;-- 外键被删除后本表内对应值为NULL3​4ON DELETE CASCADE;-- 外键删除后本表内对应值删除5​6-- 如果一个属性既是primary key 又是 foreign key,则不能设置为SEY NULL7​mysql
~~~python
age = 15
match age:
    case x if x < 10:
        print(f'< 10 years old: {x}')
    case 10:
        print('10 years old.')
    case 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18:
        print('11~18 years old.')
    case 19:
        print('19 years old.')
    case _:
        print('not sure.')

~~~

~~~python
args = ['gcc', 'hello.c', 'world.c']
# args = ['clean']
# args = ['gcc']

match args:
    # 如果仅出现gcc，报错:
    case ['gcc']:
        print('gcc: missing source file(s).')
    # 出现gcc，且至少指定了一个文件:
    case ['gcc', file1, *files]:
        print('gcc compile: ' + file1 + ', ' + ', '.join(files))
    # 仅出现clean:
    case ['clean']:
        print('clean')
    case _:
        print('invalid command.')

~~~

## 循环

### for
for x in y
	do... 

- range()生成一个整数序列
- list（）将一个序列转化为list

### while
while n>0
	do

- break
- continue

## dict和set
### dict 字典
相当于map
键值对\[key,value]
~~~python
d={"mio":100,"mugi":99,"retsu":80,"yui":59}
~~~
修改key对应的值
~~~python
d['yui']=100
~~~

~~~python
'miku' in d
#判断key是否存在
~~~

~~~python
#get函数判断
#第二个参数为自定义返回值，默认找不到不返回
d.get(key,retrun_value)
#不在返回return_value，在返回键值
~~~

key值计算位置采用hash算法
key值需为字符串，整数等不可变值

pop（key），删除指定键值对
### set
集合
内无重复元素
需要list作为输入集合，使用set（）函数转化
~~~python
s=set([1,2,3,3])
~~~

集合自动忽略重复元素

add（key）
remove（key）

~~~python
s1=set([1,2,3])
s2=set([2,3,4])
s1 & s2 #交集操作
s1 | s2 #并集操作
~~~

# 函数

[python内置函数](https://docs.python.org/zh-cn/3/library/functions.html#exec)

[常用函数](https://www.cnblogs.com/yund/p/17370975.html)

请记得在函数之后加上":"
~~~python
def love(name):
	if name=="miku":
		return "you love miku"
	else:
		return "you love nothing"
~~~

- 空函数 
pass语句作为占位置，表示什么也不做，不加会报错

- 请对参数进行错误检查

- 函数可以返回多值，实际上返回的是多个值组成的tuple

## 默认参数
~~~python
power(x,n=2)
#这里n=2意思是不传入第二个参数n取2，传入n后按传入的参数进行计算
~~~
- 默认参数在必选参数之后
- 有多个默认参数时可以在调用时指定参数名给出需要传入的参数
- 默认参数需要为不变对象

## 可变参数
可以传入任意个数变量，在函数内自动组装为一个tuple
~~~python
def Cal(*num)
	sum=0
	for x in num:
		sum=sum+num*2
	return sum
a=[1,2,3,5]
B=Cal(1,2,3,5)
C=Cal(*a)
~~~

\*告诉函数这是一个可变参数，可以接收任意多参数
在列表或元组前加上\*可将list或tuple的所有元素作为可变参数传入

## 关键字参数
允许传入0个或任意多个含参数名的参数
~~~python
def createUser(name,age,**ot)
	print("name:",name,"age:",age,"other:",ot)

createUser("mio",18)#只传入必要参数
createUser("miku",17,gender="female",career="singer")#传入任意个数关键字参数
dt={"gender":"female","career":"singer"}#传入元组
createUser("miku",17,**dt)
~~~

注意ot获得的是元组的拷贝，修改ot不会修改元组的值

## 命名关键字参数
为了限制使用关键字参数传入的参数名称
~~~python
def person(name,age,*,city,job):
	print(name,age,city,job)

person("miku",17,city="Tokyo",job="virtual singer")
~~~
\*后的参数被解释为命名关键字参数，只接受名称为city和job的参数
若函数定义中有一个可变参数，则可变参数后的参数视为命名关键字参数
- 命名关键字必须传入参数名
命名关键字参数可以有默认值，从而可以在调用时不传入

## 参数组合
- 参数定义的顺序为：
==必选参数、默认参数、可变参数、命名关键字参数、关键字参数==

## 函数递归
没什么好说的，注意写递归出口
递归深度过深会溢出

# 切片
取列表或者元组的一部分
L\[x:y:z]
x表示起始下标
y表示取元素个数
z表示步长

- 注：x可以为负数，表示从后开始取，此时仍然向后增长
- 字符串也支持切片操作

# 迭代

只要是迭代对象，即可使用for循环遍历

可以使用collections.abc模块的Iterable类型判断
~~~python
from collections.abc import Iterable
isinstance('abc',Iterable)#str是否可迭代
~~~

对于字典可以迭代不同的对象
~~~python
d={'a':1,'b':2}
for key in d：
	print(key)#默认为key值
for value in d.values()
	print(key)
for k,v in d.items()
	print(k,':',value)
	
~~~

# 列表生成式
简单而又强大的生成工具
下面是一些例子，一看便知
~~~python
[x*x for x in range(1,11)]
[x*x for x in range(1,11) if x%2==0]#条件筛选
[m+','+n for m in "ABC" for n in "DEF"]#二重循环
~~~

- 注：
	- for后if不能加else
	- for前if需要else
# 高阶函数

## map
map()函数接受两个参数
第一个参数是一个函数名f，第二个参数是一个可迭代对象l
map函数对l中的每一个元素进行f操作
返回一个迭代器iterator
可以使用list ()函数将迭代器转化为列表

## reduce
 定义在functools中
reduce函数接收的参数和map函数相同
不同是，reduce函数将结果和序列的下一个元素继续做累积计算

## filter
同样接受一个函数f和一个序列
函数返回值需要是一个bool类型
filter函数根据f的返回值为true 还是 false对序列中的元素进行删除或者保留
filter函数返回一个iterator，需要用list函数将其转化为一个列表


## sorted
sorted可以对列表内元素进行排序
可以接受第二参数key（函数）
key=f，f定义了对列表中元素进行的操作
可以传入第三个函数 reverse=True，实现从大到小排序

# 返回函数

# 面向对象
~~~python
class Student(object):
	def _init_(self,name,score)
		self.name=name
		self.score=score
~~~

\_init_函数相当于c++里的构造函数，第一个参数永远是self

其他函数类似，第一个参数也是self，其他和普通函数没有区别

## 访问限制
变量名前加两个下划线__，成为私有变量，外部无法访问

## 继承和多态

子类的函数会覆盖父类同名函数

## 获取对象信息

type（）函数可以判断对象类型，返回值为对应的类
type中定义了一系列常量用以判断对象的函数类型：
- FunctionType    自定义函数
- BuiltinFunctionType  内建函数
- LambdaType  Lambda表达式
- GeneratorType  生成器类型

isinstance（）函数判断一个变量是否是一个或某个类型
ed.
~~~python
isinstance([1,2,3],(list,tuple))
~~~


dir()方法获得一个对象的所有属性和方法
返回一个包含字符串的list

hasattr（obj，x）方法判断对象obj中是否含有属性x
setattr（obj，x）方法在对象中设置一个属性y
getattr（obj，y，404）在对象obj中获取属性y，第三个参数是属性不存在时的设置返回值

## 实例属性和对象属性

对象可以任意添加属性

类的属性需要在类中添加

## __slots__

 slots是一个特殊的变量，它指示当前类的变量可以绑定哪些属性

eg.
~~~python
class Student(object):
	__slots__=('name','age')

s=Student()
s.name='li'
s.score=99 #错误
~~~

注意：子类不继承父类的__slots__属性

## MinIn
在继承的类后加上MixIn
可以实现在一个类中拥有另一个类的属性和方法

## 枚举类

~~~python
from enum import Enum

Month = Enum('Month', ('Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'))
~~~

可以从枚举类派生出自定义类
~~~python
class WeekDay(Enum)
	Sun=0
	Mon=1
	Tue=2
~~~

# 错误处理

基本try except finally和c++相同

错误是类

使用raise抛出异常

# 调试
使用assert（断言）来测试
assert bool ‘message’

在程序头使用'$python -O err.py'关闭assert，可以当作pass看待

# IO操作

## 文件操作

### 读

~~~python
f=open('test.txt','r')
~~~

'r'表示读文件

打开成功后采用read（）方法可以读取文件的全部内容，把内容读到内存中
close（）方法关闭文件

由于文件读写时可能出错，后面就不会再关闭文件了，所以可以使用try 语句实现
~~~python
try：
	f=open('/usr/file','r')
	print(f.read())
finally:
	if f:
		f.close()
~~~

这样太过繁琐，使用内置with语句
~~~python
with open('file','r') as f:
	print(f.read())
~~~

- read()一次性会读取文件的全部内容
- read(size)读取指定size字节的内容
- readline()读取一行内容
- readlines()读一行
==strip()==方法删除末尾的‘\\n’

- 二进制文件：使用rb模式打开即可
- 字符编码：读取非UTF-8编码的文本文件时，给open（）函数传入encoding参数
- errors=‘ignore’参数可以忽略非法字符

### 写
使用open（）函数时传入w或者wb参数，进行写文本文件或者二进制文件
write（）函数将内容写入文件，建议也使用with语句完成
~~~python
with open('/Users/michael/test.txt', 'w') as f:
    f.write('Hello, world!')
~~~

传入encoding参数将字符串自动转换为指定编码

传入a参数以追加（append）模式写入，不会覆盖
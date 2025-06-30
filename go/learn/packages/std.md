# 9.1 标准库概述

像 `fmt`、`os` 等这样具有常用功能的内置包在 Go 语言中有 150 个以上，它们被称为标准库，大部分(一些底层的除外)内置于 Go 本身。完整列表可以在 [Go Walker](https://gowalker.org/search?q=gorepos) 查看。

- `archive/tar` 和 `/zip-compress`：压缩（解压缩）文件功能。
- `fmt`-`io`-`bufio`-`path/filepath`-`flag`:  
	- `fmt`: 提供了格式化输入输出功能。  
	- `io`: 提供了基本输入输出功能，大多数是围绕系统功能的封装。  
	- `bufio`: 缓冲输入输出功能的封装。  
	- `path/filepath`: 用来操作在当前系统中的目标文件名路径。  
	- `flag`: 对命令行参数的操作。　　
- `strings`-`strconv`-`unicode`-`regexp`-`bytes`:  
	- `strings`: 提供对字符串的操作。  
	- `strconv`: 提供将字符串转换为基础类型的功能。
	- `unicode`: 为 unicode 型的字符串提供特殊的功能。
	- `regexp`: 正则表达式功能。  
	- `bytes`: 提供对字符型分片的操作。  
	- `index/suffixarray`: 子字符串快速查询。
- `math`-`math/cmath`-`math/big`-`math/rand`-`sort`:  
	- `math`: 基本的数学函数。  
	- `math/cmath`: 对复数的操作。  
	- `math/rand`: 伪随机数生成。  
	- `sort`: 为数组排序和自定义集合。  
	- `math/big`: 大数的实现和计算。  　　
- `container`-`/list-ring-heap`: 实现对集合的操作。  
	- `list`: 双链表。
	- `ring`: 环形链表。
- `time`-`log`:  
	- `time`: 日期和时间的基本操作。  
	- `log`: 记录程序运行时产生的日志，我们将在后面的章节使用它。
- `encoding/json`-`encoding/xml`-`text/template`:
	- `encoding/json`: 读取并解码和写入并编码 JSON 数据。  
	- `encoding/xml`: 简单的 XML1.0 解析器，有关 JSON 和 XML 的实例请查阅第 [12.9](12.9.md)/[10](10.0.md) 章节。  
	- `text/template`:生成像 HTML 一样的数据与文本混合的数据驱动模板（参见[第 15.7 节](15.7.md)）。  
- `net`-`net/http`-`html`:（参见[第 15 章](15.0.md)）
	- `net`: 网络数据的基本操作。  
	- `http`: 提供了一个可扩展的 HTTP 服务器和基础客户端，解析 HTTP 请求和回复。  
	- `html`: HTML5 解析器。  
- `runtime`: Go 程序运行时的交互操作，例如垃圾回收和协程创建。  
- `reflect`: 实现通过程序运行时反射，让程序操作任意类型的变量。  

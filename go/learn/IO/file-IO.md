# 12.2 文件读写

## 12.2.1 读文件

在 Go 语言中，文件使用指向 `os.File` 类型的指针来表示的，也叫做文件句柄。我们在前面章节使用到过标准输入 `os.Stdin` 和标准输出 `os.Stdout`，他们的类型都是 `*os.File`。让我们来看看下面这个程序：

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

func main() {
    inputFile, inputError := os.Open("input.dat")
    if inputError != nil {
        fmt.Printf("An error occurred on opening the inputfile\n" +
            "Does the file exist?\n" +
            "Have you got access to it?\n")
        return // exit the function on error
    }
    defer inputFile.Close()

    inputReader := bufio.NewReader(inputFile)
    for {
        inputString, readerError := inputReader.ReadString('\n')
        fmt.Printf("The input was: %s", inputString)
        if readerError == io.EOF {
            return
        }
    }
}
```

变量 `inputFile` 是 `*os.File` 类型的。该类型是一个结构，表示一个打开文件的描述符（文件句柄）。然后，使用 `os` 包里的 `Open()` 函数来打开一个文件。该函数的参数是文件名，类型为 `string`。在上面的程序中，我们以只读模式打开 `input.dat` 文件。

如果文件不存在或者程序没有足够的权限打开这个文件，Open函数会返回一个错误：`inputFile, inputError = os.Open("input.dat")`。如果文件打开正常，我们就使用 `defer inputFile.Close()` 语句确保在程序退出前关闭该文件。然后，我们使用 `bufio.NewReader()` 来获得一个读取器变量。

通过使用 `bufio` 包提供的读取器（写入器也类似），如上面程序所示，我们可以很方便的操作相对高层的 `string` 对象，而避免了去操作比较底层的字节。

接着，我们在一个无限循环中使用 `ReadString('\n')` 或 `ReadBytes('\n')` 将文件的内容逐行（行结束符 `'\n'`）读取出来。

**注意：** 在之前的例子中，我们看到，Unix 和 Linux 的行结束符是 `\n`，而 Windows 的行结束符是 `\r\n`。在使用 `ReadString` 和 `ReadBytes` 方法的时候，我们不需要关心操作系统的类型，直接使用 `\n` 就可以了。另外，我们也可以使用 `ReadLine()` 方法来实现相同的功能。

一旦读取到文件末尾，变量 `readerError` 的值将变成非空（事实上，其值为常量 `io.EOF`），我们就会执行 `return` 语句从而退出循环。

**其他类似函数：**

**1) 将整个文件的内容读到一个字符串里：**

如果您想这么做，可以使用 `io/ioutil` 包里的 `ioutil.ReadFile()` 方法，该方法第一个返回值的类型是 `[]byte`，里面存放读取到的内容，第二个返回值是错误，如果没有错误发生，第二个返回值为 `nil`。请看示例 12.5。类似的，函数 `WriteFile()` 可以将 `[]byte` 的值写入文件。


```go
package main
import (
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    inputFile := "products.txt"
    outputFile := "products_copy.txt"
    buf, err := ioutil.ReadFile(inputFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
        // panic(err.Error())
    }
    fmt.Printf("%s\n", string(buf))
    err = ioutil.WriteFile(outputFile, buf, 0644) // oct, not hex
    if err != nil {
        panic(err.Error())
    }
}
```

**2) 带缓冲的读取**

在很多情况下，文件的内容是不按行划分的，或者干脆就是一个二进制文件。在这种情况下，`ReadString()` 就无法使用了，我们可以使用 `bufio.Reader` 的 `Read()`，它只接收一个参数：

```go
buf := make([]byte, 1024)
...
n, err := inputReader.Read(buf)
if (n == 0) { break}
```

变量 `n` 的值表示读取到的字节数.

**3) 按列读取文件中的数据**

如果数据是按列排列并用空格分隔的，你可以使用 `fmt` 包提供的以 `FScan...` 开头的一系列函数来读取他们。请看以下程序，我们将 3 列的数据分别读入变量 `v1`、`v2` 和 `v3` 内，然后分别把他们添加到切片的尾部。

```go
package main
import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("products2.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var col1, col2, col3 []string
    for {
        var v1, v2, v3 string
        _, err := fmt.Fscanln(file, &v1, &v2, &v3)
        // scans until newline
        if err != nil {
            break
        }
        col1 = append(col1, v1)
        col2 = append(col2, v2)
        col3 = append(col3, v3)
    }

    fmt.Println(col1)
    fmt.Println(col2)
    fmt.Println(col3)
}
```

输出结果：

```
[ABC FUNC GO]
[40 56 45]
[150 280 356]
```

**注意：** `path` 包里包含一个子包叫 `filepath`，这个子包提供了跨平台的函数，用于处理文件名和路径。例如 `Base()` 函数用于获得路径中的最后一个元素（不包含后面的分隔符）：

```go
import "path/filepath"
filename := filepath.Base(path)
```

## 12.2.3 写文件

```go
package main

import (
	"os"
	"bufio"
	"fmt"
)

func main () {
	outputFile, outputError := os.OpenFile("output.dat", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	outputString := "hello world!\n"

	for i:=0; i<10; i++ {
		outputWriter.WriteString(outputString)
	}
	outputWriter.Flush()
}
```

除了文件句柄，我们还需要 `bufio` 的 `Writer`。我们以只写模式打开文件 `output.dat`，如果文件不存在则自动创建：

```go
outputFile, outputError := os.OpenFile("output.dat", os.O_WRONLY|os.O_CREATE, 0666)
```

可以看到，`OpenFile` 函数有三个参数：文件名、一个或多个标志（使用逻辑运算符 `|` 连接），使用的文件权限。

我们通常会用到以下标志：

- `os.O_RDONLY`：只读
- `os.O_WRONLY`：只写
- `os.O_CREATE`：创建：如果指定文件不存在，就创建该文件。
- `os.O_TRUNC`：截断：如果指定文件已存在，就将该文件的长度截为 0 。

在读文件的时候，文件的权限是被忽略的，所以在使用 `OpenFile()` 时传入的第三个参数可以用 0 。而在写文件时，不管是 Unix 还是 Windows，都需要使用 `0666`。

然后，我们创建一个写入器（缓冲区）对象：

```go
outputWriter := bufio.NewWriter(outputFile)
```

接着，使用一个 `for` 循环，将字符串写入缓冲区，写 10 次：`outputWriter.WriteString(outputString)`

缓冲区的内容紧接着被完全写入文件：`outputWriter.Flush()`

如果写入的东西很简单，我们可以使用 `fmt.Fprintf(outputFile, "Some test data.\n")` 直接将内容写入文件。`fmt` 包里的 `F...` 开头的 `Print()` 函数可以直接写入任何 `io.Writer`，包括文件（请参考[章节 12.8](12.8.md))。

程序 `filewrite.go` 展示了不使用 `fmt.FPrintf()` 函数，使用其他函数如何写文件：

```go
package main

import "os"

func main() {
	os.Stdout.WriteString("hello, world\n")
	f, _ := os.OpenFile("test", os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	f.WriteString("hello, world in a file\n")
}
```

使用 `os.Stdout.WriteString("hello, world\n")`，我们可以输出到屏幕。

我们以只写模式创建或打开文件 `"test"` ，并且忽略了可能发生的错误：`f, _ := os.OpenFile("test", os.O_CREATE|os.O_WRONLY, 0666)`

我们不使用缓冲区，直接将内容写入文件：`f.WriteString()`

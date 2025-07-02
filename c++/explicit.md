
## ✅ 什么是 `explicit`？

`explicit` 是一个 **C++ 的关键字**，用来修饰类的构造函数（尤其是只有一个参数的构造函数），告诉编译器：

> “请不要允许隐式地用这个构造函数进行类型转换。”

换句话说：**加了 `explicit` 后，只能显式调用构造函数**，不能让编译器自动帮你做“看起来合理但可能出错”的转换。

---

## 🧠 为什么要用 `explicit`？

### 👀 情景一：没有 `explicit` 的时候

假设我们有这样一个类：

```cpp
class String {
public:
    String(const char* str) { /* 构造字符串 */ }
};
```

然后我们写代码：

```cpp
void print(String s) {
    // 打印字符串
}

print("hello");  // 这里会自动变成 String("hello")
```

你看，我们只是传了一个 `const char*` 类型的 `"hello"`，而函数 `print` 接收的是一个 `String` 类型的对象。这时候编译器会自动调用 `String("hello")` 来创建对象，然后再传递给函数。

这看起来没问题，但有时候这种**隐式转换**会导致意外行为，比如：

- 构造函数本来是用来做一些复杂的初始化的；
- 或者我们不希望用户轻易把一个原始指针转换成类对象；
- 或者我们想明确表示某些操作必须显式调用。

所以，为了避免这些潜在的问题，我们可以加上 `explicit`。

---

## ✅ 加上 `explicit` 后会发生什么？

```cpp
class String {
public:
    explicit String(const char* str) { /* 构造字符串 */ }
};

void print(String s) {
    // 打印字符串
}
```

现在如果你再写：

```cpp
print("hello");  // ❌ 编译错误！
```

编译器就会报错，因为它不能自动将 `const char*` 转换为 `String` 对象了。

你必须**显式调用构造函数**：

```cpp
print(String("hello"));  // ✅ 正确
```

---

## 🔍 示例对比

| 无 `explicit` | 有 `explicit` |
|---------------|----------------|
| ```cpp<br>print("hello");<br>``` | ❌ 报错 |
| ```cpp<br>print(String("hello"));<br>``` | ✅ 正确 |
| ```cpp<br>String s = "hello";<br>``` | ✅ 自动转换 | ❌ 报错 |

注意：即使你用了赋值语法 `=`，只要构造函数是单参数，也会发生隐式转换，除非加了 `explicit`。

---

## 🚫 隐式转换带来的问题

想象一下这个例子：

```cpp
class Resource {
public:
    Resource(int id) {
        if (id < 0) throw std::invalid_argument("ID must be positive");
    }

    void use() { std::cout << "Using resource with ID" << std::endl; }
};
```

如果我们不小心写了：

```cpp
Resource r = -1;  // ❗️ 会自动调用 Resource(-1)
r.use();          // 可能抛出异常
```

虽然你本意不是这样写的，但是因为构造函数只接受一个参数，所以它被自动调用了。如果加了 `explicit`，就可以防止这种意外。

---

## ✅ 总结

| 内容 | 解释 |
|------|------|
| `explicit` 是做什么的？ | 禁止隐式类型转换 |
| 通常用在哪里？ | 构造函数，尤其是只有一个参数的构造函数 |
| 为什么用？ | 防止误操作、提高代码安全性 |
| 如何使用？ | 在构造函数前加上 `explicit` 关键字 |
| 不加的话有什么风险？ | 允许隐式转换，可能导致难以发现的 bug |

---

## 🤔 回到你的代码中的 `explicit`

在你提供的代码中，构造函数是这样的：

```cpp
explicit Geometry(const void* wkb, size_t size);
explicit Geometry(const char* wkt);
```

这两个构造函数都加了 `explicit`，意思是：

- 不能通过简单的赋值或函数调用自动转换成 `Geometry`。
- 必须显式地用 `Geometry(...)` 创建对象。

例如：

```cpp
Geometry g = "POINT(1 2)";  // ❌ 错误！不能隐式转换
Geometry g("POINT(1 2)");   // ✅ 正确
```

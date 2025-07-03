
## 🧠 什么是 `noexcept`？

在 C++ 中，函数可以**抛出异常（throw exceptions）**，比如：

```cpp
void do_something() {
    if (something_wrong) {
        throw std::runtime_error("Something went wrong!");
    }
}
```

而 `noexcept` 是一个关键字，用来告诉编译器和程序员：**这个函数不会抛出任何异常。**

---

## ✅ 简单理解

你可以把 `noexcept` 想象成一个“承诺”：

> “我保证这个函数不会抛出异常。”

如果违反了这个承诺，C++ 标准规定会调用 `std::terminate()`，也就是程序直接崩溃。

---

## 🔧 如何使用 `noexcept`

### 示例：
```cpp
void my_function() noexcept {
    // 这个函数承诺不抛异常
    // 如果里面写了 throw，程序会直接崩溃
}
```

或者你也可以把它放在函数定义的末尾：

```cpp
void my_function();

void my_function() noexcept {
    // 实现
}
```

---

## 🚫 不加 `noexcept` 的默认行为

如果不写 `noexcept`，那函数**默认是可以抛出异常的**。

也就是说：

```cpp
void my_function();  // 默认可抛异常
```

等价于：

```cpp
void my_function() noexcept(false);
```

---

## 🔄 `noexcept` 和移动语义的关系

在你的代码中看到这样的写法：

```cpp
Geometry(Geometry&& other) noexcept { ... }
```

这是因为 **移动构造函数通常应该写成 `noexcept`**，原因如下：

- 移动操作应该是**快速、高效、不可失败**的。
- 如果移动过程中抛出了异常，那么可能会导致资源泄漏或对象处于不一致状态。
- 所以我们通过 `noexcept` 来告诉编译器：“这个移动是安全的”。

---

## 🧪 `noexcept` 的实际好处

1. **提高性能**
   - 编译器知道函数不会抛异常后，可以做一些优化（比如省略异常处理的开销）。
2. **更清晰的接口**
   - 告诉使用者“这个函数很安全，不会打断流程”。
3. **避免意外错误**
   - 如果你在 `noexcept` 函数中不小心写了 `throw`，编译器会报错或运行时直接终止程序。

---

## ⚠️ 注意事项

虽然 `noexcept` 很有用，但也有一些陷阱需要注意：

### ❌ 不要滥用 `noexcept`
如果你不能确保函数真的不会抛异常，就不要随便加上 `noexcept`。否则一旦抛出异常，程序就会崩溃。

### ✅ 什么时候适合加？
- 资源转移类函数（如移动构造/赋值）
- 一些简单的逻辑判断
- 你知道所有内部调用都不会抛异常的时候

---

## 📦 在你的代码中怎么体现？

在你提供的代码中：

```cpp
Geometry(Geometry&& other) noexcept
    : wkb_data_(std::move(other.wkb_data_)),
      size_(std::move(other.size_)),
      geometry_(std::move(other.geometry_)) {
}
```

这是 **移动构造函数**，并且加上了 `noexcept`，意思是：

> “这个构造函数不会抛出任何异常”，所以可以放心地用于需要高性能的场景。

---

## ✅ 小结一下

| 名称 | 含义 |
|------|------|
| `noexcept` | 表示函数不会抛出异常 |
| `noexcept(true)` | 明确表示不会抛异常 |
| `noexcept(false)` | 可以抛异常（默认情况） |
| 用途 | 提高性能、增强安全性、帮助编译器优化 |

---

## 💡 示例对比

```cpp
// 可能抛异常
void maybe_throw() {
    if (rand() % 2 == 0)
        throw std::runtime_error("Oops");
}

// 不会抛异常
void safe_function() noexcept {
    // 安全操作
}
```

---

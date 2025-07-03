这是一个非常经典又重要的问题，理解好 `static`、`const` 和 `constexpr` 这三个关键字，是写出 **高质量、高性能、易维护 C++ 代码的基础**。

## ✅ 一、`const`：**常量限定符**

### 🔹 作用：

- 表示**只读**，禁止修改。
    
- 可以用于变量、指针、函数参数、成员函数等。
    

---

### 📌 用法一：变量只读

```cpp
const int x = 5;
x = 6; // ❌ 错误：x 是只读的
```

---

### 📌 用法二：修饰函数参数

```cpp
void print(const std::string& msg);
```

表示函数中不会修改传入的字符串，**提高安全性**和**可读性**。

---

### 📌 用法三：修饰成员函数

```cpp
class A {
public:
    void show() const;  // 表示这个函数不会修改成员变量
};
```

---

### 📌 用法四：修饰指针

```cpp
const int* p;     // 指向的内容不能变
int* const p2;     // 指针地址不能变
const int* const p3;  // 都不能变
```

---

## ✅ 二、`static`：**静态生命周期或作用域**

### 🔹 作用（分为几种场景）：

#### 📌 1. 在函数中：**静态局部变量**

```cpp
void foo() {
    static int count = 0;
    count++;
    std::cout << count << std::endl;
}
```

- `count` 只会初始化一次，并且在多次调用中**保留值**（不像普通局部变量）
    

---

#### 📌 2. 在类中：**类的静态成员**

```cpp
class A {
public:
    static int count;
};
int A::count = 0;
```

- 所有对象共享同一个 `count`。
    
- 不需要实例化对象就能访问：`A::count`
    

---

#### 📌 3. 在文件作用域（全局）：**限制符号访问**

```cpp
static int x = 5;  // 只能在当前 .cpp 文件中使用（内部链接）
```

用于**隐藏实现细节**，**防止命名冲突**。

---

## ✅ 三、`constexpr`：**编译期常量（C++11 起）**

### 🔹 作用：

- 表示变量、函数、构造函数可以在**编译期就被求值**。
    
- 可用于定义**常量表达式、数组大小、模板参数等**。
    

---

### 📌 示例：

```cpp
constexpr int square(int x) {
    return x * x;
}

constexpr int a = square(5);  // ✅ 编译期计算
```

### 📌 和 const 的区别：

|比较|`const`|`constexpr`|
|---|---|---|
|主要用途|只读变量|编译期常量/函数|
|求值时间|可以是运行时或编译期|一定要在编译期|
|支持函数吗|❌|✅ 支持|
|用作模板参数|有限制|✅ 可用|

---

## ✅ 四、组合使用的场景和含义

|写法|意义|
|---|---|
|`static const int a = 5;`|常量 + 静态变量，生命周期长，不能修改|
|`static constexpr int a = 5;`|编译期常量 + 静态变量，作用域限制|
|`const static int a = 5;`|顺序无影响，与上面等价|
|`constexpr const int a = 5;`|❌ 冗余，`constexpr` 本身已经包含 `const` 含义|
|`const int*`|指向常量的指针（指针能改，内容不能）|
|`int* const`|常量指针（内容能改，指针不能）|

---

## 🎯 实战演示：区别感受

```cpp
#include <iostream>

constexpr int square(int x) {
    return x * x;
}

int main() {
    const int a = 5;
    int arr1[a];  // ❌ 非标准（不保证支持 VLA）

    constexpr int b = 6;
    int arr2[b];  // ✅ 编译期常量，合法

    static const int c = 7;
    int arr3[c];  // ✅ 有些编译器接受

    static int counter = 0;  // 静态变量，在 main() 多次调用时保留值
}
```

---

## ✅ 五、总结对比

|关键字|含义|生命周期|求值时间|可否修饰函数|是否隐含 const|
|---|---|---|---|---|---|
|`const`|只读变量/只读接口|根据上下文|编译期或运行时|❌（不能修饰函数）|❌|
|`static`|静态存储 / 文件作用域|程序全程 / 文件内|不相关|✅（修饰成员函数）|❌|
|`constexpr`|编译期常量或函数|编译期|编译期|✅（函数、构造）|✅|

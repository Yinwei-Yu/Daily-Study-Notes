
## 一、为什么需要 `unique_ptr`？

在 C++ 中，我们通常用 `new` 和 `delete` 来手动管理内存，比如：

```cpp
int* ptr = new int(10);  // 分配内存
// 使用 ptr ...
delete ptr;              // 释放内存
```

但这样做有**两个大问题**：

1. **容易忘记 `delete`，造成内存泄漏**。
    
2. **程序异常中断时，可能来不及释放内存**。
    

为了解决这个问题，C++11 引入了 **智能指针**，自动管理资源，程序不再需要手动 `delete`。

---

## 二、什么是 `unique_ptr`？

`std::unique_ptr` 是一种 **独占式** 智能指针：

- 它**拥有**所指向的对象。
    
- **只有一个 `unique_ptr` 可以拥有这个对象**。
    
- 当 `unique_ptr` 被销毁时，它会自动调用 `delete`。
    

### 类比生活：

你买了一本书（对象），这本书只能属于你一个人（`unique_ptr`）。你不能同时把它“借”给别人用（不能复制）。如果你想让别人拥有这本书，**你必须把书“转交”出去**（用 `std::move()`）。

---

## 三、如何使用 `unique_ptr`？

### ✅ 创建 `unique_ptr`：

```cpp
#include <memory>  // 需要包含这个头文件

std::unique_ptr<int> p1 = std::make_unique<int>(42);  // 推荐用 make_unique
```

也可以这样写（不推荐）：

```cpp
std::unique_ptr<int> p2(new int(42));
```

### ✅ 自动释放资源：

```cpp
{
    std::unique_ptr<int> p = std::make_unique<int>(100);
    // 使用 p，比如 *p = 200;
} // 作用域结束，p 自动销毁，delete 自动调用
```

### ❌ 不允许复制：

```cpp
std::unique_ptr<int> p1 = std::make_unique<int>(10);
std::unique_ptr<int> p2 = p1; // ❌ 错误！不能复制 unique_ptr
```

### ✅ 允许移动：

```cpp
std::unique_ptr<int> p1 = std::make_unique<int>(10);
std::unique_ptr<int> p2 = std::move(p1);  // OK：p2 拥有对象，p1 变为空
```

---

## 四、`unique_ptr` 的常用操作

### 解引用：

```cpp
std::unique_ptr<int> p = std::make_unique<int>(5);
std::cout << *p << std::endl;  // 输出 5
```

### 访问成员（如果是指向对象）：

```cpp
struct Person {
    void sayHi() { std::cout << "Hi!" << std::endl; }
};

std::unique_ptr<Person> p = std::make_unique<Person>();
p->sayHi();  // 使用 -> 访问成员函数
```

### 重置（释放原来的对象）：

```cpp
p.reset();  // 手动释放内存，变为 nullptr
```

### 替换指向对象：

```cpp
p.reset(new int(20));  // 释放原来的，指向新对象
```

---

## 五、什么时候用 `unique_ptr`？

- 当你需要**局部管理资源**，只让一个指针拥有这个资源时。
    
- 常用于**工厂函数返回对象**：
    

```cpp
std::unique_ptr<Person> createPerson() {
    return std::make_unique<Person>();
}
```

---

## 六、`unique_ptr` vs 原始指针 vs 其他智能指针

|指针类型|自动释放内存|可以复制|适合场景|
|---|---|---|---|
|原始指针|❌ 否|✅ 是|不推荐新代码使用|
|`unique_ptr`|✅ 是|❌ 否|对象所有权唯一|
|`shared_ptr`|✅ 是|✅ 是|多个对象共享资源|

---

## 七、小结

- `unique_ptr` 是一种安全的替代原始指针的方式。
    
- 它确保对象只被一个指针拥有。
    
- 它在作用域结束时自动释放资源，避免内存泄漏。
    

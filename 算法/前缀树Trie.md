# 前缀树

**前缀树（Prefix Tree）**，又叫 **Trie 树**（发音为 "try"），是一种高效的用于处理字符串集合的**树形数据结构**，特别适合用来处理与**字符串前缀相关的问题**。

---

## 一、什么是前缀树（Trie）？

### 🌲 定义

前缀树是一种多叉树结构，每个节点代表一个字符。从根到某个节点的路径组成一个字符串，表示该字符串存在于字典中。

- 根节点不包含字符。
- 每个子节点包含一个字符。
- 每个节点的子节点字符各不相同。
- 如果某个节点是单词的结尾，通常会用一个标记（如 `is_end`）来表示。

---

## 二、前缀树的特点

| 特点 | 描述 |
|------|------|
| **高效前缀匹配** | 可快速查找所有具有某个前缀的字符串 |
| **节省空间** | 共享公共前缀的字符串可以共用路径，减少存储开销 |
| **插入和查询效率高** | 时间复杂度约为 O(L)，L 是字符串长度 |

---

## 三、前缀树的基本操作

1. **插入（Insert）**
   - 将字符串逐字符插入 Trie 树中。
   - 若字符对应的节点不存在，则新建节点。

2. **搜索（Search）**
   - 查找完整字符串是否存在于 Trie 中。
   - 需要判断最后节点是否是某个单词的结尾。

3. **前缀匹配（StartsWith）**
   - 判断是否存在以某个字符串为前缀的单词。
   - 不需要判断是否是单词结尾。

---

## 四、Python 实现 Trie 的简单示例

```python
class TrieNode:
    def __init__(self):
        self.children = {}  # 字符 -> 子节点映射
        self.is_end = False  # 是否是一个单词的结尾


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
        node.is_end = True

    def search(self, word: str) -> bool:
        node = self.root
        for char in word:
            if char not in node.children:
                return False
            node = node.children[char]
        return node.is_end

    def startsWith(self, prefix: str) -> bool:
        node = self.root
        for char in prefix:
            if char not in node.children:
                return False
            node = node.children[char]
        return True
```

---

## 五、应用场景

| 应用场景 | 描述 |
|----------|------|
| **自动补全（搜索引擎/输入法）** | 输入部分字符后，提示可能的完整词 |
| **拼写检查 / 纠错** | 快速查找相似拼写的词 |
| **IP 路由最长前缀匹配** | 在路由表中查找最长匹配的网络地址 |
| **词频统计** | 结合计数功能记录单词出现次数 |
| **字典查找 / 单词游戏（如 Scrabble）** | 快速判断某个单词是否存在 |

---

## 六、Trie 与其他结构的对比

| 数据结构 | 插入时间 | 查询时间 | 前缀查找 | 内存占用 |
|----------|-----------|-----------|------------|-------------|
| **哈希表** | O(L) | 平均 O(L) | ❌ 差 | ⬇️ 较低 |
| **平衡 BST** | O(L log N) | O(L log N) | ❌ 一般 | ⬇️ 中等 |
| **Trie** | O(L) | O(L) | ✅ 强大 | ⬆️ 较高 |

虽然 Trie 占用内存稍高，但它的**前缀查找性能非常优秀**，非常适合做搜索推荐、拼写检查等任务。

---

## 七、优化版本：压缩 Trie（Radix Tree / Patricia Trie）

为了节省空间，有些变种对 Trie 进行了压缩：

- 合并只有一个子节点的连续节点。
- 使用更紧凑的结构存储路径。
- 适用于大规模字符串集合。

---

## 八、总结一句话

> **前缀树（Trie）是一种利用字符串前缀共享路径的树形结构，能高效地进行字符串插入、查找和前缀匹配操作，常用于自动补全、拼写检查等场景。**

# 8.3 for-range 的配套用法

可以使用 `for` 循环读取 `map`：

```go
for key, value := range map1 {
	...
}
```

第一个返回值 `key` 是 `map` 中的 key 值，第二个返回值则是该 key 对应的 value 值；这两个都是仅 `for` 循环内部可见的局部变量。其中第一个返回值 `key` 值是一个可选元素。如果你只关心值，可以这么使用：

```go
for _, value := range map1 {
	...
}
```

如果只想获取 `key`，你可以这么使用：

```go
for key := range map1 {
	fmt.Printf("key is: %d\n", key)
}
```

示例 8.5 [maps_forrange.go](examples/chapter_8/maps_forrange.go)：

```go
package main
import "fmt"

func main() {
	map1 := make(map[int]float32)
	map1[1] = 1.0
	map1[2] = 2.0
	map1[3] = 3.0
	map1[4] = 4.0
	for key, value := range map1 {
		fmt.Printf("key is: %d - value is: %f\n", key, value)
	}
}
```

输出结果：

	key is: 3 - value is: 3.000000
	key is: 1 - value is: 1.000000
	key is: 4 - value is: 4.000000
	key is: 2 - value is: 2.000000

注意 `map` 不是按照 key 的顺序排列的，也不是按照 value 的序排列的。

> 译者注：map 的本质是散列表，而 map 的增长扩容会导致重新进行散列，这就可能使 map 的遍历结果在扩容前后变得不可靠，Go 设计者为了让大家不依赖遍历的顺序，每次遍历的起点--即起始 bucket 的位置不一样，即不让遍历都从某个固定的 bucket0 开始，所以即使未扩容时我们遍历出来的 map 也总是无序的。

问题 8.1： 下面这段代码的输出是什么？

```go
capitals := map[string] string {"France":"Paris", "Italy":"Rome", "Japan":"Tokyo" }
for key := range capitals {
	fmt.Println("Map item: Capital of", key, "is", capitals[key])
}
```

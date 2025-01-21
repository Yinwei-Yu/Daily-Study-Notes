# hashmap

```rust
use std::collections::HashMap;//引入HashMap
fn main() {
    //HashMap<K,V>方法存储一个键K和一个值V

    //创建一个空的HashMap,并插入元素
    let mut scores = HashMap::new();
    scores.insert(String::from("Blue"), 10);
    scores.insert(String::from("Yellow"), 50);

    //访问HashMap中的元素,get方法,返回一个Option<&V>
    let team_name = String::from("Blue");
    let score = scores.get(&team_name).copied().unwrap_or(0);//copied()方法返回一个拷贝的值,这样score获得的是一个值而不是一个引用,unwrap_or()方法在scores中没有对应的键的时候将其设置为0

    //遍历HashMap
    for(key,value) in &scores {
        println!("{key},{value}");
    }

    //hashmap的所有权
    let field_name = String::from("Favorite color");
    let field_value = String::from("Blue");
    let mut map = HashMap::new();
    map.insert(field_name, field_value);
    //println!("{},{}",field_name,field_value);//这里会报错,因为field_name和field_value的所有权已经转移给了map

    /*
    更新hashmap的值,有三种情况:
    1.覆盖一个值
    2.在键没有对应的值的时候插入键值对
    3.根据旧值更新一个值
     */

    //1.覆盖一个值,使用相同的键插入一个不同的值,则旧值会被覆盖掉
    scores.insert(String::from("Blue"), 25);

    //2.在键没有对应的值的时候插入键值对
    scores.entry(String::from("Yellow")).or_insert(50);//如果键Yellow对应的值不存在,则插入50,否则不做任何操作
    //entry以想要检查的键作为参数,返回一个枚举Entry,代表可能存在或者不存在的值,若键存在则返回相应的值
    //Entry有两个方法,or_insert和or_insert_with
    //or_insert方法在键对应的值不存在的时候插入一个值,or_insert_with方法在键对应的值不存在的时候根据一个闭包插入一个值

    //3.根据旧值更新一个值
    let text="hello world wonderful world";
    let mut map=HashMap::new();
    for word in text.split_whitespace() {
        let count = map.entry(word).or_insert(0);
        *count+=1;
    }
    /*
    这段代码统计了text中每个单词出现的次数
    split_whitespace()方法将text根据空格分割成单词
     */
    println!("{:?}",map);
}
```

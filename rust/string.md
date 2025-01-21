# string

代码演示

```rust
fn main() {
    /*
    1. rust核心语言中只有一种字符串类型：str，字符串切片
    2. 字符串切片是对字符串的引用，它的类型是&str,它指向一个存在内存某处的字符串
    3. 字符串的字面值也是字符串切片
    4. String类型是标准库提供的，它是可变的，可增长的，可拥有的，它是UTF-8编码的字符串类型
    */

    //新建字符串
    let mut s = String::new();

    //使用to_string方法
    let data = "initial contents";
    let s = data.to_string();

    //直接用于字符串字面值
    let s = "initial contents".to_string();

    //使用String::from方法
    let s = String::from("initial contents");

    //字符串是UTF-8编码的，所以可以包含任何字节
    let hello = String::from("你好");
    let happy = String::from("😊");
    let dog = String::from("dog");

    //更新字符串
    //1. push_str方法,附加字符串slice
    let mut s = String::from("foo");
    s.push_str("bar");

    //2. push方法,附加单个字符
    let mut s = String::from("foo");
    s.push('l');

    //连接字符串
    //1. 使用+运算符
    let s1 = String::from("hello,");
    let s2 = String::from("world!");
    let s3 = s1 + &s2; //注意这里s1的所有权被转移,不能再使用
                       /*
                       这里+运算符调用了add方法，add方法的签名如下：
                       fn add(self, s: &str) -> String
                       这里s1是String类型，s2是&String类型，所以s2的所有权没有被转移
                       但是&s2是&String类型，而add方法需要的是&str类型，所以&s2会被隐式转换为&str类型,也就是说&String可以隐式转换为&str
                        */

    //2. 使用format!宏
    let s1 = String::from("tic");
    let s2 = String::from("tac");
    let s3 = String::from("toe");
    let s = format!("{}-{}-{}", s1, s2, s3); //format!宏不会获取任何参数的所有权

    //索引字符串
    //字符串是一个Vec<u8>的封装，但是不可以使用索引来访问字符串,因为一个字符可能占用多个字节
    //rust中有三种理解字符串的方法：字节、标量值、字形簇
    /*
    1. 字节：字符串的底层字节表示
    2. 标量值：rust中的char类型，它是Unicode标量值的表现形式,其中可能包含某些没有显示的字符,比如发音符号等
    3. 字形簇：人们认为的字符
     */

    //切片字符串
    let hello = "Здравствуйте";
    let s = &hello[0..4]; //会得到前四个字节,即前两个字符
                          //如果使用&hello[0..3]会panic

    //遍历字符串
    for c in "नमस्ते".chars() {
        println!("{}", c);
    } //这里会打印出每个字符

    for b in "नमस्ते".bytes() {
        println!("{}", b);
    } //这里会打印出每个字节

    //对于字型簇,需要使用第三方库
}

```

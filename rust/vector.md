# vector

代码使用如下:

```rust
fn main () {
    //vec的两种声明方式
    let v1:Vec<i32> = Vec::new();//new方法需要指明类型
    let v2=vec![1,2,3];//使用宏vec!声明,编译器会自动推断类型

    let mut v3=Vec::new();//声明可变vec,push元素后编译器可以推断类型
    v3.push(5);
    v3.push(6);
    v3.push(7);

    let third: &i32 = &v3[2];//获取元素,获取的是引用,如果越界会panic
    println!("The third element is {}", third);

    match v3.get(2) {//获取元素,获取的是Option<&T>,如果越界会返回None,不会panic
        Some(third) => println!("The third element is {}", third),
        None => println!("There is no third element."),
    }

    //遍历vec
    for i in &v3 {
        println!("{}", i);
    }
    //可变vec遍历
    for i in &mut v3 {
        *i += 50;
    }

    //vec中的元素可以是枚举,这样vec中的元素类型可以不同
    enum SpreadsheetCell {
        Int(i32),
        Float(f64),
        Text(String),
        
    }
    let row = vec![
        SpreadsheetCell::Int(3),
        SpreadsheetCell::Text(String::from("blue")),
        SpreadsheetCell::Float(10.12),
    ];
    
    /*注意
    1. vec<T>是一个泛型类型,所以vec中的元素类型必须相同
    2. vec<T>是一个堆分配的数据结构,所以vec<T>的元素在堆上分配
    3. vec<T>的元素在vec<T>的生命周期内有效,vec<T>的生命周期结束,vec<T>的元素也会被释放
    4.因为对一个变量的引用不能同时有可变和不可变的引用,所以vec<T>的元素不能同时有可变和不可变的引用
    */ 
}
```

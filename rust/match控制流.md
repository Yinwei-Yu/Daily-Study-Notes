# match

基本使用:

```rust
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter,
}

fn value_in_cents(coin:Coin)-u8 {
    match coin {
        Coin::Penny => {
            prinrtln!("Lucky penny!");
            1//match中的分支可以加大括号
        }
        Coin::Nickel => 5,//可做返回值
        Coin::Dime => 10,
        Coin::Quarter => 25,
    }
}
```

也可以给匹配到的值进行绑定:

```rust
#[derive(Debug)]
enum UsState {
    Alabama,
    Alaska,
    // --snip--
}

enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter(UsState),
}

fn value_in_cents(coin:Coin)->u8 {
    match coin {
        Coin::Penny => 1,
        Coin::Nickel =>5,
        Coin::Dime => 10,
        Coin::Quarter(state)=>{//这里的state是一个UsState类型的变量，匹配之后就会把coin里面的UsState类型的值赋值给state
            println!("State quarter from {:?}!", state);
            25
        }
    }
}
```

## 利用match匹配Option

```rust
fn main() {
    fn plus_one(x:Option<i32>) -> Option<i32> {
        match x {
            None => None,
            Some(i) => Some(i + 1),
        }//match要匹配所有可能的情况
    }

    let five = Some(5);
    let six = plus_one(five);
    let none = plus_one(None);
}
```

## 使用通配模式和_占位符

```rust
let dice_roll = 9;
match dice_roll {
    3 => add_fancy_hat(),
    7 => remove_fancy_hat(),
    other => move_player(other),//这里把其他的值绑定到other上然后使用这个变量
}

fn add_fancy_hat() {}
fn remove_fancy_hat() {}
fn move_player(num_spaces: u8) {}
```

在不希望使用其他的值(上述绑定到other的值)时,可以使用_占位符like this:

```rust
let dice_roll = 9;
match dice_roll {
    3 => add_fancy_hat(),
    7 => remove_fancy_hat(),
	_=>(),
}

fn add_fancy_hat() {}
fn remove_fancy_hat() {}
fn move_player(num_spaces: u8) {}
```

## if let语法

if let语法用于match中只需要匹配其中一种模式而忽略其他模式的情况,如下:

原来的代码是:

```rust
    let config_max = Some(3u8);
    match config_max {
        Some(max) => println!("The maximum is configured to be {max}"),
        _ => (),
    }
```

可以改写为:

```rust
    let config_max = Some(3u8);
    if let Some(max) = config_max {//注意这里是一个等号而非两个
        println!("The maximum is configured to be {max}");
    }

```

if let忽略了穷尽检查,可以认为是match的一个语法糖,也可以和else配合使用:

```rust
    let mut count = 0;
    if let Coin::Quarter(state) = coin {
        println!("State quarter from {state:?}!");
    } else {
        count += 1;
    }
```
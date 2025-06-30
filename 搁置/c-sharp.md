# 枚举

整形常量的集合
用于表示状态,类型等
可以让使用整数的判断变得可读

```c#
enum E_myEnum{
	RED,// 0
	BLUE,// 1
	BLACK=3,// 3
	YELLOW// 4
}

```

枚举在==namespace==中声明
或者类和结构体中(不常用)
==枚举不能在函数中声明==

```c#
enum E_MonsterType
{
	Normal,
	Boss
}

public void Main(string args[])
{
	E_MonsterType monsteType = E_MonsterType.Normal;
	if(monsterType==E_MonsterType.Boss)
		return;
	else
		return;
}

// 常用于switch

```

枚举可以和int互转
```c#
int i = (int)monsterType;
monsterType=0;
```

枚举和string互转
```c#
string str= monsterType.ToString();//str = name of enum variable
//str = Normal
monsterType = (E_Monster)Enum.Parse(typeof(E_Monster),"Boss");//Enum类型使用Parse函数进行类型转换,需要搭配强转

```

# 数组

数组声明方式:
```c#
int [] arr1;
int [] arr2 = new int [5];
int [] arr3 = new int [5] {1,2,3,4,5};//长度指定
int [] arr4 = new int [] {1,2,3,4,5};//长度不指定,由大括号中元素个数决定
int [] arr5 = {1,2,3,4,5};
```

数组访问
```c#
int[] arr = new int[5]{1,2,3,4,5};

int len=arr.Length;

int num=arr[2];

for(int i=0;i<arr.Length;i++)
	Console.Write(arr[i]);
	
foreach(num in arr)
	Console.Write(num);
```

增删查改

# 二维数组

```c#
int[,] arr1;
int[,] arr2=new int [2,3];
int[,] arr3=new int [2,3]{{1,1,1},
						  {2,2,2},
						  {3,3,3}};
int[,] arr4=new int[,]{{1,1,1},
					   {2,2,2},
					   {3,3,3}};
int[,] arr5 = {{1,1,1},
				{2,2,2},
				{3,3,3}};
				
```

获得长度
```c#
arr.GetLength(0);//行数
arr.GetLength(1);//列数
```

访问数组元素

```c#
int a = arr[1,2];

```

#  交错数组

每个数组元素中元素个数可以不同
```c#
int[][] arr1;
int[][] arr2 = new int[3][];
int[][] arr3 = new int[3][]{ new int[] {1,2,3},
							new int[] {4,5,6},
							new int[] {7,8}};
int[][] arr4 = new int[][]{new int[] {1,2,3},
							new int[] {4,5,6},
							new int[] {7,8}};

int[][] arr5= {new int[] {1,2,3},
				new int[] {4,5,6},
				new int[] {7,8}};

```

长度,访问,遍历
```c#
int row=arr.GetLength(0);
int col1=arr[0].Length;

int num = arr[0][1];

for(int i=0;i<arr.GetLength(0);i++)
	for(int j=0;j<arr[i].length;j++)
		Console.Write(arr[i][j]);
```

# 值和引用类型

引用类型:string,数组,类
值类型:基本类型,结构体,枚举
==本质区别==在于存储位置是stack还是heap

```c#
int a = 10;
int arr = new int[]{1,2,3,4};

int b = a;//复制
int arr2 = arr;//arr2指向arr内存相同位置

//值类型存储在stack
//引用类型存储在heap
```

## string

自动使用new分配地址
更改字符串后原有内存空间不释放

# 函数

定义在class和struct中

## ref和out

```c#
static void ChangeValue(ref int num)
{
	num = 0;
}
//ref修饰值类型变量相当于cpp &引用运算符
int a=1;
ChangeValue(ref a);

//out和ref很相似
/*
区别在于:
	1.ref传入的变量必须初始化,out不用
	2.out传入的变量必须在内部赋值,ref不用
*/
```
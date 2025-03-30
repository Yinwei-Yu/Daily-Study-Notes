# 关系型数据库
primary key
可以有多个，唯一标识

# 基本命令
CREATE \`database\`; 创建
SHOW DATABASES;  展示所有数据库
DORP DATABASE \`database\` 删除某个数据库
USE \`database\`

# 数据类型
INT   整数
DECIMAL(m,n)  小数，m：总位数，n：小数位数
VARCHAR(n)  字符串 n：字符串的长度
BLOB   Binary Large Object  图片 视频等
DATE  日期 ‘YYYY-MM-DD’ 2024-07-09
TIMESTAMP 时间 ‘YYYY-MM-DD HH:MM:SS’

# 数据表格
```mysql
CREATE TABLE student(
	`student`_id INT PRIMARY KEY,
	
	`name` VARCHAR(10),

	`major` VARCHAR(20),

	PRIMARY KEY(student_id);-- 也可以作为指定PRIMARY /KEY 的方式
);
```

DESCRIBE \`student\`;   显示表格信息
删除表格和删除数据库相同

- 新增属性
```mysql
ALTER TABLE `student` ADD gpa DECIMAL(3,2);
```
- 删除属性
```mysql
ALTER TABLE `student` DROP COLUMN gpa;
```

# 存储资料

```mysql
INSERT INTO `student` VALUES(1，'小白'，'历史'，3.99)-- 顺序和创建表格时相同 
```

```mysql
SELECT * FROM `student`;-- 搜索资料
```

NULL表示空

```mysql
INSERT INTO `student`(`name`,`major`,`student_id`) VALUES ('小白'，`数学`,4)
```
（\`name\`。。。）表示参数填入的顺序


# 限制 约束

```mysql
CREATE TABLE student(
	`student`_id INT AUTO_INCREMENT,-- 自动补全
	
	`name` VARCHAR(10) NOT NULL ,-- 非空
	
	`major` VARCHAR(20) UNIQUE,-- 唯一

	`gpa` INT DEFAULT 0-- 默认值

	PRIMARY KEY(`student_id`)
		
```

# 修改 删除资料

```mysql
SET SQL_SAFE_UPDATES = 0;//关闭默认更新行为

UPDATE `student`-- 表格名
SET `major` = '英语文学'
WHERE `major` = '英语'; -- 表条件判断,注意分号位置
```

OR AND等判断符

```mysql
DELETE FROM `student` 
WHERE `student_id` = 4;
```

<> :不等号

# 查找资料

```mysql

SELECT `name` FROM `student`;

SELECT * FROM `student`;

SELECT `name`,`major` FROM `student`

SELECT * FROM `student` ORDER BY `score` DESC;
-- 查找时同时排序,/ORDER /BY 后跟排序标准 /DESC表示是否反序,ORDER BY后可加多个属性,在第一个属性相同时,按照后一个属性相同

SELECT * FROM `student` LIMIT 3;-- 限制返回资料的数量
```


# 聚合函数

~~~mysql
SELECT COUNT(*) FROM `student`;
-- 可综合使用条件判断

SELECT AVG(`Total_Grades`) FROM `student`; -- 平均值

SELECT SUM(`Grades`) FROM `student`; -- 总和

SELECT MAX(`Total_Grades`) FROM `student`; -- 最大值

-- MIN 最小值

~~~

# wildcards 万用字符

~~~mysql
-- %多个字符
-- _代表一个字符

SELECT *
FROM `client`
WHERE `phone` LIKE `%335`;-- LIKE 用于模糊查询

~~~


# union 并集

~~~mysql

SELECT `name`
FROM `employ`
UNION
SELECT `client_name`
FROM `client`
UNION
SELECT `branch_name`
FROM `branch`
-- 有限制 ,每个分支内的属性数目应该相同
-- 类型应该相同

~~~


# join 连接

连接两张表

~~~mysql

SELECT *
FROM `emoloyee` LEFT JOIN `branch`-- LEFT ->无论判断条件是否成立,都会返回所有数据
-- Right同理反向
ON `employee`.`emp_id` = `manager_id`;-- 条件判断
-- `employee`.`emp_id`,表示后者属于前者


~~~

# 子查询 subquery

嵌套查询


~~~mysql

SELECT `name`
FROM `employee`
WHERE `emp_id` =(

	SELECT `manager_id`
	FROM `branch`
	WHERE `branch_name`=`研发`
);

~~~

# on delete

~~~mysql

ON DELETE SET NULL;-- 外键被删除后本表内对应值为NULL

ON DELETE CASCADE;-- 外键删除后本表内对应值删除

-- 如果一个属性既是primary key 又是 foreign key,则不能设置为SEY NULL

~~~
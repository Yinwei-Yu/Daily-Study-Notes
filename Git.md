[[progit.pdf]]
# Git config
```
git config --list  #显示环境配置
```

# Git基础

## 状态浏览
```
git status -s
git status --short
```

?? 新添加未跟踪
M 修改过的文件
当有两列时,左侧代表暂存区状态,右侧代表工作区状态
A 新添加到暂存区中的文件

## 忽略文件

在.gitignore中
可以使用正则表达式


## 查看不同
```
git diff #比较工作区与暂存区快照之间的差异

git diff --staged #比较已暂存的与最后一次提交的文件差异

git diff --cached #与上述相同
```


## 跳过暂存区

```
git commit -a #跳过暂存区,将所有已经跟踪过的文件暂存起来并一并提交
```

## 移除文件

必须从暂存区域移除然后提交
使用 git rm命令,会连带删除工作目录中的文件

如果简单地手动从工作目录中删除,会显示为跟踪删除记录

```
git rm --cached README
# 从git仓库中删除但是不从工作区删除
# 适用于不小心添加了不必要的配置文件
-f #强制选项,在删除之前修改过或已经存放到暂存区的文件时需要使用
```

## git重命名

```
git mv file_from file_to
```

实际上执行了三个操作

```
mv file_from file_to
git rm file_from
git add file_to
```

- 在使用其他的重命名工具,记得使用git rm 和git add

## git 提交历史

```
git log
```

常用选项:
```
-p/-patch #显示每次提交引入的差异(以补丁的格式输出)
-n #显示n条最近的提交
--stat # 文件修改统计信息
--shortstat # 只显示--stat中最后行数修改添加移除统计
--pretty # 使用内建选项展示提交历史
--pretty=oneline/short(只包括提交信息,作者信息)/full（包括提交信息，作者和提交者信息）/fuller(full加上日期)
```

***

format选项可以定制显示格式，对后期提取分析有用
```
git log --pretty=format:"%h - %an,%ar : %s"
```
选项如下：

| 选项  | 说明                   |
| --- | -------------------- |
| %H  | 提交完整的哈希值             |
| %h  | 简写哈希值                |
| %an | 作者名字                 |
| %ae | 作者电子邮件               |
| %ad | 作者修订日期，使用--date=选项定制 |
| %ar | 作者的修订日期，多久之前         |
| %s  | 提交说明                 |

将ae等选项的a换成c即为提交者

```
--graph # 图形化显示
```

- 与上述--pretty结合使用

***

限制输出长度：

```
--since
--until

git log --since=2.weeks

```

***

限制路径：
- 在只关心某些文件或目录的修改时使用
在git log的最后指定他们的路径，使用两个短划线隔开之前的选项和后面限定的路径名

***

其他有用的选项：

![[Pasted image 20240722215737.png]]

为了避免显示合并提交的提交，可以使用```--no-merges ```选项

## 撤销操作

***撤销操作可能会导致之前的工作丢失***

```
git commit --amend
```

会将暂存区内的文件提交，如果自上次提交以来还没有任何修改，则则个命令修改的只是提交信息

若忘记暂存某些需要的修改，则可以下面这样
```
git commit -m "initial commit"
git add forgotten_file
git commit --amend
```
 there will be only the last commit message

**最有价值的是稍微修改最后的提交**

- 其实更像是用一个全新的提交,因为旧的提交信息完全不存在

***

取消暂存的文件:
```
git reset HEAD file
```

***

撤销对文件的修改

```
git checkout -- file
```

会使用上次提交的信息

## 远程仓库

```
git remote -v #查看远程仓库名称和url
```

## 打标签

```
git tag #列出所有存在的标签
```

-l以字母顺序列出标签

可以指定特定模式(正则表达式

***

**创建标签**

轻量标签
附注标签

## git别名

```
git config --global alias.co checkout
```

# 分支

## 分支管理

```
git branch -v #查看所有分支最后一次提交
git branch -r #列出所有远程分支

```

## 远程分支

## 变基

```
rebase命令可以将一个分支上的所有修改提交到另一分支上

such as:
git checkout experiment
git rebase main

```
如下图:
![[Pasted image 20240723145412.png]]
然后到main分支上,进行一次快速合并
```
git checkout main
git merge experiment
```

原理是:rebase命令执行时,先找到当前分支和目标基底分支的最近共同祖先分支(c2),然后根据当前分支的历次提交的修改信息与祖先的不同并创建一个临时文件,然后将分支指向基底分支c3,根据临时文件进行一次fast-forword合并,创建一个新的snapshot

这样做的结果与merge没有任何区别,只是提交信息更加简洁,不用有合并的提交信息

# 服务器git

# 分布式git

# github

# git工具
## 选择修订版本


使用git show查看某一分支的顶端

```
git show branch1
git show 59vd464
```

直接加名字和使用sha-1值效果相同

```
git rev-parse branch1
```

显示指定分支的sha-1值

***

使用
```
git reflog # 查看引用日志,每当head位置变化时该日志就会记录
```

可以使用
```
git log -g #查看类似git log的日志引用信息
```

**引用日志仅存在于本地仓库**

***

祖先引用

使用'^'查看父引用

git ~相同
都是第一父提交

git HEAD^2 会使用第二父提交(仅在合并后有用

***
双点

```
git log master..experiment
```

可以查看在experiment中但不在master中的提交

可以使用origin/master来查看当前分支和远程分支的关系

![[QQ_1722003041897.png]]

***

三点

```
git log A...B #查看在两个分支中的一个有而不被同时包含的
```

参数--left-right,指出某个提交具体属于哪个分支

## 交互式暂存

```
git add -i
```

## 贮藏与清理

## 签署工作

## 搜索

```
git grep <pattern> # 搜索当前目录文件中包含pattern字符串的位置
```

```
git grep -n/--line-number str #显示包含str的行数
git grep -c/--count str # 仅显示包含str的文件以及多少处包含了目标str
git grep -p/--show-function str #显示字符串上下文信息
```

***

git 日志搜索

```
git log -S str #搜索新增和删除该字符串的提交
```

使用-G选项进行正则表达式搜索

***

git函数log搜索

```
git log -L :funcname:filename
git log -L start,end:filename
```

## 高级合并

```
git merge -Xignore-space-change filename # -Xignore-space-change选项忽略合并文件中的空白符
```

## Rerere

## git调试

```
git blame -L begin,end filename #查看指定文件指定行数内的修改者及修改时间
```
 ***
二分查找

## 子模块

## git打包

## 替换

# 自定义git

## 配置Git

[配置git](https://git-scm.com/book/zh/v2/%E8%87%AA%E5%AE%9A%E4%B9%89-Git-%E9%85%8D%E7%BD%AE-Git)

## git属性

### 导出版本库

### 合并策略

## git钩子




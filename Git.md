# Git config
```
git config --list  #显示环境配置
```

# Git基础

## 状态浏览
```
git status -s
git status -short
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

```


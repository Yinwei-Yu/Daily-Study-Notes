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

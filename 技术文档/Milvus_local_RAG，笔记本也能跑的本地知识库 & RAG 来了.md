> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [mp.weixin.qq.com](https://mp.weixin.qq.com/s/st4fALlUVRr9Blmlo82ilw)

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh4ek0GU1Snpd5xyiahZAUvz7OBtgIlbTu6EPpyEfG6V2kESjkXKK3fSHP6voC8hyCcn0RXhCkB2dicg/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1)

![](https://mmbiz.qpic.cn/mmbiz_jpg/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8Lcubib31OG3bmqdzrKX34RhBlCz3ucwgjEW75xCTicqWwWsDn43J2P3Q/640?wx_fmt=jpeg)

多数前端开发工程师可能都面临这样一个困境：每天需要查阅大量技术文档、项目规范和学习资料。传统的文件夹分类和搜索方式效率低下，经常为了找一个 API 用法翻遍整个项目文档。

一些大公司，可能会采用企业级知识库方案，通过智能问答来解决这个问题。但问题是：

1、不是所有公司都有这个预算

2、个人部署一套企业级知识库，环境配置复杂、学习门槛高，对新手极不友好

3、使用企业级知识库的平替，在线服务又会出现数据隐私泄露风险。

当然，以上问题不止是前端会遇到，所有有复杂文档管理、检索需求的朋友，其实都会遇到。

那怎么解决？这篇 “milvus_local_rag” 指南正是为你准备的。

这个轻量级 RAG 方案，用一台普通笔记本就能搭建起个人知识库，查询响应时间，也可以从几分钟缩短到几秒钟。

（备注：本项目是基于 Shubham Saboo 作者开源的 awesome-llm-apps 项目二次开发完成的。）

01
==

核心概念解释
======

在开始之前，让我们先了解几个关键概念，这样后续的操作会更加清晰：

**RAG（检索增强生成）：**简单来说，就是让 AI 在回答问题时，先从你的文档库中找到相关信息，再基于这些信息给出答案。就像考试时可以 "翻书" 一样，让 AI 的回答更准确、更有依据。

**向量数据库：**把文档转换成数字形式存储的 "智能仓库"。它能理解文档的含义，当你提问时，能快速找到最相关的内容片段。

**嵌入模型：**负责把文字转换成数字的 "翻译官"。它能理解文字的语义，让计算机也能 "读懂" 文档内容。

02
==

RAG 工作原理：从文档到智能问答的完整流程
======================

了解了基本概念后，让我们看看整个系统是如何工作的：

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8zl9SISa6iaSwyNeSDovNr9dR9rnpBloW9vNMDaIxVIj5HA8z5FFXOiaQ/640?wx_fmt=png&from=appmsg)

这个流程确保了 AI 的回答既基于你的文档内容，又具备良好的理解能力。

03
==

为什么选择轻量级方案？
===========

这是一个专为个人用户设计的轻量级 RAG 项目，核心思路是用最少的依赖实现最完整的功能。本文作者对 awesome-llm-apps 项目源代码进行了调整，整个系统只需要 Ollama 和 Qdrant 两个组件，一条命令就能启动完整的本地知识库。

**核心特点**：

*   **真正的本地化**：支持 Qwen、Gemma 等多种本地模型，数据完全不出本地
    
*   **极简部署**：无需复杂环境配置，Docker 一键启动向量数据库
    
*   **智能检索**：文档相似度搜索 + 网络搜索双重保障，确保答案质量
    
*   **灵活切换**：可在纯 RAG 模式和直接对话模式间自由切换
    

**实际价值**：让你用最小的成本获得企业级 RAG 能力，适合处理个人文档、学习资料或项目知识库，既保护隐私又提供智能问答体验。

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8xYDl2vjia95wV0p8IYiaWXyVBDt50Q9qkmK0sQsIeDDpo1lmicSUpWpXA/640?wx_fmt=png&from=appmsg)

04
==

实践部署
====

**（1）环境准备要求**

本教程不含 Python3、Conda 以及 Ollama 安装展示，请自行按照官方手册进行配置。

相关官网链接：

*   Python3 官网：https://www.python.org/
    
*   Conda 官网：https://www.anaconda.com/
    
*   Milvus 官网：https://milvus.io/docs/prerequisite-docker.md
    
*   Ollama 官网：https://ollama.com
    
*   Docker 官网：https://www.docker.com/
    

**（2）系统环境配置表**

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8jFhYenyYI9cZ3I9X80TztOldsYmpl0W7So4iccQhe88Jz1M4aUiaSuOQ/640?wx_fmt=png&from=appmsg)

#### （3）Milvus 向量数据库部署

 **Milvus 简介**

Milvus 是由 Zilliz 开发的全球首款开源向量数据库产品，能够处理数百万乃至数十亿级的向量数据，在 Github 获得 3 万 + star 数量。基于开源 Milvus，Zilliz 还构建了商业化向量数据库产品 Zilliz Cloud，这是一款全托管的向量数据库服务，通过采用云原生设计理念，在易用性、成本效益和安全性上实现了全面提升。

**部署环境要求**

必要条件：

*   软件要求：docker、docker-compose
    
*   CPU：8 核
    
*   内存：至少 16GB
    
*   硬盘：至少 100GB
    

**下载部署文件**

```
wget https://github.com/milvus-io/milvus/releases/download/v2.5.12/milvus-standalone-docker-compose.yml -O docker-compose.yml

```

**启动 Milvus 服务**

```
docker-compose up -d
```

```
docker-compose ps -a
```

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8O0YGJx9PjbnrYdNdKWyic5pteczRIvfntQYOrDDGRlghqgwHVcvtOjg/640?wx_fmt=png&from=appmsg)

#### （4） 模型下载与配置

**下载大语言模型**

```
# 下载Qwen3模型
ollama pull qwen3:1.7b
```

 **下载嵌入模型**

```
# 下载embedding模型
ollama pull snowflake-arctic-embed
```

**验证模型安装**

```
# 查看已安装模型列表
ollama list

```

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8SNTlZj0WzqDIPa35ibo5Iia596zAFS4VjwhqbMzBdprXjGic2o1gDsTSw/640?wx_fmt=png&from=appmsg)

#### （5）Python 环境配置

**创建虚拟环境**

```
# 创建conda虚拟环境
conda create -n milvus
# 激活虚拟环境
conda activate milvus

```

```
# 克隆项目代码
git clone https://github.com/yinmin2020/milvus_local_rag.git

```

**依赖包安装**

```
# 安装项目依赖
pip3 install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple/

```

**参数配置说明**

关键配置参数：

```
COLLECTION_NAME：自定义集合名称（必须配置）
"uri": "tcp://192.168.7.147:19530"：Milvus连接地址（必须修改为实际地址）

```

```
# 启动Streamlit应用
streamlit run release.py

```

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8azUz4F4E1vXjDYE0t5R6gQ6iawtodoYXdpialHFpyNHKibzVn3zUKVW4g/640?wx_fmt=png&from=appmsg)

#### （7）功能测试与验证

**访问应用界面**

应用启动后会自动跳转到 Web 界面，通常地址

```
http://localhost:8501

```

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8RZVRLfrbPiaictZFEVMyBpsO0ftcmoEAlWsgXEy4szgOOibdM43wp4AqA/640?wx_fmt=png&from=appmsg)

**文档上传测试**

1.  在 Web 界面中选择文档上传功能
    
2.  上传测试 PDF 文档（建议使用 Milvus 相关介绍文档）
    
3.  等待文档处理完成
    

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8w9IaBSOnicOhtufXbYKDwQPYPCKibh0HI8wGRKzYBqsOUeaeBjKib5kHg/640?wx_fmt=png&from=appmsg)

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8IedLcdiaEia5zf7CCsAA2rjJnsJnrUEQ4qIAuRQWK2MeSu2r1FU96g0Q/640?wx_fmt=png&from=appmsg)

**RAG 功能验证**

测试查询示例：

```
milvus向量查询能力有哪些？

```

通过此查询可以验证：

*   向量数据库检索功能
    
*   RAG（检索增强生成）能力
    
*   问答系统的准确性
    

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8XsW6DhkHlNnbPb26ZtA4RlL6cfF3nyapibpPCXjlXAQgRVxFdrjBnZA/640?wx_fmt=png&from=appmsg)

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8GNibrlxjw7TSbKib1whVKeGzcr9ibhN5MTMFpha1dUmyb61sc0lntL8Lg/640?wx_fmt=png&from=appmsg)

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8qkVUfxt6cbsvgOT8KMA41EgA3GdtwIO2Erliciah87ySgIqlMiaiaTiaMRw/640?wx_fmt=png&from=appmsg)

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia8kaw2qcmtxjiaeRsKlf67Au6ibBu6ZvylTqceRklm7VNibmJ41O55mY5ww/640?wx_fmt=png&from=appmsg)

05
==

写在最后
====

回望文章开头提到的那些令人望而却步的部署障碍：做 RAG 为什么要让简单的事情变得复杂？

其实，企业级知识库流行的同时，轻量级 RAG 也逐渐成为了个人侧的主流趋势。

轻量级 RAG 最大的价值在于各种成本低，能解决的问题很实在。几行代码就能让文档 "活" 起来，能问能答，而且简单好用，是很多中小企业或者个人用户入门 RAG 的第一步。

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh7X0LLYMYoSYVWUYDXwHXia81Pz366jVYh0Jos9DZicM6MzzG2PAE5pufiaB0oibmNo1JnhWgibENOiciaJw/640?wx_fmt=png&from=appmsg)

如对本教程仍有不理解的地方，欢迎扫描文末二维码交流。

作者介绍

![](https://mmbiz.qpic.cn/mmbiz_jpg/MqgA8Ylgeh4UzkeAkCmkn4gibkkyfr3Xy7NbZ8Bnt1e5Cf4OAiaqf6d8ZHZGf8k9BrEMzSxPXXicJcWHFUAHfS9fg/640?wx_fmt=other&from=appmsg&wxfrom=5&wx_lazy=1&wx_co=1&tp=webp)

Zilliz 黄金写手：尹珉

推荐阅读

[不止语义检索，Milvus+LangChain 全文检索 RAG 教程来了](https://mp.weixin.qq.com/s?__biz=MzUzMDI5OTA5NQ==&mid=2247508852&idx=1&sn=dee4222afddaa7e734001c90888d7a18&scene=21#wechat_redirect)

[用 Milvus 构建 RAG 系统，N8N VS dify 如何选？](https://mp.weixin.qq.com/s?__biz=MzUzMDI5OTA5NQ==&mid=2247509157&idx=1&sn=78edccc7618e4b1c4731d113bc6138f0&scene=21#wechat_redirect)

[金融支付 × 实时推荐：Milvus 如何支撑全球百亿交易的 “猜你喜欢”](https://mp.weixin.qq.com/s?__biz=MzUzMDI5OTA5NQ==&mid=2247509091&idx=1&sn=e28842e976a43171d2c7af0cb1c0349a&scene=21#wechat_redirect)

点击 “阅读原文” 即可体验 zillz cloud

![](https://mmbiz.qpic.cn/mmbiz_png/MqgA8Ylgeh4wRYzJPkrxOP2wMciaLmG7mJnLWypF9CfPhaKvgXMmfNEYoUTk7pJJ9Niak6qAqZSjsSk6mpDB5Hyg/640?wx_fmt=other&from=appmsg&wxfrom=5&wx_lazy=1&wx_co=1&tp=webp)
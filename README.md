#### Go基于有向无环图的并发执行流

参考文章：[Go基于有向无环图的并发执行流的实现-地鼠文档 (topgoer.cn)](https://cc.topgoer.cn/blog-230.html)

基于该文章进行编写，修复文章部分代码错误。



```go
Run(i interface{})
```

中的`i`可以为任意数据，用于传递数据给执行函数。



有向无环图如下：

![](https://cc.topgoer.cn/uploads/blog/202208/attach_17094f3a271a8e2e.png)
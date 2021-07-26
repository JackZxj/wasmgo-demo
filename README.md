# Go WebAssembly 尝试

本项目是一个实验性测试，目标是使用 webassembly 实现一个纯前端的 ssh 客户端，结果是： 即使是 webassembly 也不能逃脱浏览器的限制，不能随意访问任意地址。

但另一方面，作为第一次尝试使用 webassembly 来说，却也是值得记录，方便今后有用到时参阅。

**使用方法:**

```bash
# 注: 需要较高版本的 go 才支持 wasm

# 开启服务
$ make serve
# 浏览器访问 IP:8080
# 浏览器 F12 打开控制台, 编译生成的wasm文件很大,所以要等一会儿页面才能响应操作

# 修改 web/gowasm 中的代码后重新生成 wasm 文件
$ make wasm
```

> ref:
> https://geektutu.com/post/quick-go-wasm.html

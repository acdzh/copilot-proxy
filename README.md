该项目目的是使用 CloudFlare worker 加速 Copilot 的 api 接口 `https://copilot-proxy.githubusercontent.com/*`.

server 部分为部署在云上的代理脚本, 目前只有一个 CloudFlare worker. 可以自行部署.

client 作用是替换本地 ide / 编辑器的 copilot 插件中的 api 地址. 二进制可以在 release 下载, 直接运行是替换, `./copilot-proxy recover` 是恢复. 因为插件会更新, 插件每次更新后都要重新运行次程序. 目前仅支持 VSCode 与 JetBrains 的 Copilot 插件. 如果你的代理 server 是自行部署的话, 请自己修改代码中的 api 地址, 自行编译.

有问题请看源码, 谢谢茄子🍆.
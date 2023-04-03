该项目目的是使用自建服务加速 Copilot 的 api 接口, 即把 `https://copilot-proxy.githubusercontent.com/*` 替换为你自己的服务接口.

server 部分需要自行部署.

client 作用是替换本地 ide / 编辑器的 copilot 插件中的 api 地址. 二进制可以在 release 下载, 直接运行是替换, `./copilot-proxy recover` 是恢复. 因为插件会更新, 插件每次更新后都要重新运行次程序. 目前仅支持 VSCode 与 JetBrains 的 Copilot 插件. 如果你的代理 server 是自行部署的话, 请自己修改代码中的 api 地址, 自行编译.

有问题请看源码, 谢谢茄子🍆.

btw.
不放心 client 程序也可以自己去改, 给两个插件地址的例子, 其余的请举一反三:

```sh
/Users/admin/.vscode/extensions/github.copilot-1.78.9758/dist/
/Users/admin/Library/Application Support/JetBrains/WebStorm2022.1/plugins/github-copilot-intellij/copilot-agent/dist/
```

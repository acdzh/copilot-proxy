# 概述

该项目可以帮助你使用自建服务加速 Copilot 的 api 接口, 即把 Copilot 对 `https://copilot-proxy.githubusercontent.com/*` 的请求改为你自己的服务接口.

# 使用

项目分为两部分, client 与 server. client 旨在修改本地 Copilot 插件的 api 地址, 使其指向自己的服务. server 是自建服务的代码.

有问题请看源码, 谢谢茄子🍆.

# 编译

go 相关的实现可以直接在 release 里下载预编译好的二进制, 也可以用 build.sh 自行编译.

# Server

server 请自建. 需要注意的是, copilot 的 api 请求频率相当之高, 如果你有计划部署在 cloudflare worker 或 vercel 等公共服务上, 请**密切注意你的 api 限额**.

## 自建服务

请部署在可以畅通访问 Copilot 的服务器上, 可以通过如下命令判断:

```sh
curl -L https://copilot-proxy.githubusercontent.com/_ping
```

部署在本机也是可以的, 但需要同时使用你的代理 (vpn).

### go

相关文件对应 [proxy.go](./proxy.go).

`transport` 相关代码为代理的逻辑, 如果需要的话, 请取消注释, 并自行修改为你的代理配置. 如果不需要使用代理的话, 可以删掉.

交叉编译:

```sh
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build proxy.go
# linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build proxy.go
# windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build proxy.go
```

### node

相关文件对应 [proxy.js](./proxy.js).

`clashAgent` 相关代码为代理的逻辑, 如果需要的话, 请取消注释, 并自行修改为你的代理配置. 如果不需要使用代理的话, 可以删掉.


## CloudFlare Workers

相关文件对应 [cloudflare-worker.js](./cloudflare-worker.js).

注意: 由于众所周知的原因, CloudFlare Workers 的公用域名 `*.workers.dev` 在赛里斯无法访问, 请使用你自己的域名.

注意 2: CloudFlare Workers 的限额是 10,0000 次每天, 请注意你的 api 限额.

部署请参考: [使用 Cloudflare Workers 让 OpenAI API 绕过 GFW 且避免被封禁 · noobnooc/noobnooc · Discussion #9](https://github.com/noobnooc/noobnooc/discussions/9)


## 为什么没有公共服务

- 不安全
- copilot 的 api 请求频率太高了, 顶不住.

# Client

查找各个编辑器或 IDE 的 Copilot 插件地址, 替换其中的 api 地址为自己的服务地址. 目前仅支持 VSCode 与 JetBrains 的 Copilot 插件. 因为插件会更新, 插件每次更新后都要重新运行次程序. 如果发现失效, 也可以重新运行程序.

修改内容如下:
1. `copilot-agent/dist` 中的 js 文件, 替换 api 地址.
2. `copilot-agent/bin` 中的二进制(不包括 windows 平台), 改为 `node ../dist/agent.js`.

**VSCode 不需要使用此程序, 目前只需要在配置文件中增加如下配置即可**

```json
{
  "github.copilot.advanced": {
    "debug.testOverrideProxyUrl": "http://127.0.0.1:9394",
    "debug.overrideProxyUrl": "http://127.0.0.1:9394"
  }
}
```

**一定要确保本地有 Node 环境, 否则插件不会生效. 如果没有 Node, 请不要使用该程序. 如果你运行了, 请使用 `-r` 来还原**

## 使用

二进制文件可以从 release 下载.

设置 api 地址:

```sh
./copilot-proxy__macos -u "http://127.0.0.1:9394"
```

恢复默认地址:

```sh
./copilot-proxy__macos -r
```

## 自行修改

不放心 client 程序也可以自己去改, 给两个插件地址的例子, 其余的请举一反三:

```sh
/Users/admin/.vscode/extensions/github.copilot-1.78.9758/
/Users/admin/Library/Application Support/Google/AndroidStudio2022.1/plugins/github-copilot-intellij/copilot-agent/
/Users/admin/Library/Application Support/JetBrains/WebStorm2022.1/plugins/github-copilot-intellij/copilot-agent/
/Users/admin/Library/Application Support/JetBrains/Toolbox/apps/WebStorm/ch-0/222.4167.31/WebStorm.app.plugins/github-copilot-intellij/copilot-agent/
```

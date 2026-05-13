# s06 — Agent 派发

## Problem

上游框架允许用户(或另一条命令)用 `@gsd-planner` 引用一个 sub-agent。框架要把这个
token 变成真实的一次调用:找到 markdown 文件,解析 frontmatter 和 `<role>` 块,
组装一段包含 persona 的提示词,跑 pre-hook,把它交给 transport(transport 才真正
和 Claude Code 说话),再跑 post-hook 处理返回文本。这一整套就是一个机制——**派发器**。

## Solution

`agents/s06-agent-dispatcher/main.go` 定义了 `Dispatcher`,初始化时给它 agents 目录
和一个 `gsd.Transport`。`Dispatcher.Dispatch(ref, userMsg)` 跑完整条流水线。本章
还附带 `EchoTransport`,把提示词前面加 `ECHO:\n` 原样返回——非常适合测试,也方便
你直接观察"如果接真模型,送出去的字节是什么样"。

`Dispatch` 是面向使用者的一站式入口:

```go
d := NewDispatcher("./agents", EchoTransport{})
d.AddPre(s05hooks.Redact())
d.AddPost(s05hooks.Log(os.Stderr, "out"))
out, err := d.Dispatch("@gsd-planner", "Plan the auth phase")
```

## How It Works

`Resolve(ref)` 去掉前导 `@`,拼上 `.md` 路径,读文件,然后用第四章引入的"frontmatter
+ `<role>`"解析器解析。返回一个 `gsd.Agent`。

`Dispatch(ref, userMsg)` 五步走:

1. `Resolve(ref)` —— 加载 agent;
2. `buildPrompt(agent, userMsg)` —— 把 persona 和用户消息组装成确定性字符串
   (第四章 `Assemble` 的精简版);
3. 顺着 `d.preHooks` 把 prompt 依次穿过每一个 `HookPre`;
4. `d.Transport.Send(prompt)` —— `EchoTransport` 直接回显;
5. 同样的方式跑 post-hook 处理 response。

Hook 这一段我们写成了 dispatcher 本地代码,没有去构造第五章的 `Runner`。形状一样
(`for h := range hooks {…}`),但内联在本章里读起来不用来回跳文件。真实工程里应该
抽到共享的 `gsd/` 包。

`EchoTransport` 的存在让本章不依赖任何外部服务,测试是封闭的。以后想接真正的 Claude
Code 或 HTTP API,只要实现 `Send(prompt string) (gsd.TransportResult, error)` 即可,
dispatcher 不关心你怎么实现。

## What Changed

第五章给了我们 hook,第六章给了 hook 一个落地的地方。我们现在有一个函数:
"给我一个用户消息和一个 agent 引用,我给你一段回复" —— 正好是第七章 workflow
执行器想要的接口。

## Try It

```sh
go run ./agents/s06-agent-dispatcher
go run ./agents/s06-agent-dispatcher @gsd-executor "Apply task 2"
go test ./agents/s06-agent-dispatcher
```

预期输出:`ECHO:` 标记,后面跟着组装出来的提示词
(`You are gsd-planner.` … `=== USER ===` … 你的消息),最前最后两个 demo hook
打的 tag。

## Upstream Source Reading

- `sdk/src/cli-transport.ts` —— 上游主 transport(向 stdout 输出带 ANSI 的丰富日志);
  上游其实有多个 transport(`ws-`、`cli-`),都实现同一个接口,跟我们的 `gsd.Transport`
  完全是一个套路。
- `sdk/src/gsd-transport.ts` —— 通用 transport 契约。
- `agents/gsd-planner.md`、`agents/gsd-executor.md` 等 —— 我们派发器要解析的 agent 定义。
- `sdk/src/prompt-builder.ts` —— `parseAgentTools`、`parseAgentRole`。

带注释的摘录:[`upstream-readings/s06-transport.go.md`](../../upstream-readings/s06-transport.go.md)。

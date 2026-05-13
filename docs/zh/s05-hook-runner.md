# s05 — Hook 执行器

## Problem

第四章组装好提示词之后,你几乎不会原样发出去。真实工作流要求在最后一公里做一系列
处理:脱敏、按预算截断、做日志、去 PII、规整空白。把这些写死在组装器里是错的——
它们是横切关注点,会随着部署来来去去。标准的逃生口就是 **hook**:一个具名的小函数,
框架在发送前和接收后跑一遍。

## Solution

`agents/s05-hook-runner/main.go` 定义了 `Runner`,内部是一个 `gsd.Hook` 的有序切片。
每个 hook 声明自己的 phase(`HookPre` 或 `HookPost`)和一个变换字符串的 `HookFunc`。
`Runner.Run(phase, payload)` 按注册顺序遍历所有匹配 phase 的 hook,把上一个的输出
喂给下一个。任何 hook 返回错误就立刻终止整个链路。

我们提供三个具体 hook:`Redact()`(基于正则的脱敏,pre)、`Truncate(max)`(按 rune
数截断,pre)、`Log(w, label)`(透传式日志,post)。

## How It Works

`Runner.Run` 是核心。它按 phase 过滤,把"运行中的 payload"传给每个匹配的 hook,
错误通过 `fmt.Errorf` 包装,这样调用方一眼能看出是哪个 hook 失败。Hook 本身是值,
不是 interface 实现,这意味着你可以非常方便地内联声明一次性 hook(错误传播测试就用到了)。

`Redact` 在包初始化时通过 `regexp.MustCompile` 编译一次正则,然后用 `ReplaceAllString`。
模式刻意写得宽松:任何与 `api_key`、`password`、`secret`、`token` 临近的 token 都
会被替换成 `[REDACTED]`。误报便宜,漏报昂贵。

`Truncate` 在 `[]rune` 上操作而非字节,避免把多字节字符切成半个。低于上限就原样返回。

`Log` 很有意思:它是一个**不修改 payload** 的 hook。把 payload 的字节数写到配置好的
writer,然后原样返回输入。这样你可以在两个修改型 hook 之间插入 `Log`,得到免费的
观测能力。

## What Changed

我们现在能用横切式的变换"夹住"组装器输出。第六章会用这个 runner 把派发包起来:
prompt 在交给 transport 之前先 `Run(HookPre, …)`,response 拿回来后 `Run(HookPost, …)`。

## Try It

```sh
go run ./agents/s05-hook-runner
go test ./agents/s05-hook-runner
```

demo 会先在一段带"假 API key"的 prompt 上跑 Redact + Truncate,打印清理后的文本,
然后对一段较长的合成响应跑 Log。

## Upstream Source Reading

- `hooks/lib/git-cmd.js` —— 上游把 hook 写成独立的 JS 文件;runtime 是 Node,
  不是 Go,但契约一致:从 stdin 读 payload,把可能修改过的 payload 写到 stdout。
- `sdk/src/prompt-sanitizer.ts` —— 上游"发送前最终清理"那一层,即使内联实现,
  本质也是 hook。
- `agents/gsd-planner.md` —— 注意 frontmatter 里被注释掉的 `# hooks:` 段。agent
  可以声明自己的 hook,框架在加载 agent 时会把它们接进来。

带注释的摘录:[`upstream-readings/s05-hooks.go.md`](../../upstream-readings/s05-hooks.go.md)。

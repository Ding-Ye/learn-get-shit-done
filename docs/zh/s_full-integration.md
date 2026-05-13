# 端到端串讲 —— `/gsd plan "fix the auth bug"` 全流程

## Problem

八章已经分别讲完八个机制。考试题是:当用户敲下一条命令,这些机制是怎么组合起来的?
这份文档跟踪一次假想的调用 `/gsd plan "fix the auth bug"` 穿过每一站。本章不引入
新代码,只是把你已经写过的章节按执行顺序串成一份阅读地图。

## Solution

跟着数据流走:

```
用户输入字符串
   │
   ▼  第二章: ParseCommand
gsd.Command{Name:"plan", Flags:{}, Args:["fix","the","auth","bug"]}
   │
   ▼  第一章: LoadSkill("commands/gsd/plan.md")
gsd.Skill{Name:"gsd:plan", AllowedTools:[...], Body:"..."}
   │
   ▼  第三章: engine.Render("dev", {User:"Ding"})
渲染后的 context 字符串
   │
   ▼  第六章: dispatcher.Resolve("@gsd-planner")
gsd.Agent{Name:"gsd-planner", Role:"You produce PLAN.md..."}
   │
   ▼  第四章: Assemble(Inputs{skill, command, agent, context, user})
最终 prompt 字符串
   │
   ▼  第五章: runner.Run(HookPre, prompt) —— 脱敏、截断
清理后的 prompt
   │
   ▼  第六章: transport.Send(prompt)
gsd.TransportResult{Response:"..."}
   │
   ▼  第五章: runner.Run(HookPost, response) —— 日志
最终展示给用户的回复
```

## How It Works

按章节注明来源,一步步走:

1. **解析用户输入。** `agents/s02-command-parser` 把原始行变成
   `gsd.Command{Name: "plan", Args: ["fix", "the", "auth", "bug"]}`。
   引号、flag 都在这里处理完,后面的环节就不必再操心字符串。

2. **加载对应 skill。** 命令的 `Name` 告诉我们要去读 `commands/gsd/plan.md`
   (我们会配置一个 skills 目录)。第一章的 `LoadSkill` 返回带 `AllowedTools`
   和 `Body` 的 `gsd.Skill`。

3. **决定渲染哪个 context。** Skill body 或一份 config 决定使用 `"dev"` 或
   `"review"`。我们把它喂给第三章的 `engine.Render("dev", vars)`,它会读
   `contexts/dev.md`,合并变量(调用方覆盖 frontmatter 默认值),
   用 `text/template` 渲染,输出一段可以直接贴进去的字符串。

4. **解析 agent persona。** Skill 可能引用 `@gsd-planner`(workflow 经常这么干)。
   第六章 dispatcher 的 `Resolve("@gsd-planner")` 读 `agents/gsd-planner.md`,
   返回带 `Tools`、`Description` 和 `<role>` 的 `gsd.Agent`。

5. **组装最终 prompt。** 第四章的 `Assemble(Inputs{...})` 把 SYSTEM(agent)、
   CONTEXT(已渲染)、COMMAND(skill body + 结构化调用)、USER(原始输入)
   按固定顺序拼起来。输出就是 transport 实际要发送的那串确定性字符串。

6. **跑 pre-hook。** 第五章的 `Runner.Run(HookPre, prompt)` 顺着注册的 hook 走:
   `Redact()` 抹掉 API key,`Truncate(N)` 让 prompt 不超预算,可选的内联 hook
   做 trace 打标。输出是真正离开本机的字节。

7. **派发。** 第六章的 `Transport.Send(prompt)` 返回 `gsd.TransportResult`。
   测试里是 `EchoTransport`;生产里是真正的 Claude Code transport。

8. **跑 post-hook。** 第五章的 `Runner.Run(HookPost, response)` 是回程的对称步:
   日志、脱敏、规整。

如果用户敲的是 `/gsd workflow note "Pay the bill"`,第七章的
`Run(workflow, vars, dispatcher.Dispatch)` 就会取代第 7 步,把多条命令串起来跑。
如果用户敲的是 `/gsd plan-phase auth-refactor`,第八章的 `LoadSpec` + `ExecutePhase`
会在前面再加一层:先把目标拆 phase,再对每个 phase 走上面的整条流水线。

## What Changed

本文不改代码。它是一份**地图**,让你按"执行顺序"而不是"我介绍它们的顺序"来看
八个章节。你能复述这条 trace,就算把这套框架内化了。

## Try It

可以自己写一段短的 Go 程序把整条 trace 串起来(留给读者作为练习)。或者就按上面的
顺序读每一章的 `main.go`:每一章都把自己的那一段打印了出来,顺着读一遍就在脑子里
重演一次集成。

```sh
go run ./agents/s02-command-parser '/gsd plan fix the auth bug'
go run ./agents/s01-skill-loader  agents/s04-prompt-assembler/testdata/plan.md
go run ./agents/s03-context-engine dev Ding
go run ./agents/s06-agent-dispatcher @gsd-planner "fix the auth bug"
```

## Upstream Source Reading

- `sdk/src/cli.ts` —— 上游真正的端到端入口。最后再读它;读到这里你会发现它的形状
  非常熟悉,因为它的每个零件都是你亲手写过的。
- `sdk/src/index.ts` —— 顶层 `GSD` 类,把所有零件接起来。
- `sdk/src/session-runner.ts` —— 长生命周期 session,会驱动 workflow。

带注释的摘录:每一章对应 [`upstream-readings/`](../../upstream-readings) 下的一个文件。

# s04 — 提示词组装

## Problem

到第三章我们有了三份结构化的值:`gsd.Skill`、`gsd.Command`,以及一段渲染好的
context 字符串。我们也有 agent persona 以同样的 markdown+frontmatter 形式躺在磁盘上。
任何 transport 在发送之前,都得先把这些**组装**成一份最终提示词 —— 一个段落顺序
确定的字符串,这样缓存和 diff 才有可预期的行为。本章就是写这个组装器。

## Solution

`Assemble(in Inputs) string` 是对外的入口。它的 `Inputs` 结构体里装着已解析的 skill、
已解析的 command、agent persona、已渲染好的 context 片段,以及原始用户输入。函数
按固定顺序写出一段字符串,共四个段落:`=== SYSTEM ===`(agent persona)、
`=== CONTEXT ===`(渲染好的片段)、`=== COMMAND ===`(skill body + 结构化调用)、
`=== USER ===`(原始用户消息)。

本章还附带一个约 20 行的 `loadAgent` 小工具,负责解析 agent markdown 文件
(frontmatter 形态和 skill 一样,但用 `tools` 代替 `allowed-tools`,而且 body 里
可能有 `<role>…</role>`)。把它写在本章里是有意的:第四章自洽,第六章(派发器)
会用更完整的版本替换它。

## How It Works

组装器就是一个 `strings.Builder` 加一串 `fmt.Fprintf`。关键决定是**段落顺序**,
为什么是这个顺序?

- SYSTEM 在最前,因为真正的 Claude API 调用要求 system prompt 排在最前面;
  即使我们还没在构建 API payload,也要让组装出来的文本读起来像真实调用的形状;
- CONTEXT 第二,因为它可复用、确定、值得缓存;
- COMMAND 第三,这是"本次调用"最强相关的内容 —— skill body 加上结构化的调用参数;
- USER 最后,因为这里是自由文本。

两个小细节:每段写之前都 trim 尾部空白,避免空行堆积;COMMAND 段把 flag 打成
`--name = value`(布尔 flag 单独写 `--name`)。这样 `go run` 之后 `less` 看会非常清晰。

## What Changed

我们终于有一个函数,它的返回值可以**直接交给 transport**了。后续章节都在它上面叠加:
第五章用 pre/post hook 包住 `Assemble`;第六章把硬编码的 fixture 加载换成真正的派发器;
第七、八章加上 workflow 和 spec 编排。下游所有东西都依赖这个组装器:同样的输入,
每次都产出同样的字符串。

## Try It

```sh
go run ./agents/s04-prompt-assembler
go run ./agents/s04-prompt-assembler | less
go test ./agents/s04-prompt-assembler
```

预期输出以 `=== SYSTEM ===` 开头,出现 `gsd-planner`,打印 agent 的 `<role>` 块,
然后是 CONTEXT,然后是 COMMAND 块(`--phase = 03` 和 args `03-auth`),最后是 USER 段。

## Upstream Source Reading

- `sdk/src/prompt-builder.ts` —— 看 `buildExecutorPrompt`、`formatTask`、
  `parseAgentRole`。上游的 prompt-builder 是本章的精神近亲。
- `sdk/src/assembled-prompts.test.ts` —— 黄金文件测试,固化上游组装结果的格式;
  顺手翻一下,体会上游对"确定性"的重视。
- `sdk/src/phase-prompt.ts` —— phase 维度的提示词组装,搭在基本 builder 之上。

带注释的摘录:[`upstream-readings/s04-prompt-builder.go.md`](../../upstream-readings/s04-prompt-builder.go.md)。

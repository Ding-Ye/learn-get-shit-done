# s03 — 上下文引擎

## Problem

元提示词框架很少每次都从零拼一份提示词,通常是用一堆"命名好的、可复用的"片段
组合而成 —— "dev 上下文"、"review 上下文"、"planning 上下文"。每个片段都是一份
markdown,中间留着每次调用都会变的变量(用户名、任务、目标文件)。我们要的就是
一个小引擎:按名字加载片段,把调用方提供的变量填进去再返回。

## Solution

`agents/s03-context-engine/main.go` 暴露一个以目录为根的 `Engine`。
`Engine.Load(name)` 读 `<name>.md`,剥掉可选的 `---` frontmatter,把 body 留下,
并把 frontmatter 里所有 `key: value` 当作默认变量。
`Engine.Render(name, vars)` 把 body 当作 Go `text/template`,把默认变量和调用方
变量合并(调用方优先),然后用 `Option("missingkey=error")` 执行模板:缺变量直接
报错,而不是渲染成 `<no value>`。

## How It Works

`parseContext` 基本上是第一章 frontmatter 解析器的简化版:只关心 `key: value`,
列表暂时不要。解析失败的行直接忽略,做向前兼容。

`Render` 顺序做三件事:

1. 加载模板(调用 `Load`);
2. 构造合并后的 `map[string]string`(默认值 + 调用方变量,后者覆盖前者);
3. 把 body 解析成 `text/template`,设置 `Option("missingkey=error")`,然后用合并
   map 执行,返回渲染结果。

为什么选 `text/template` 而不是写正则替换?两个理由。第一,`{{.Var}}` 这个分隔符
Go 同行都很熟,在 markdown 里走一遭也不会被改坏。第二,我们免费拿到了 if/range 等
能力 —— 后续章节里如果要在 context 片段里循环一组文件,直接就能用。

## What Changed

我们现在可以**从命名模板生成文本**了。配合第一章解析出的 `gsd.Skill` 和第二章解析
出的 `gsd.Command`,第四章就能搭"prompt 组装器":command 告诉我们要加载哪个 context,
引擎把它渲染出来,组装器把所有片段粘起来。

## Try It

```sh
go run ./agents/s03-context-engine dev Ding
go run ./agents/s03-context-engine review
go test ./agents/s03-context-engine
```

`dev` 那条命令的预期输出里应该有 "Hello Ding." 以及替换好的任务句。

## Upstream Source Reading

- `sdk/src/context-engine.ts` —— 上游的 context 引擎。比我们的复杂得多:
  按 phase 维护文件清单、对超大文件做截断、把 ROADMAP.md 收窄到当前 milestone。
  我们刻意把这些都丢掉,只保留核心思想。
- `get-shit-done/contexts/dev.md`、`get-shit-done/contexts/research.md`、
  `get-shit-done/contexts/review.md` —— 示例片段。

带注释的摘录:[`upstream-readings/s03-context-engine.go.md`](../../upstream-readings/s03-context-engine.go.md)。

# s01 — Skill 加载器

## Problem

像上游 `get-shit-done` 这样的 Claude Code 框架,会把能力以 **skill** 的形式暴露
给模型。每个 skill 都是一份 markdown 文件,文件开头有一段 YAML frontmatter,
描述 skill 的 name、description、参数提示,以及它允许使用哪些工具;`---` 之后
的正文就是模型在被调用时会看到的提示词。在我们做提示词组装、Hook、派发之前,
首先要能把这份文件变成一个类型化的内存对象。第一章干的就是这件事:读一个文件,
返回一个 `gsd.Skill`。

## Solution

我们在 `agents/s01-skill-loader/main.go` 里写一个非常小的 loader。它接收一个路径,
读出字节,找到 `---` frontmatter 围栏,解析上游真正用到的那一小撮 YAML(`key: value`
和 `allowed-tools:` 下的字符串列表),返回一个 `gsd.Skill` 结构体(定义在
共享的 `gsd/` 包里)。我们刻意不引入 `gopkg.in/yaml.v3`:上游用到的语法非常小,
自己手写一个解析器既能控制在百行左右,又不会引入外部依赖。

## How It Works

`LoadSkill(path)` 是对外入口。它把文件读出来,然后委托给 `ParseSkill(text)`,
这样测试就可以直接传字符串。ParseSkill 的流程:

1. 统一行尾,按行切分。
2. 第一行必须是 `---`,否则返回错误。
3. 向下找到下一个 `---`,这个索引就是 frontmatter 的结束位置,之后的就是 body。
4. 把两个围栏中间的内容交给 `parseFrontmatter`。

`parseFrontmatter` 一遍扫描。它只维护一个状态:一个名叫 `currentList` 的指针,
当上一行非空行是 `allowed-tools:` 这种"列表标题行"时,它非 nil。在 `currentList`
非空、当前行以 `"  - "` 开头时,把整行 trim 后追加到列表里。其他任何行都会把
`currentList` 重置为 nil,并把当前行当作 `key: value` 来解析。我们也支持内联形式
`allowed-tools: Read, Write, Edit`:当 value 非空时按逗号切分。**未知 key 默默忽略**,
这样上游加字段时这个 loader 也不会坏。

`main` 函数提供一个非常薄的 CLI:不带参数时读取 `testdata/capture.md` 这个 fixture,
带一个参数时按路径读取。把解析出的字段打出来,`go run` 就能看到效果。

## What Changed

这是第一章,基线是空的 `gsd/` 共享类型包。增量是:我们现在有了一个真正的函数,
能把磁盘上的字节变成一个 `gsd.Skill`。后续章节直接消费这个结构体,不必再次解析文件。

## Try It

```sh
go run ./agents/s01-skill-loader
go run ./agents/s01-skill-loader agents/s01-skill-loader/testdata/capture.md
go test ./agents/s01-skill-loader
```

预期输出(节选):

```
Loaded skill: gsd:capture
  description : Capture ideas, tasks, notes, and seeds to their destination
  arg-hint    : [--note | --backlog | --seed | --list] [text]
  tools       : Read, Write, Edit, Bash
  body (first 80 chars): # Capture
```

## Upstream Source Reading

- `commands/gsd/capture.md` —— 我们要模拟的文件形态。
- `agents/gsd-planner.md` —— agent 文件,frontmatter 约定相同,只是字段略有不同
  (`tools`、`color`、可选的 `hooks`)。
- `sdk/src/prompt-builder.ts` —— 留意 `parseAgentTools` 与 `parseAgentRole`;
  上游用正则去抓取同样的字段。

带注释的摘录:[`upstream-readings/s01-skill-frontmatter.go.md`](../../upstream-readings/s01-skill-frontmatter.go.md)。

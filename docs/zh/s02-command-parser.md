# s02 — 命令解析

## Problem

第一章我们能把 skill 文件加载进来了,但还得有人**调用**它。上游框架里,用户输入
的是 `/gsd capture --note "Pay the electric bill"` 这样一行字。它是整个系统的入口:
带命令名、带 flag、带位置参数。在它被解析成结构化值之前,下游的 context 引擎、
prompt 组装器、派发器都拿不到任何东西可以处理。本章就是把这行字变成结构体的那个解析器。

## Solution

`agents/s02-command-parser/main.go` 暴露 `ParseCommand(line string) (gsd.Command, error)`。
它先按 token 切分(尊重 `"…"` 和 `'…'`),去掉可选的前导斜杠,同时接受
`gsd capture`(空格形式)和 `gsd:capture`(冒号形式)两种 namespace,然后扫描剩余
token,识别 `--flag`、`--flag value`、`--flag=value`。不是 flag 的裸 token 按顺序
进入 `cmd.Args`。最终得到一个 `gsd.Command`,带 `Name`、`Flags`(map)和 `Args`
(slice)—— 后续章节需要的形状。

## How It Works

`tokenize` 是个非常小的状态机。它按 rune 扫一遍输入,只跟踪一个状态:是否在引号里。
遇到引号就翻转状态;在引号外遇到空白就 flush 当前 token;其余情况累加。我们不支持
转义序列——引号成对就够了。遇到不闭合的引号会直接返回错误,让调用方拿到确定性的失败。

`ParseCommand` 再处理第一个 token 作为 namespace/name。如果第一个 token 含 `:`,
冒号之后就是命令名;如果恰好是 `gsd` 后面还跟一个 token,就把下一个 token 当作
命令名;否则第一个 token **就是**命令名(这样在 workflow 文件里写 `fast --text hi`
也合法)。

Flag 处理:以 `--` 开头就是 flag。`--flag=value` 按 `=` 切;`--flag` 后面跟一个
非 flag token 就消费成 value;后面紧跟另一个 `--flag` 或结束,就保留为空——我们把
它理解成"布尔:存在"。flag 之后的第一个非 flag token 仍然进入 Args,我们不做重排序。

## What Changed

我们已经能用真正的解析器填充 `gsd` 包里的第二个运行时类型了。和第一章的 loader 合起来,
用户输入的命令行 + 磁盘上的 skill 文件,都已经被装进类型化的 Go 值。第三、四章可以
开始搭建 context 和 prompt 层。

## Try It

```sh
go run ./agents/s02-command-parser
go run ./agents/s02-command-parser '/gsd capture --note "Pay the bill"'
go run ./agents/s02-command-parser 'gsd:plan --phase=03 --gaps'
go test ./agents/s02-command-parser
```

第二条命令的预期输出:

```
raw   : /gsd capture --note "Pay the bill"
name  : capture
flags : map[note:Pay the bill]
args  : []
```

## Upstream Source Reading

- `sdk/src/cli.ts` —— 注意 `parseArgs` 与 `parseCliArgsQueryPermissive`。上游用
  Node 的 `parseArgs` + 一层"宽松"包装层来兼容未知 flag。
- `commands/gsd/capture.md` —— 看 `<context>Arguments: $ARGUMENTS</context>` 部分,
  这个 `$ARGUMENTS` 字符串就是被解析的对象。

带注释的摘录:[`upstream-readings/s02-cli-args.go.md`](../../upstream-readings/s02-cli-args.go.md)。

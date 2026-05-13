# s07 — Workflow 执行器

## Problem

有些用户意图一条命令就够(第二章覆盖)。但很多有用的操作是"链":追加一条 note,
然后列出所有 note;捕获一个 todo,再把它挂到 milestone;先跑 planner,再路由到
executor。上游框架把这些叫 **workflow**,用 markdown 文件声明;运行时按文件顺序
派发每一步。

## Solution

`agents/s07-workflow-runner/main.go` 提供 `LoadWorkflow(path)`:解析我们简化的
workflow markdown 形态(frontmatter 写元数据 + 编号列表写每一步,步骤内部是
`key: value`);以及 `Run(w, vars, runCmd)` 来执行。每一步给一个命令名、一个 args
字符串(支持 `{{.Var}}` 替换)、可选 `next` 步骤 ID 用于跳转、可选 `on_error`
步骤 ID 用于恢复。runner 返回一份文本 trace,第一个未被处理的错误会被原样冒出来。

## How It Works

`parseWorkflow` 两遍走:先剥 frontmatter(与第一章 loader 同样的围栏检测),然后
扫 body,遇到 `<数字>.` 或 `<数字>)` 开头的行就开新步骤。后续缩进行按 `key: value`
解析,映射到 `WorkflowStep` 的字段(`id`、`command`、`args`、`next`、`on_error`、
`describe`)。第一个编号之前的散文一律忽略——这样你可以在同一份文件里把"叙事说明"
和"可执行步骤"放一起,解析器不会被绕晕。

`Run` 用 `map[string]WorkflowStep` 按 ID 索引步骤,从声明顺序的第一个开始,顺着
`Next` ID 走。对每一步:

1. 用 `vars` map 替换 args 里的 `{{.Var}}`;
2. 调用 `runCmd(step.Command, args)`;
3. 成功:写一行 trace,跟着 `Next` 走。`Next` 为空就终止;
4. 失败:`OnError` 设了就跳到那一步继续,否则把错误包上步骤 ID 抛出去。

我们设了 100 步的安全上限,防止无限循环。任何真实 workflow 都不会用满,
开发中误写 `next:` 自指的时候这个上限会立刻把问题暴露出来。

## What Changed

我们抬到了上一层粒度。一到六章每章处理一个机制;第七章把这些机制(主要是命令派发)
按"用户没亲手写出来"的顺序排起来跑。`runCmd` 这个回调就是 workflow 层和第六章
dispatcher 的接缝:真正使用时你会闭包持有一个 `Dispatcher`,把每一步路由给它。

## Try It

```sh
go run ./agents/s07-workflow-runner
go run ./agents/s07-workflow-runner agents/s07-workflow-runner/testdata/with-error.md
go test ./agents/s07-workflow-runner
```

默认 workflow 是 `note.md`:两步,假装追加一条 note 再列出来。
`with-error.md` workflow 演示 `OnError` 触发时走的恢复路径。

## Upstream Source Reading

- `get-shit-done/workflows/note.md` —— 上游真正的 note workflow。比我们的丰富很多
  (多模式、交互式问询),但步骤结构还是一眼能对上。
- `get-shit-done/workflows/add-todo.md` —— 多步 workflow 示例:推断 area,落盘文件。
- `sdk/src/session-runner.ts` —— session 级 driver,会调 workflow。

带注释的摘录:[`upstream-readings/s07-workflow.go.md`](../../upstream-readings/s07-workflow.go.md)。

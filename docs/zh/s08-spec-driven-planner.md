# s08 — Spec 驱动规划

## Problem

上游框架对"大工作怎么干"有自己的强观点:先写 SPEC,把需求锁住,然后拆成多个 phase,
每个 phase 都小到能独立验证、独立交付。Planner 读 spec,Executor 实现一个 phase,
Verifier 验,再走下一个 phase。这就是 spec-driven 的纪律。本章用最小可信的形式
把它移植过来:把 spec 解析成结构化的 phase,然后用可插拔的 verifier 一次执行一个 phase。

## Solution

`agents/s08-spec-driven-planner/main.go` 暴露三个函数:

- `LoadSpec(path)` —— 把 markdown 文件解析成 `gsd.Spec`;
- `ExecutePhase(spec, phaseID, verifier)` —— 执行单个 phase,成功就标记 `Complete`,
  返回更新后的 spec 和一份日志;
- `NextIncomplete(spec)` —— 找到下一个 `Complete` 仍为 false 的 phase;UI 的"跑下一个"
  按钮要用的。

Spec 文件格式是扁平 markdown:顶层 `# Title`、一个 `## Goal` 段,然后每个 phase
是一个 `## Phase <ID>: <Title>`,里面有 `Goal:` / `Tasks:`(带短横线的列表) / `Verify:`。

## How It Works

`parseSpec` 是个小型状态机,按行扫描,维护一个状态变量(`outer`、`goal`、`phase`、`tasks`)。
phase 标题正则 `^##\s*Phase\s+(\S+):\s*(.+?)\s*$` 同时捕获 ID(`P01`)和标题
(`Token issuance`)。任何别的 `## ` 都会关闭当前的 Goal/Phase。phase 内部:
前缀匹配 `Goal:`、`Verify:` 直接写字段;`Tasks:` 切到 tasks 子状态;子状态下短横线
开头的行追加到 `current.Tasks`。

`ExecutePhase` 刻意写得小、几乎纯函数 —— 它自己不决定 phase 是否通过。接受一个
`verifier func(gsd.SpecPhase) error` 回调,这就是接缝:测试里 verifier 确定性地
返回 `nil` 或 `errors.New(...)`;真实情况下它可能跑 `go test`、打个 healthcheck、
或者 grep 一段日志。成功了 ExecutePhase 把 `p.Complete` 置 true 并写日志;失败了
返回 verifier 的错误,前面包一层 phase ID。

`NextIncomplete` 就是个 for 循环。存在的意义是让"跑下一个 phase 试试"这个集成故事
有一个干净的入口。

## What Changed

第七章是把用户自己写得出来的命令串起来;第八章串的是用户**没写**的命令——结构由 SPEC
文件给出来。有了本章,框架可以从"这是我的目标"产出"第一个可执行的步骤是什么,
我们怎么验证它做完了"。这就是 spec-driven 开发的精髓。

## Try It

```sh
go run ./agents/s08-spec-driven-planner
go test ./agents/s08-spec-driven-planner
```

默认 demo 加载 `auth-refactor.md`,打印解析出的 spec,然后用一个"`Verify:` 非空就过"
的 stub verifier 执行 `P01`。

## Upstream Source Reading

- `get-shit-done/templates/spec.md` —— 上游 SPEC 模板的正典版本。比我们丰富很多
  (ambiguity report、falsifiability 要求、boundary 条款),但区段形状还是一眼能对上。
- `sdk/src/plan-parser.ts` —— 上游解析器。我们按 `## ` 切,上游会追踪任务依赖、
  按 wave 安排并行。
- `sdk/src/phase-runner.ts` —— 一阶段一阶段执行的主循环。
- `agents/gsd-planner.md`、`agents/gsd-verifier.md` —— 产出 phase 和验证 phase 的两个 agent。

带注释的摘录:[`upstream-readings/s08-spec-phase.go.md`](../../upstream-readings/s08-spec-phase.go.md)。

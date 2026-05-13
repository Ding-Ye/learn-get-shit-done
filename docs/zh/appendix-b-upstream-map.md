# 附录 B —— 上游源码地图:去哪里读什么

一页纸的索引,把每一章映射到上游最值得读的文件。
上游提交:`ba625c0978e1133ead8a26546b5e6a435e886edb`(MIT)。

## 顶层入口

- `CLAUDE.md` —— Claude 最先读的顶层指令。
- `CONTEXT.md` —— 示例项目的 context 文件(也是上游第三章对应组件按 phase 加载的
  文件之一)。
- `package.json`、`tsconfig.json` —— 工程形态(依赖图很大,先不看,把其他都理解了再回来)。
- `sdk/src/cli.ts` —— 主 CLI 入口。`parseArgs` 的形状和我们第二章 `gsd.Command` 同款。
- `sdk/src/index.ts` —— 顶层 `GSD` 类,把所有东西接起来。

## 按章对照

### s01 —— Skill loader

- `commands/gsd/capture.md`、`commands/gsd/plan-phase.md`、
  `commands/gsd/execute-phase.md` —— 典型 command 形态。
- `agents/gsd-planner.md`、`agents/gsd-executor.md` —— 典型 agent 形态。
- `sdk/src/prompt-builder.ts` —— 看 `parseAgentTools`、`parseAgentRole`
  (和我们一样,用正则啃 frontmatter)。

### s02 —— Command parser

- `sdk/src/cli.ts` —— `ParsedCliArgs`、`parseCliArgsQueryPermissive`。
- `sdk/src/init-runner.ts` —— `init` 的特殊解析(支持 `@file`)。
- `commands/gsd/capture.md` —— `<context>Arguments: $ARGUMENTS</context>` 这段,
  能看到"被解析后的输入到模型眼里长什么样"。

### s03 —— Context engine

- `sdk/src/context-engine.ts` —— 完整引擎,按 phase 维护文件清单。
- `sdk/src/context-truncation.ts` —— 感知 heading 的截断。
- `get-shit-done/contexts/dev.md`、`research.md`、`review.md` —— 真正的 context profile。
- `get-shit-done/templates/context.md` —— 新项目 `CONTEXT.md` 的种子模板。

### s04 —— Prompt assembler

- `sdk/src/prompt-builder.ts` —— `buildExecutorPrompt`、`formatTask`。
- `sdk/src/assembled-prompts.test.ts` —— 组装结果的黄金文件测试。
- `sdk/src/phase-prompt.ts` —— phase 维度的 prompt 组装,搭在基础 builder 之上。
- `sdk/src/prompt-sanitizer.ts` —— 发送前最终清理。

### s05 —— Hook runner

- `hooks/lib/git-cmd.js` —— 示例 hook 实现。
- `sdk/src/prompt-sanitizer.ts` —— 发送前清理层。
- `sdk/src/logger.ts` —— 可观测性 hook 栈。
- `agents/gsd-planner.md` 里被注释掉的 `# hooks:` 段 —— agent 自己声明 per-event hook 的写法。

### s06 —— Agent dispatcher

- `sdk/src/cli-transport.ts` —— ANSI/终端 transport。
- `sdk/src/gsd-transport.ts` —— transport 接口契约。
- `sdk/src/event-stream.ts` —— 事件类型和流式管道。
- `sdk/src/query/query-subprocess-adapter.ts` —— 上游真正 shell 出 `claude` 二进制的地方。
- `agents/gsd-*.md` —— dispatcher 要解析的全部 agent 定义。

### s07 —— Workflow runner

- `get-shit-done/workflows/note.md` —— 真正的 note workflow。
- `get-shit-done/workflows/add-todo.md` —— 多步 + area 推断。
- `get-shit-done/workflows/plan-phase.md` —— 调用 planner sub-agent。
- `sdk/src/session-runner.ts` —— session 级驱动。
- `sdk/src/milestone-runner.test.ts` —— 多 workflow 的 milestone 测试。

### s08 —— Spec-driven planner

- `get-shit-done/templates/spec.md` —— 正典 SPEC 模板。
- `get-shit-done/templates/phase-prompt.md` —— phase 执行 prompt 模板。
- `sdk/src/plan-parser.ts` —— 带依赖分析的完整 plan 解析器。
- `sdk/src/phase-runner.ts` —— 完整 phase 循环。
- `sdk/src/planning-journal.ts` —— 只增不改的 phase 迁移审计日志。
- `agents/gsd-planner.md`、`agents/gsd-verifier.md`、`agents/gsd-executor.md`
  —— 产 phase、验 phase、执行 phase 的三个 agent。

## 学完八章之后怎么读上游

1. 从 `sdk/src/cli.ts` 开始。找到 `parseArgs`——那是第二章。
2. 跟着 `GSD` 的 import 跳到 `sdk/src/index.ts`。找到调 prompt builder 的那一行——
   那是第四章。
3. 打开 `sdk/src/prompt-builder.ts`。里面 agent 加载那部分——那是第一章。
4. 沿着 transport 调用跳——`sdk/src/cli-transport.ts` 是第六章。
5. 找 spec/plan 解析器——`sdk/src/plan-parser.ts` 是第八章。
6. 翻一遍 `get-shit-done/workflows/` 下的 workflow 文件——那是第七章。

按这个顺序走完,你应该能基本读懂上游里的每个文件。

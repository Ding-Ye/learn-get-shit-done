# learn-get-shit-done

一份从零开始、动手实践的 **Go 语言**学习仓库,目标是配套阅读上游项目
[gsd-build/get-shit-done](https://github.com/gsd-build/get-shit-done) ——
一个面向 Claude Code 的元提示词 / 上下文工程框架(TypeScript)。

本仓库会带你**一章一个机制**地在 Go 里搭建一个袖珍版的 Claude Code 元提示词框架。
每一章都是一个独立、聚焦的 Go 小程序(约 80–150 行),只演示上游系统中的一个核心点。

> 英文版请见 [`README.md`](./README.md)。

## 为什么写这个仓库

上游项目体量大、跨语言、以 TypeScript 为主。如果你从来没读过一个元提示词框架,
直接跳进约 19 万行代码里会非常劝退。本仓库把它精简到骨架,把每个核心机制改写
为约百行的地道 Go 代码,并清楚地指回上游源码,让你之后再去读时心里有底。

读完你会理解:

- 一个带 YAML 头的 markdown 文件如何变成**运行时 Skill**;
- 一条 `/gsd <command>` 斜杠命令如何被**解析**;
- **上下文模板**如何与变量合成,产出可复用的提示词片段;
- 最终提示词如何由 command + context + agent + user input **组装**而成;
- **Hook** 如何在请求前、响应后对内容做加工;
- **Agent** 如何被查找、加载,并通过 Transport 接口**派发**;
- **Workflow** 如何用声明式语法把多条命令串起来;
- 一份 **Spec** 怎么转成多个 Phase,然后一次只执行一个 Phase。

## 课程总览

| #   | 章节                    | 你会动手实现的内容                                                         |
| --- | ----------------------- | -------------------------------------------------------------------------- |
| 01  | Skill loader            | 解析带 YAML frontmatter 的 markdown skill 文件 → `gsd.Skill`。             |
| 02  | Command parser          | 把 `/gsd <name> [args]` 解析成结构化的 `gsd.Command`。                     |
| 03  | Context engine          | 从磁盘加载上下文模板,并用 `text/template` 渲染。                           |
| 04  | Prompt assembler        | 把 command + context + agent persona + user input 拼成最终提示词。         |
| 05  | Hook runner             | Pre/post hook:脱敏、截断、日志。                                           |
| 06  | Agent dispatcher        | 解析 `@agent` → 加载 → 通过 `EchoTransport` 派发。                          |
| 07  | Workflow runner         | 按 workflow markdown 中的声明顺序串起多条命令并执行。                       |
| 08  | Spec-driven plan        | 解析 spec → 拆成多个 phase → 一次执行一个 phase。                          |

端到端串讲见 [`docs/zh/s_full-integration.md`](./docs/zh/s_full-integration.md)。

两份附录:

- [附录 A — 元提示词:Claude Code 框架的形状](./docs/zh/appendix-a-meta-prompting.md)
- [附录 B — 上游源码地图](./docs/zh/appendix-b-upstream-map.md)

## 如何使用

```sh
git clone https://github.com/Ding-Ye/learn-get-shit-done.git
cd learn-get-shit-done
go vet ./...
go test ./...

# 运行任意一章的小程序
go run ./agents/s01-skill-loader
go run ./agents/s02-command-parser '/gsd capture --note 交水电费'
# ...
```

每个章节目录都包含:

- `main.go` —— 可直接运行的小 demo,有意控制在 80–150 行;
- `*_test.go` —— 覆盖主要行为的单元测试;
- `README.md` —— 指向长文文档的入口;
- `testdata/` —— 仿照上游 markdown 形态的小 fixture。

每章还配有 `docs/en/sNN-<name>.md`(英文长文)和 `docs/zh/sNN-<name>.md`(中文镜像)。
两边都按六个固定段落组织:

1. **Problem / 问题** —— 我们在模拟哪个机制?为什么它重要?
2. **Solution / 解法** —— 演示该机制所需的最小 Go 构造。
3. **How It Works / 工作原理** —— 逐行讲解代码。
4. **What Changed / 增量是什么** —— 与上一章相比新增了什么。
5. **Try It / 动手试试** —— 可直接复制的命令。
6. **Upstream Source Reading / 上游源码导读** —— 指向上游对应文件。

上游源码的**带注释摘录**放在 [`upstream-readings/`](./upstream-readings)。

## 目录结构

```
.
├── go.mod
├── gsd/                       # 共享领域类型(Skill / Command / Agent 等)
├── agents/                    # 每一章一个目录
│   ├── s01-skill-loader/
│   ├── s02-command-parser/
│   ├── s03-context-engine/
│   ├── s04-prompt-assembler/
│   ├── s05-hook-runner/
│   ├── s06-agent-dispatcher/
│   ├── s07-workflow-runner/
│   └── s08-spec-driven-planner/
├── docs/                      # en/ + zh/ 长文文档
├── upstream-readings/         # 上游源码带注释摘录(.go.md)
├── web/                       # 静态文档浏览页
└── .github/workflows/ci.yml
```

## 许可证与上游致谢

- 本仓库的 Go 源码与文档采用 **Apache 2.0** 许可证。
- 上游项目 [`gsd-build/get-shit-done`](https://github.com/gsd-build/get-shit-done)
  (提交 `ba625c0978e1133ead8a26546b5e6a435e886edb`)采用 **MIT** 许可证。
  本仓库不重发上游源码,仅在教学目的下引用少量带注释的片段。

完整许可证见 [`LICENSE`](./LICENSE)。

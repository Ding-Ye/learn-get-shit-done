# 附录 A —— 元提示词:Claude Code 框架的形状

## 框架和 chat client 的本质区别

Chat client 的工作流是:接用户消息,前面拼一段固定 system prompt,发给模型,
渲染返回。提示词是用户写的,client 只负责传输。

**元提示词框架**把这个关系反了过来:**框架**来写提示词,常常是从好几份文件里拼出来的。
用户敲一行短命令——`/gsd plan`、`@gsd-executor do task 2`、`gsd:capture --note hi`
——框架把它展开成一长段结构化的提示词,带 system 指令、context 片段、允许使用的工具、
解析过的调用、加上用户自由输入。告诉模型的不只是"做什么",还有"做的时候是谁",
以及"手上有什么工具"。

上游 `get-shit-done` 就是这种框架。把它(挺大量的)TypeScript 运行时、打包、
Claude Code 事件管道剥掉,概念骨架其实很小:

- **Skill** —— 一份带 YAML frontmatter 的 markdown,frontmatter 给框架用,
  body 是模型会看到的文本。
- **Command** —— 一种特殊 skill,`/gsd <name>` 对应 `commands/gsd/` 下的 md 文件。
- **Agent** —— 另一种特殊 skill,`@agent-name` 对应 `agents/` 下的 md 文件,
  body 里多一个 `<role>` 块。
- **Context** —— `contexts/` 下可复用的文本片段,可选带模板变量。
- **Workflow** —— `workflows/` 下的 markdown,声明一串命令的执行顺序。
- **Spec** —— `templates/` 下的 markdown,把一个目标拆成多个 phase。
- **Hook** —— 在 prompt 或 response 上做变换的函数。
- **Transport** —— 真正调用模型的适配器。

## Skill 文件的解剖

```yaml
---
name: gsd:capture                              # 用户怎么调用
description: Capture ideas to their destination  # 一句话说明
argument-hint: "[--note] [text]"               # tab 补全提示
allowed-tools: [Read, Write, Bash]             # 允许模型用哪些工具
---

<objective>...</objective>                      # 提示词正文
<process>...</process>
```

三个观察:

1. **Frontmatter 给框架,body 给模型**。这一刀让我们可以在代码里校验、渲染、路由,
   同时把散文原封不动地交给模型。
2. Body 可以用 **XML 风格标签**做"软段落"。Claude 很吃这一套,框架也可以按需
   grep(`<role>…</role>`、`<process>…</process>`)。
3. 这份文件**同时是文档**。没有第二份事实来源:`cat` 一下你就能看懂它。

## Workflow 文件的解剖

```markdown
1. id: append
   command: capture
   args: "--note {{.Text}}"
   next: list

2. id: list
   command: capture
   args: "--list"
```

Workflow 文件是**声明式**的:每一步都是结构化的 key/value。运行时按顺序走、替换变量、
调用 dispatcher。用户不参与循环。

(上游 workflow 文件写得更松散,大段散文里穿插着 `<step>` 标签,因为它们经常是给模型
本身解释的,不是给严格解析器的。我们第七章的格式形状一致,只是为代码硬化过。)

## Spec 文件的解剖

```markdown
# Auth refactor — Specification

## Goal
Replace cookie auth with JWTs.

## Phase P01: Token issuance
Goal: Mint signed JWTs on login.
Tasks:
- add /token endpoint
- sign with HS256
Verify: hitting /token returns a JWT whose exp is ~12h away.
```

Spec 是框架的"承诺工具":**每条需求都必须是可证伪的**。"提升性能"不是需求,
"P95 延迟从 800ms 降到 ≤300ms"才是。Verifier(第八章的回调)只输出 PASS 或 FAIL,
没有第三种结果。

## 组合 vs 拼接

朴素的做法是把一切都拼成一个大字符串。`get-shit-done` 走的是组合:

- 同一个 **context 片段**可以被十个不同的命令复用;
- 同一个 **agent persona** 可以被三个不同的 workflow 调用;
- 同一份 **spec** 可以被三个不同的角色(planner / executor / verifier)分别消费。

组合要求阶段之间流动的是类型化的值,而不是原文本——正好就是 `gsd/` 包做的事。
这就是"Unix 管道"和"一段巨长的 shell 脚本"的差别。

## 元提示词什么时候才值回票价

对一次性聊天来说,元提示词是过度工程。它在以下场景才挣回复杂度:

- 同一类任务你要在多个 session 里反复做;
- 你希望同样的输入每次都产出**字节相同**的提示词(缓存、回放、回归测试);
- 你想换 persona 但不想重写提示词;
- 你想要可统一管理的 hook(脱敏、日志),写一次到处复用;
- 你想让新同事不必先把全队摸索过的"提示词玄学"全部内化就能开始干活。

以上多数你都在点头,那元提示词就是你要的工具;如果不是,写一个聊天循环然后下班。

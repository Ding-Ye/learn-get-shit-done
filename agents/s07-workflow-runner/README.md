# s07 — Workflow runner

Load a workflow markdown file and execute its steps in order, with variable
substitution and on-error routing.

```sh
go run ./agents/s07-workflow-runner
go run ./agents/s07-workflow-runner agents/s07-workflow-runner/testdata/with-error.md
go test ./agents/s07-workflow-runner
```

Long-form: [`docs/en/s07-workflow-runner.md`](../../docs/en/s07-workflow-runner.md)
· [中文](../../docs/zh/s07-workflow-runner.md).

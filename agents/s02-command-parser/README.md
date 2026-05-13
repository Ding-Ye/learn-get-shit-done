# s02 — Command parser

Turn a raw line like `/gsd capture --note "Pay the bill"` into a `gsd.Command`.

```sh
go run ./agents/s02-command-parser '/gsd plan --phase=03 --gaps'
go test ./agents/s02-command-parser
```

Long-form: [`docs/en/s02-command-parser.md`](../../docs/en/s02-command-parser.md)
· [中文](../../docs/zh/s02-command-parser.md).

---
name: note
description: Append a note then list pending notes
---

1. id: append
   command: capture
   args: "--note {{.Text}}"
   next: list

2. id: list
   command: capture
   args: "--list"

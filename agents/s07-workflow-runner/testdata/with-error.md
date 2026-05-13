---
name: with-error
description: Demonstrate on_error routing
---

1. id: try
   command: broken
   args: ""
   next: done
   on_error: recover

2. id: recover
   command: ok
   args: "recovered"
   next: done

3. id: done
   command: ok
   args: "done"

---
# d6roll-whqk
title: Make curl installer command fail on download errors
status: completed
type: bug
priority: normal
created_at: 2026-10-26T22:14:41Z
updated_at: "2026-10-26T22:16:33Z"
---
README's curl -fsSL ... | sh uses the shell command's exit status and can return success even if curl fails before producing a script. Document a safe curl-download-then-execute command that gates execution on curl success, preserve inspect-before-running advice, and add an offline failure-path test.

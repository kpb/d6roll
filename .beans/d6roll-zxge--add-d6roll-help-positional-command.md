---
# d6roll-zxge
title: Add d6roll help positional command
status: completed
type: feature
priority: normal
created_at: 2026-10-26T23:00:43Z
updated_at: "2026-10-26T23:01:23Z"
---
Support `d6roll help` as a single special-case CLI invocation without adding a subcommand framework. It should print usage to stdout, exit 0, and not generate a passphrase; retain -h/--help behavior. Reject extra positional arguments such as `d6roll help extra` and keep ordinary CLI stdout script-friendly. Add deterministic CLI tests for status, output streams, and invalid arguments, update README, and verify Go tests/build/vet.
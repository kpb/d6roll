---
# d6roll-hv3c
title: Reject impractically large library word counts
status: completed
type: bug
priority: normal
created_at: 2026-10-26T22:14:40Z
updated_at: "2026-10-26T22:16:33Z"
---
Generate and GenerateWords accept any positive int and allocate a slice of that capacity, so math.MaxInt can panic instead of returning an error. Introduce a documented library-only maximum independent of the CLI 3–32 range, check before allocation or reads, and test its boundary and oversized counts across both APIs.

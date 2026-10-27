---
# d6roll-cjy6
title: Add Homebrew installation via existing kpb tap
status: completed
type: feature
priority: normal
created_at: 2026-10-27T14:29:33Z
updated_at: "2026-10-27T14:33:56Z"
---
Publish a d6roll formula in the existing kpb/homebrew-tap repository using the v0.2.1 release archives for macOS amd64/arm64 and Linux amd64, with pinned SHA-256 checksums. Install only the d6roll executable. Add a meaningful Homebrew test without logging generated passphrases; validate the formula on supported platforms. Document `brew install kpb/tap/d6roll` in the d6roll README. Update the kpb/homebrew-tap README to link to both https://github.com/kpb/beanstalk and https://github.com/kpb/d6roll and show how to install each formula. Update the kpb/homebrew-tap GitHub repository About description to reflect that it hosts both Beanstalk and d6roll rather than only Beanstalk. Maintain version/checksums on subsequent releases; automation can follow later. No CLI or library API changes or new release are needed.
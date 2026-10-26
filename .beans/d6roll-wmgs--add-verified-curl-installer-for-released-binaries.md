---
# d6roll-wmgs
title: Add verified curl installer for released binaries
status: completed
type: feature
priority: normal
created_at: 2026-10-26T21:59:19Z
updated_at: "2026-10-26T22:01:35Z"
---
Ship a POSIX install.sh for Linux amd64 and macOS amd64/arm64 using existing GitHub release archives. Resolve latest release by default, support VERSION and INSTALL_DIR overrides, download archive and checksums.txt, verify SHA-256 before extraction and installation, reject unsupported platforms and failed downloads, and avoid root by default. Document curl | sh with inspect-before-running caution, checksum trust limits, PATH guidance, and manual Windows ZIP installation. Add deterministic offline installer tests and run them on Linux/macOS CI. Keep dependencies minimal; do not add release automation or unsupported architectures.

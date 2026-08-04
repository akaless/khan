# Changelog

All notable changes to Khan project will be documented in this file.

Based on [Keep a Changelog](https://keepachangelog.com/fa-IR/1.1.0/) and [Semantic Versioning](https://semver.org/).

## [1.0.1] - 2026-08-02

### Changed
- New repository: codeberg.org/adiib/khan1.0.1 (separate from v1.0.0)
- All v1.0.0 features preserved
- Fresh repository for v1.0.1 development cycle

### Added
- Security check script: `scripts/khan-security-check.sh`

## [1.0.0] - 2026-08-02

### Added
- Complete LAN chat server (Go, WebSocket, JSON storage)
- 3 user roles: User, Supervisor, Admin
- Ed25519 license system (20 free / valid / 5 penalty)
- AES-256-GCM message encryption
- Argon2id password hashing
- Persian RTL + English UI
- Bilingual login page with language switcher (FA/EN)
- 4 platform binaries: Windows, Linux, macOS Intel, macOS ARM
- Installers: NSIS (.exe), .deb, shell scripts, DMG builder
- Professional documentation (API, installation, licensing)
- CI/CD workflows (GitHub Actions)
- Docker support

### Security
- Hidden super admin (aDiB) in all layers
- No super admin mentions in public docs
- Private keys never in repo

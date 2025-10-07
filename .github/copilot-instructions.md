# Copilot Instructions for the `kodik` Repository

This repository defines a **spec-driven, agent-powered development workflow**. It is not a traditional codebase, but a meta-tooling and documentation resource for AI coding agents and developers using platforms like GitHub Copilot and Roo Code.

## Key Concepts & Architecture
- **Agent Modes:** Four core modes—`spec`, `design`, `tasks`, `agent`—each with a dedicated `.chatmode.md` file in `.github/chatmodes/` and mirrored in `.roomodes` for Roo Code.
- **Workflow:**
  1. **Spec Mode:** Write feature specs with user stories and Cucumber scenarios (`.spec.md`).
  2. **Design Mode:** Produce detailed technical designs with ASCII diagrams, data flows, and component breakdowns (`.design.md`).
  3. **Tasks Mode:** Break designs into granular, atomic tasks with explicit validation and dependencies (`.tasks.md`).
  4. **Agent Mode:** Execute tasks strictly in order, updating status and validating each step.
- **All planning and execution artifacts are stored in the `.kodik` directory at the repo root.**

## Developer & Agent Workflows
- **Setup:** Use the `kodik` CLI to install or update `.github`, `.roomodes`, and `.opencode` in a target project. Install the CLI, then run `kodik all` (or `kodik github` / `kodik roo`) from the project root.
- **No build/test commands:** This repo contains no executable code or tests. It is a configuration and workflow definition resource.
- **To use agent modes:**
  - For Copilot: Run `kodik github` in your project’s root to install/update `.github/chatmodes/` and managed prompts (use `--dry-run` to preview).
  - For Roo Code: Run `kodik roo` in your project’s root to install/update `.roomodes` (manual copy is also supported).
- **Prompts and templates:** See `.github/prompts/` for reusable prompt templates for each mode.

## CLI Installation & Usage

Use the `kodik` CLI to manage installation and updates of this repository’s configurations in target projects.

Installation (see `README.md` for platform-specific commands):
- Download the latest binary for your OS from Releases and place it on your PATH as `kodik`.

Basic usage (run in the target project root):
- `kodik all` — Install/update all components (`.github`, `.roomodes`, `.opencode`)
- `kodik github` — Install/update GitHub Copilot chat modes and prompts
- `kodik roo` — Install/update Roo Code `.roomodes`
- `kodik opencode` — Install/update OpenCode configuration

Global flags:
- `--dry-run` — Preview actions without executing
- `--force` — Skip backups and confirmations

Key features:
- Cross-platform binaries (macOS, Linux, Windows)
- Automatic backups in `.kodik-state/backups/`
- Modification detection with SHA256 checksums
- Selective merge that preserves workflows, CODEOWNERS, and other user files
- Atomic operations with rollback and state management in `.kodik-state/`

Note: The CLI is used to manage kodik configurations in target projects; it is not embedded into those projects.

## Project-Specific Conventions
- **All planning files (`.spec.md`, `.design.md`, `.tasks.md`) must live in the `.kodik` directory.**
- **ASCII art diagrams are mandatory** in design docs for UI and data flows.
- **Task lists must be extremely granular, with status, context, dependencies, and validation for each task.**
- **Agents must never deviate from the written plan**—no improvisation or skipping steps.
- **All status changes must be explicitly announced and reflected in the task file.**
- **Review plans** are generated using the `kodik review planner` mode, synthesizing spec, design, and tasks with a detailed checklist.
- Add `.kodik-state/` to `.gitignore` in target projects to exclude backups and local state.

## Integration Points
- **GitHub Copilot:** Custom chat modes in `.github/chatmodes/`.
- **Roo Code:** Unified `.roomodes` file.
- **Install/update automation:** Use the `kodik` CLI—`kodik github`, `kodik roo`, `kodik opencode`, or `kodik all`—which performs selective merges, creates automatic backups, detects modifications via checksums, and supports dry-run/rollback.

## Examples & References
- See `.github/chatmodes/` for mode definitions and usage patterns.
- See `.github/prompts/` for prompt templates.
- See `.roomodes` for Roo Code integration and detailed agent instructions.
- See `README.md` for high-level project overview and setup instructions.

---

**Agents: Always follow the mode-specific instructions in `.chatmode.md` or `.roomodes`. Never assume, always verify against the actual files.**

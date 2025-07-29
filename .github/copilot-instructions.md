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
- **Setup:** Use `./scripts/install_update.sh` to install or update `.github` and `.roomodes` from the canonical source.
- **No build/test commands:** This repo contains no executable code or tests. It is a configuration and workflow definition resource.
- **To use agent modes:**
  - For Copilot: Copy `.github/chatmodes/*.chatmode.md` into your project’s `.github/chatmodes/`.
  - For Roo Code: Copy `.roomodes` to your project root.
- **Prompts and templates:** See `.github/prompts/` for reusable prompt templates for each mode.

## Project-Specific Conventions
- **All planning files (`.spec.md`, `.design.md`, `.tasks.md`) must live in the `.kodik` directory.**
- **ASCII art diagrams are mandatory** in design docs for UI and data flows.
- **Task lists must be extremely granular, with status, context, dependencies, and validation for each task.**
- **Agents must never deviate from the written plan**—no improvisation or skipping steps.
- **All status changes must be explicitly announced and reflected in the task file.**
- **Review plans** are generated using the `kodik review planner` mode, synthesizing spec, design, and tasks with a detailed checklist.

## Integration Points
- **GitHub Copilot:** Custom chat modes in `.github/chatmodes/`.
- **Roo Code:** Unified `.roomodes` file.
- **Install/update automation:** `scripts/install_update.sh` manages setup and updates for `.github` and `.roomodes`.

## Examples & References
- See `.github/chatmodes/` for mode definitions and usage patterns.
- See `.github/prompts/` for prompt templates.
- See `.roomodes` for Roo Code integration and detailed agent instructions.
- See `README.md` for high-level project overview and setup instructions.

---

**Agents: Always follow the mode-specific instructions in `.chatmode.md` or `.roomodes`. Never assume, always verify against the actual files.**

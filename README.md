# kodik
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

An open-source implementation of the Amazon Kiro Agent modes, designed and maintained by Nik Kyriakidis.

**k**iro **o**pen source**d** by n**ik**

`kodik` provides a structured, specification-driven development workflow through specialized agent modes. It enables you to replicate the spec-driven development behavior of Kiro.dev using GitHub Copilot, and is designed for easy porting to other platforms such as Roo Code or OpenCode.

## About The Project

`kodik` is a meta-tooling repository that defines agent modes for a highly organized, spec-driven software development lifecycle. Each mode is responsible for a distinct phase, from initial specification to final implementation, ensuring clarity, traceability, and automation throughout the process.

## Key Features

- **One-line installer**: Global CLI tool installation via curl
- Specification-driven workflow for software projects
- Five specialized agent modes: Spec, Design, Tasks, Agent, and Review Planner
- Markdown-based planning and execution documents
- Behavior-driven development (BDD) with Cucumber notation
- Automated, sequential task execution
- Backup and rollback protection for existing configurations
- Easily extensible and portable to other agent platforms

## Built With

- Markdown
- GitHub Copilot & Roo Code (for agent execution)
- [Kiro.dev](https://kiro.dev) (conceptual inspiration)

## Getting Started

`kodik` can be installed as a global CLI tool to easily set up agent modes in any project, or used directly from the repository.

### Prerequisites

- A compatible agent platform (such as GitHub Copilot, Roo Code, or OpenCode)
- Basic knowledge of Markdown and Git
- `curl` and `bash` (for CLI installation)

### Installation

#### Option 1: CLI Installation (Recommended)

Install `kodik` as a global command-line tool with a single command:

```sh
curl -fsSL https://raw.githubusercontent.com/nkyriakidis/kodik/main/kodik | sudo tee /usr/local/bin/kodik > /dev/null && sudo chmod +x /usr/local/bin/kodik
```

Then use it in any project directory:

```sh
cd your-project
kodik                    # Install both .github and .roomodes
kodik --github          # Install only GitHub Copilot modes
kodik --roo             # Install only Roo Code configuration
```

#### Option 2: Manual Installation

1. Clone the repo:
   ```sh
   git clone https://github.com/nkyriakidis/kodik.git
   ```
2. Navigate to the project directory:
   ```sh
   cd kodik
   ```
3. Run the install script directly:
   ```sh
   ./kodik
   ```

The script will automatically create the necessary directory structure and copy the agent mode definitions to your project.

## Usage

### Quick Start with CLI

After installing `kodik` globally, navigate to any project directory and run:

```sh
kodik  # Sets up both GitHub Copilot and Roo Code configurations
```

### CLI Options

```sh
kodik [OPTIONS] [MODE]

OPTIONS:
  -h, --help      Show usage information
  -f, --force     Skip backups and confirmations
  -d, --dry-run   Preview actions without execution

MODE:
  --github, -g    Install/update only .github directory
  --roo, -r       Install/update only .roomodes file
  (no mode)       Install/update both components (default)
```

### Using Agent Modes

Once installed, the agent modes will be available in your chosen platform:

- **GitHub Copilot**: Access custom chat modes via the mode dropdown in VS Code
- **Roo Code**: Use `/kodik-spec`, `/kodik-design`, `/kodik-tasks`, `/kodik-agent` commands

### Important: Update Your .gitignore

After running `kodik`, you should add the state directory to your `.gitignore` file to avoid committing backup and temporary files:

```sh
echo ".kodik-state/" >> .gitignore
```

Or manually add this line to your `.gitignore`:

```ignore
.kodik-state/
```

The `.kodik-state/` directory contains:
- **Backup files** of your previous configurations
- **Checksums** for tracking modifications
- **Installation history** and temporary files
- **Restoration logs** for rollback operations

This directory should not be committed to version control as it contains local state and backup data specific to your machine.

### Uninstalling

To remove the `kodik` CLI tool:

```sh
sudo rm /usr/local/bin/kodik
```

To remove installed configurations from a project:

```sh
rm -rf .github/chatmodes/kodik*
rm -f .roomodes
rm -rf .kodik-state  # Remove state/backup directory
```

## Agent Modes Available

The `kodik` system provides five specialized agent modes:

- **Spec Mode (`kodik spec`):** Generate specification documents with user stories and Cucumber scenarios
- **Design Mode (`kodik design`):** Create detailed technical design documents from specs  
- **Tasks Mode (`kodik tasks`):** Break down designs into granular, executable tasks
- **Agent Mode (`kodik agent`):** Execute tasks sequentially, updating status and validating completion
- **Review Planner Mode (`kodik review planner`):** Generate comprehensive review plans for features and changes

Each mode has specific instructions and workflows defined in the corresponding `.chatmode.md` files in `.github/chatmodes/`.

## Supported Providers

### GitHub Copilot
`kodik` supports **GitHub Copilot** through custom chat modes (`.chatmode.md` files) located in the `.github/chatmodes/` directory.

### Roo Code
`kodik` also supports **Roo Code** through a comprehensive configuration file (`.roomodes`) that defines all specialized modes.

### Future Provider Support
The mode definitions are designed to be easily portable to other AI coding assistants.

## How to Use kodik

### Quick Start (Recommended)

1. **Install kodik globally** (one-time setup):
   ```sh
   curl -fsSL https://raw.githubusercontent.com/nkyriakidis/kodik/main/kodik | sudo tee /usr/local/bin/kodik > /dev/null && sudo chmod +x /usr/local/bin/kodik
   ```

2. **Navigate to your project** and run:
   ```sh
   cd your-project
   kodik
   ```

This automatically creates the `.github/chatmodes/` directory with all chat mode definitions and the `.roomodes` file.

### Setting Up GitHub Copilot Chat Modes

1. **Automated setup** (creates the files automatically):
   ```sh
   kodik --github
   ```

2. **Manual setup** (if you prefer to customize):
   - Open Command Palette (`⇧⌘P` on Mac, `Ctrl+Shift+P` on Windows/Linux)
   - Run the command `Chat: New Mode File`
   - Choose "Workspace" to make the mode available to your project team
   - Copy mode definitions from this repository's `.github/chatmodes/` directory

3. **Default Location**: Chat modes are stored in `.github/chatmodes/` directory:
   ```
   your-project/
   └── .github/
       └── chatmodes/
           ├── kodik agent.chatmode.md
           ├── kodik design.chatmode.md
           ├── kodik review planner.chatmode.md
           ├── kodik spec.chatmode.md
           └── kodik tasks.chatmode.md
   ```

4. **Access the Modes**: In VS Code with GitHub Copilot:
   - Open the Chat view (`⌃⌘I` on Mac, `Ctrl+Alt+I` on Windows/Linux)
   - Use the **chat mode dropdown** at the top of the chat panel
   - Select your custom kodik modes from the dropdown list

### Setting Up Roo Code Modes

1. **Automated setup** (creates the file automatically):
   ```sh
   kodik --roo
   ```
   
   Or manually copy the `.roomodes` file from this repository to the root of your project directory.

2. **Access the Modes**: Roo Code automatically detects and loads the custom modes. Use them with the `/` prefix:
   - `/kodik-spec` for specification writing
   - `/kodik-design` for technical design
   - `/kodik-tasks` for task breakdown
   - `/kodik-agent` for automated implementation

### Best Prompts for Each Mode

#### 📝 Spec Mode (`kodik spec`)
**Purpose:** Initial feature requirements and user story creation

**Example Prompts:**
- "Create a spec for a user authentication system with login and password reset"
- "I need a specification for a shopping cart feature with add/remove items and checkout"
- "Write specs for a file upload system that supports drag and drop"

**Tips:**
- Be descriptive about the feature you want to build
- Mention any specific user types or roles involved
- Include any business rules or constraints
- Focus on what the feature should do, not how

#### 🏗️ Design Mode (`kodik design`)
**Purpose:** Creating technical architecture from existing specs

**Example Prompts:**
- "Create a design document for the user-authentication.spec.md"
- "Design the technical implementation for shopping-cart.spec.md"
- "Generate design for the file-upload-system.spec.md using our existing React components"

**Tips:**
- Always reference the specific `.spec.md` file you want designed
- Mention your tech stack if not obvious from the codebase
- Include any architectural constraints or preferences

#### ✅ Tasks Mode (`kodik tasks`)
**Purpose:** Breaking down designs into actionable tasks

**Example Prompts:**
- "Generate tasks for user-authentication based on the spec and design files"
- "Create a task list for implementing the shopping cart feature"
- "Break down the file-upload-system design into development tasks"

**Tips:**
- Reference both spec and design files for best results
- Tasks are created in the `.kodik` directory automatically
- Be specific about the feature name

#### 🤖 Agent Mode (`kodik agent`)
**Purpose:** Automated implementation execution

**Example Prompts:**
- "Implement the User Authentication System"
- "Continue working on the shopping cart implementation"
- "Execute the file upload system tasks"

**Tips:**
- Use the exact feature name from your task files
- The agent works through tasks sequentially
- Don't provide implementation details - it follows the pre-defined plan

#### 📋 Review Planner Mode (`kodik review planner`)
**Purpose:** Generate comprehensive review plans for features and changes

**Example Prompts:**
- "Create a review plan for the user authentication feature"
- "Generate review checklist for shopping cart implementation"
- "Plan review process for file upload system"

**Tips:**
- Reference existing spec, design, and task files
- Helps ensure thorough code review coverage
- Creates structured review checklists

### Recommended Workflow

1. **Start with Spec Mode** → Define what you want to build
2. **Move to Design Mode** → Plan how to build it technically  
3. **Use Tasks Mode** → Break it down into actionable steps
4. **Execute with Agent Mode** → Let the agent implement it automatically
5. **Use Review Planner Mode** → Create comprehensive review plans before merging

This creates a complete spec-driven development cycle that ensures thorough planning and review before implementation.

## Repository Structure

This repository contains the following key components:

- **`.github/chatmodes/`** - GitHub Copilot chat mode definitions
- **`.roomodes`** - Roo Code mode configuration file
- **`kodik`** - CLI installer script
- **`.kodik-state/`** - State directory (created after installation, should be gitignored)

## Development

This repository does not contain executable code or automated tests. It is a configuration and workflow definition resource for agent-powered development.

## Contributing

Contributions are welcome! If you have suggestions for new agent modes, improvements, or documentation, please open an issue or submit a pull request.

> **Note:** No `CONTRIBUTING.md` file exists yet. Please follow standard GitHub contribution practices.

## License

Distributed under the MIT License. See `LICENSE.md` for more information.

# kodik
An open-source implementation of the Amazon Kiro Agent modes, designed and maintained by Nik Kyriakidis.

**k**iro **o**pen source**d** by n**ik**

`kodik` provides a structured, specification-driven development workflow through specialized agent modes. It enables you to replicate the spec-driven development behavior of Kiro.dev using GitHub Copilot, and is designed for easy porting to other platforms such as Roo Code or OpenCode.

## About The Project

`kodik` is a meta-tooling repository that defines agent modes for a highly organized, spec-driven software development lifecycle. Each mode is responsible for a distinct phase, from initial specification to final implementation, ensuring clarity, traceability, and automation throughout the process.

## Key Features

- Specification-driven workflow for software projects
- Four specialized agent modes: Spec, Design, Tasks, Agent
- Markdown-based planning and execution documents
- Behavior-driven development (BDD) with Cucumber notation
- Automated, sequential task execution
- Easily extensible and portable to other agent platforms

## Built With

- Markdown
- GitHub Copilot & Roo Code (for agent execution)
- [Kiro.dev](https://kiro.dev) (conceptual inspiration)

## Getting Started

This repository does not contain executable code or dependencies. It provides agent mode definitions and workflow documentation for use with compatible agent platforms (e.g., GitHub Copilot, Roo Code, OpenCode).

### Prerequisites

- A compatible agent platform (such as GitHub Copilot)
- Basic knowledge of Markdown and Git

### Installation

1. Clone the repo:
   ```sh
   git clone https://github.com/nkyriakidis/kodik.git
   ```
2. Navigate to the project directory:
   ```sh
   cd kodik
   ```
3. Review the agent mode definitions in `.github/chatmodes/`

## Usage

To use the `kodik` agent modes:

1. Copy or reference the desired `.chatmode.md` files from `.github/chatmodes/` in your agent platform.
2. Follow the instructions in each mode file to implement the corresponding phase of your development workflow:
   - **Spec Mode:** Generate specification documents with user stories and Cucumber scenarios.
   - **Design Mode:** Create detailed technical design documents from specs.
   - **Tasks Mode:** Break down designs into granular, executable tasks.
   - **Agent Mode:** Execute tasks sequentially, updating status and validating completion.

See the documentation in each `.chatmode.md` file for detailed operational instructions.

## Supported Providers

### GitHub Copilot
`kodik` supports **GitHub Copilot** as an AI provider through custom chat modes (`.chatmode.md` files) located in the `.github/chatmodes/` directory.

### Roo Code
`kodik` also supports **Roo Code** through a comprehensive configuration file (`.roomodes`) that defines all four specialized modes.

### Future Provider Support
The mode definitions are designed to be easily portable to other AI coding assistants, such as **OpenCode** (planned).

## How to Use This

### Setting Up Custom Chat Modes in GitHub Copilot

1. **Create Chat Mode Files**: Use VS Code's built-in command to create the chat modes:
   - Open Command Palette (`⇧⌘P` on Mac, `Ctrl+Shift+P` on Windows/Linux)
   - Run the command `Chat: New Mode File`
   - Choose "Workspace" to make the mode available to your project team
   - Enter the mode name (e.g., "kodik-spec", "kodik-design", etc.)
2. **Copy Mode Definitions**: Replace the generated template with the content from the corresponding `.chatmode.md` files in this repository's `.github/chatmodes/` directory.
3. **Default Location**: By default, VS Code stores workspace chat modes in `.github/chatmodes/` directory:
   ```
   your-project/
   └── .github/
       └── chatmodes/
           ├── kodik-agent.chatmode.md
           ├── kodik-design.chatmode.md
           ├── kodik-spec.chatmode.md
           └── kodik-tasks.chatmode.md
   ```
4. **Access the Modes**: In VS Code with GitHub Copilot:
   - Open the Chat view (`⌃⌘I` on Mac, `Ctrl+Alt+I` on Windows/Linux)
   - Use the **chat mode dropdown** at the top of the chat panel
   - Select your custom kodik modes from the dropdown list

### Setting Up Custom Modes in Roo Code

1. **Copy the Configuration File**: Copy the `.roomodes` file from this repository to the root of your project directory.
2. **Verify Mode Registration**: Roo Code will automatically detect and load the custom modes defined in the `.roomodes` file when you open your project.
3. **Access the Modes**: In your Roo Code interface:
   - Look for the custom mode selector or dropdown
   - You should see the four kodik modes available:
     - `kodik-spec` - Specification writing mode
     - `kodik-design` - Technical design mode  
     - `kodik-tasks` - Task breakdown mode
     - `kodik-agent` - Agent execution mode
4. **Mode Usage**: Use the modes with the `/` prefix followed by the mode slug:
   - `/kodik-spec` for specification writing
   - `/kodik-design` for technical design
   - `/kodik-tasks` for task breakdown
   - `/kodik-agent` for automated implementation

### Best Prompts for Each Mode

#### 📝 Spec Mode (`kodik spec` in GitHub Copilot / `/kodik-spec` in Roo Coder)
**Best for:** Initial feature requirements and user story creation

**Example Prompts:**
- "Create a spec for a user authentication system with login and password reset"
- "I need a specification for a shopping cart feature with add/remove items and checkout"
- "Write specs for a file upload system that supports drag and drop"

**Prompt Tips:**
- Be descriptive about the feature you want to build
- Mention any specific user types or roles involved
- Include any business rules or constraints
- Don't worry about technical details - focus on what the feature should do

#### 🏗️ Design Mode (`kodik design` in GitHub Copilot / `/kodik-design` in Roo Coder)
**Best for:** Creating technical architecture from existing specs

**Example Prompts:**
- "Create a design document for the user-authentication.spec.md"
- "Design the technical implementation for shopping-cart.spec.md"
- "Generate design for the file-upload-system.spec.md using our existing React components"

**Prompt Tips:**
- Always reference the specific `.spec.md` file you want designed
- Mention your tech stack (React, Vue, Node.js, etc.) if not obvious from the codebase
- Include any architectural constraints or preferences
- The mode will analyze your existing codebase automatically

#### ✅ Tasks Mode (`kodik tasks` in GitHub Copilot / `/kodik-tasks` in Roo Coder)
**Best for:** Breaking down designs into actionable tasks

**Example Prompts:**
- "Generate tasks for user-authentication based on the spec and design files"
- "Create a task list for implementing the shopping cart feature"
- "Break down the file-upload-system design into development tasks"

**Prompt Tips:**
- Reference both the spec and design files for best results
- The mode will create tasks in a `.kodik` directory automatically
- Tasks will be ordered with dependencies and validation criteria
- Be specific about the feature name - it should match your spec/design files

#### 🤖 Agent Mode (`kodik agent` in GitHub Copilot / `/kodik-agent` in Roo Coder)
**Best for:** Automated implementation execution

**Example Prompts:**
- "Implement the User Authentication System" (will look for matching task files)
- "Continue working on the shopping cart implementation"
- "Execute the file upload system tasks"

**Prompt Tips:**
- Use the exact feature name from your task files
- The agent will work through tasks sequentially
- Don't provide implementation details - the agent follows the pre-defined plan
- The agent will announce each task as it starts and completes it

### Recommended Workflow

1. **Start with Spec Mode** → Define what you want to build
2. **Move to Design Mode** → Plan how to build it technically  
3. **Use Tasks Mode** → Break it down into actionable steps
4. **Execute with Agent Mode** → Let the agent implement it automatically

This creates a complete spec-driven development cycle that ensures thorough planning before implementation.

## Running Tests

This repository does not contain code or automated tests.

## Contributing

Contributions are welcome! If you have suggestions for new agent modes, improvements, or documentation, please open an issue or submit a pull request.

> **Note:** No `CONTRIBUTING.md` file exists yet. Please follow standard GitHub contribution practices.

## License

Distributed under the MIT License. See `LICENSE.md` for more information.

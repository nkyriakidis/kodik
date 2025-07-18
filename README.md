# kodik
An implementation of the Amazon Kiro Agent modes written by Nik.

**k**iro **o**pen source**d** by n**ik**

Currently this repo plays host to some Copilot Agent modes definitions that allows you to replicate the spec driven development behaviour of Kiro.dev using github copilot, these could be easily ported to Roo Code or OpenCode (i plan to later)

## Modes

kodik implements a structured, specification-driven development workflow through four specialized agent modes. Each mode handles a specific phase of the software development lifecycle, from initial specification to final implementation.

### 📝 Spec Mode (`spec.chatmode.md`)
**Purpose:** Creates comprehensive specification documents with user stories and Cucumber notation.

The Spec Mode agent analyzes user requests and workspace context to generate detailed `.spec.md` files containing:
- Comprehensive feature descriptions
- User stories from an end-user perspective
- Cucumber scenarios using Given-When-Then syntax
- Behavior-driven development (BDD) specifications

**Location:** `.github/chatmodes/spec.chatmode.md`

### 🏗️ Design Mode (`design.chatmode.md`)
**Purpose:** Translates specifications into detailed technical design documents.

The Design Mode agent takes specification files and creates comprehensive technical blueprints that include:
- Component breakdown and architecture decisions
- ASCII art UI/UX wireframes and component hierarchies
- Data flow diagrams showing system interactions
- API endpoint specifications
- Database schema changes
- Detailed implementation guidance

**Location:** `.github/chatmodes/design.chatmode.md`

### ✅ Tasks Mode (`tasks.chatmode.md`)
**Purpose:** Generates granular, executable task lists from spec and design documents.

The Tasks Mode agent breaks down technical designs into atomic, step-by-step implementation tasks:
- Highly granular task decomposition
- Clear dependencies between tasks
- Validation criteria for each task
- Status tracking (Pending/In Progress/Complete)
- Structured task lists saved to `.kodik` directory

**Location:** `.github/chatmodes/tasks.chatmode.md`

### 🤖 Agent Mode (`agent.chatmode.md`)
**Purpose:** Executes development tasks by strictly following predefined task lists.

The Agent Mode is the implementation executor that:
- Follows task lists sequentially without deviation
- Updates task status as work progresses
- Performs self-validation against defined criteria
- Ensures strict adherence to specifications and design
- Provides automated development execution

**Location:** `.github/chatmodes/agent.chatmode.md`

## Supported Providers

### GitHub Copilot
kodik supports **GitHub Copilot** as an AI provider through custom chat modes (`.chatmode.md` files) located in the `.github/chatmodes/` directory.

**Provider Files:**
- `.github/chatmodes/agent.chatmode.md` - Agent execution mode
- `.github/chatmodes/design.chatmode.md` - Technical design mode  
- `.github/chatmodes/spec.chatmode.md` - Specification writing mode
- `.github/chatmodes/tasks.chatmode.md` - Task breakdown mode

### Roo Code
kodik also supports **Roo Code** through a comprehensive configuration file that defines all four specialized modes.

**Configuration File:**
- `.roomodes` - Complete Roo Coder configuration with all four kodik modes (spec, design, tasks, agent)

The `.roomodes` file contains custom mode definitions that replicate the same structured workflow available in GitHub Copilot, allowing teams using Roo Coder to benefit from the complete kodik development methodology.

### Future Provider Support
The mode definitions are designed to be easily portable to other AI coding assistants:
- **OpenCode** (planned)

Each provider follows the same conceptual framework but is adapted to the specific tool's configuration format and capabilities.

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

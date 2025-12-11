# projects

CLI to organize and open local/SSH/WSL projects quickly.

![Go](https://github.com/filipenos/projects/workflows/Go/badge.svg)

## Quick install

```bash
go install github.com/filipenos/projects@latest
```

After installing, run `projects init` to generate the default config file.

## Core commands

| Command | Description | Notes |
| --- | --- | --- |
| `projects init` | Initialize new config file | Creates the default configuration. Alias: `i` |
| `projects create [name] [path]` | Registers a new project | Flags: `--editor` lets you edit fields before saving; `--no-validate` skips path checks |
| `projects update <name>` | Edits an existing project | Accepts `--no-validate` to update paths that do not exist yet |
| `projects delete <name>` | Deletes an existing project | Removes the project from the configuration |
| `projects list` | Lists all registered projects | Flags: `--ssh`, `--local`, `--workspace` filter by type (can be combined with AND logic) |
| `projects code <project>` | Opens the project in the configured editor | Uses aliases defined in `editors.json`. If the editor service fails, run `projects editors reload`. |
| `projects exec <project> <command...>` | Runs a command inside the project directory | Supports `local` and `ssh` projects (including workspaces) |
| `projects shell <project>` | Opens a shell inside the project | Supports `local`, `wsl` and `ssh` projects. Aliases: `sh`, `bash`, `zsh`, `nu`. For SSH, uses remote default shell. |
| `projects session <project> [args...]` | Opens/attaches a terminal session for the project | Aliases: `tmux`, `screen`. Use `--backend` to choose backend. Only supports local/WSL projects. |
| `projects editors ...` | Manages supported editors | `projects editors init/list/reload` |
| `projects completion [shell]` | Generates completion scripts | Use `--file` to write to disk instead of stdout |
| `projects version` | Shows version and commit information | Displays the current version of the CLI |

## Examples

Initialize the configuration:

```bash
projects init
```

Add the current directory as a project and review fields before saving:

```bash
projects create --editor
```

Add a project pointing to a path that does not exist yet:

```bash
projects create my-project /path/to/future/dir --no-validate
```

Update an existing project and allow saving even if the new path is missing:

```bash
projects update my-project --no-validate
```

Delete a project:

```bash
projects delete my-project
```

Run tests inside a local project:

```bash
projects exec my-project go test ./...
```

Execute a command on a remote SSH project:

```bash
projects exec my-ssh-project ls -la
```

Open a shell in a local project using a specific shell:

```bash
projects nu my-project    # Opens nushell
projects bash my-project  # Opens bash
projects shell my-project # Opens $SHELL
```

Open a shell on an SSH project:

```bash
projects shell my-ssh-project
# Note: SSH projects always use the remote server's default shell
```

Filter projects by type:

```bash
projects list --ssh              # List only SSH projects
projects list --local            # List only local/WSL projects
projects list --workspace        # List only workspace projects
projects list --ssh --workspace  # List SSH workspaces (AND logic)
```

Open (or reattach) a terminal session for the project:

```bash
projects session my-project              # Uses default backend (tmux)
projects tmux my-project                 # Explicitly use tmux
projects screen my-project               # Explicitly use screen
projects session --backend screen my-project  # Alternative syntax
```

Run backend-specific commands when creating the session:

```bash
projects tmux my-project split-window -h
# If the session already exists, the command fails telling you to close it first.
```

## Shell completions

```bash
# Bash
projects completion bash > /etc/bash_completion.d/projects

# Zsh
projects completion zsh > "${fpath[1]}/_projects"
```

Use `fish` or `powershell` to generate the respective scripts.

## SSH Projects

SSH projects use the VS Code Remote URI format: `vscode-remote://ssh-remote+HOST/PATH`

### Creating SSH projects

```bash
# Register an SSH project
projects create my-remote vscode-remote://ssh-remote+myserver/home/user/project

# Register an SSH workspace
projects create my-workspace vscode-remote://ssh-remote+myserver/home/user/project.code-workspace
```

### Using SSH projects

```bash
# Open shell (uses remote default shell)
projects shell my-remote

# Execute commands remotely
projects exec my-remote npm test

# Open in editor (if configured)
projects code my-remote
```

**Important notes:**
- Shell aliases (`nu`, `bash`, `zsh`) are ignored for SSH; the remote server's default shell is used
- Workspace files (`.code-workspace`) are automatically handled - the parent directory is used as working directory
- The host must be configured in your SSH config (`~/.ssh/config`)

## Custom editors

1. Run `projects editors init` to create `editors.json`.
2. Edit the file to add binaries such as `cursor`, `goland`, `code-insiders`.
3. Reload the configuration with `projects editors reload`.

**Note:** All configured editors automatically become available as command aliases. For example, if you configure `cursor`, you can use `projects cursor my-project` to open the project directly in Cursor.

To see all available editors and their status:

```bash
projects editors list
```

## Development

```bash
git clone https://github.com/filipenos/projects.git
cd projects
go test ./...
```

The repo already ships with a basic `Makefile`; run `go test` before submitting PRs describing new flags/behaviors.

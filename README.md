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
| `projects create [name] [path]` | Registers a new project | Flags: `--editor` lets you edit fields before saving; `--no-validate` skips path checks |
| `projects update <name>` | Edits an existing project | Accepts `--no-validate` to update paths that do not exist yet |
| `projects list` | Lists all registered projects | Shows type, workspace flag and path status |
| `projects code <project>` | Opens the project in the configured editor | Uses aliases defined in `editors.json`. If the editor service fails, run `projects editors reload`. |
| `projects exec <project> <command...>` | Runs a command inside the project directory | Currently supports `local` and `ssh` projects; other types return an explicit error |
| `projects shell <project>` | Opens a shell inside the project | Uses `$SHELL` or the alias used to call the command |
| `projects tmux <project> [args...]` | Opens/attaches a tmux session named after the project | If the session already exists and you pass tmux args, it aborts and asks you to close the session first |
| `projects editors ...` | Manages supported editors | `projects editors init/list/reload` |
| `projects completion [shell]` | Generates completion scripts | Use `--file` to write to disk instead of stdout |

## Examples

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

Run tests inside a local project:

```bash
projects exec my-project go test ./...
```

Open (or reattach) a tmux session dedicated to the project:

```bash
projects tmux my-project
```

Run a tmux command when creating the session:

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

## Custom editors

1. Run `projects editors init` to create `editors.json`.
2. Edit the file to add binaries such as `cursor`, `goland`, `code-insiders`.
3. Reload the configuration with `projects editors reload`.

## Development

```bash
git clone https://github.com/filipenos/projects.git
cd projects
go test ./...
```

The repo already ships with a basic `Makefile`; run `go test` before submitting PRs describing new flags/behaviors.

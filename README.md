# gitcraft

<div align="center">

```
  ▄▀  █ ▀█▀ ▄▀▀ █▀▄ ▄▀▄ █▀ ▀█▀
  ▀▄█ █  █  ▀▄▄ █▀▄ █▀█ █▀  █
```

**Craft perfect git commit messages with AI**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/nathannewyen/gitcraft)](https://github.com/nathannewyen/gitcraft/releases)

[Installation](#installation) • [Quick Start](#quick-start) • [Providers](#providers) • [Configuration](#configuration)

</div>

---

## Features

- **Zero Dependencies** - Single binary, no Node.js or Python required
- **Multi-Provider Support** - Works with OpenAI, Claude, Ollama, and Gemini
- **Conventional Commits** - Generates standardized commit messages automatically
- **Git Hooks** - Install as a prepare-commit-msg hook for seamless workflow
- **Offline Mode** - Use local models with Ollama for complete privacy
- **Interactive** - Review, edit, or regenerate before committing

## Installation

### Using Go

```bash
go install github.com/nathannewyen/gitcraft@latest
```

### From Source

```bash
git clone https://github.com/nathannewyen/gitcraft.git
cd gitcraft
go build -o gitcraft .
```

### Binary Releases

Download the latest binary from [Releases](https://github.com/nathannewyen/gitcraft/releases).

## Quick Start

### 1. Configure your API key

```bash
# For OpenAI (default)
gitcraft config set openai.api_key sk-your-api-key

# For Claude
gitcraft config set claude.api_key sk-ant-your-api-key
gitcraft config set provider claude

# For Gemini
gitcraft config set gemini.api_key your-api-key
gitcraft config set provider gemini

# For Ollama (no API key needed!)
gitcraft config set provider ollama
```

### 2. Stage your changes

```bash
git add .
```

### 3. Generate commit message

```bash
gitcraft
```

That's it! gitcraft will analyze your staged changes and generate a meaningful commit message.

## Usage

```bash
# Basic usage - generates and prompts for confirmation
gitcraft

# Specify commit type
gitcraft --type feat
gitcraft --type fix

# Add gitmoji
gitcraft --emoji

# Dry run - show message without committing
gitcraft --dry-run

# Use specific provider
gitcraft --provider claude
gitcraft --provider ollama
```

### Available Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--type` | `-t` | Commit type (feat, fix, docs, style, refactor, test, chore) |
| `--emoji` | `-e` | Add gitmoji to commit message |
| `--dry-run` | `-d` | Show message without committing |
| `--provider` | `-p` | AI provider to use |
| `--model` | `-m` | Model to use for generation |
| `--max-length` | `-l` | Maximum subject line length (default: 72) |

## Providers

### OpenAI (Default)

```bash
gitcraft config set openai.api_key sk-your-key
gitcraft config set openai.model gpt-4o-mini  # or gpt-4o, gpt-4-turbo
```

### Claude (Anthropic)

```bash
gitcraft config set claude.api_key sk-ant-your-key
gitcraft config set claude.model claude-3-5-sonnet-20241022
gitcraft config set provider claude
```

### Ollama (Local/Offline)

Perfect for offline use or when you want complete privacy:

```bash
# Make sure Ollama is running
ollama serve

# Pull a model
ollama pull llama3.2

# Configure gitcraft
gitcraft config set provider ollama
gitcraft config set ollama.model llama3.2
gitcraft config set ollama.url http://localhost:11434  # optional
```

### Gemini (Google)

```bash
gitcraft config set gemini.api_key your-key
gitcraft config set gemini.model gemini-1.5-flash
gitcraft config set provider gemini
```

## Configuration

Configuration is stored in `~/.gitcraft.yaml`

### Commands

```bash
# Set a value
gitcraft config set <key> <value>

# Get a value
gitcraft config get <key>

# List all config
gitcraft config list

# Show config path
gitcraft config path
```

### Example Configuration

```yaml
provider: openai
openai:
  api_key: sk-xxx
  model: gpt-4o-mini
claude:
  api_key: sk-ant-xxx
  model: claude-3-5-sonnet-20241022
ollama:
  url: http://localhost:11434
  model: llama3.2
gemini:
  api_key: xxx
  model: gemini-1.5-flash
```

### Environment Variables

You can also use environment variables with the `GITCRAFT_` prefix:

```bash
export GITCRAFT_OPENAI_API_KEY=sk-xxx
export GITCRAFT_PROVIDER=openai
```

## Git Hooks

Install gitcraft as a git hook to automatically generate commit messages:

```bash
# Install hook
gitcraft hook install

# Check status
gitcraft hook status

# Uninstall hook
gitcraft hook uninstall
```

Once installed, simply run `git commit` and gitcraft will generate the message for you!

## Commit Types

gitcraft follows the [Conventional Commits](https://www.conventionalcommits.org/) specification:

| Type | Description |
|------|-------------|
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation changes |
| `style` | Code style changes (formatting, etc.) |
| `refactor` | Code refactoring |
| `test` | Adding or updating tests |
| `chore` | Maintenance tasks |
| `perf` | Performance improvements |
| `ci` | CI/CD changes |
| `build` | Build system changes |

## Comparison

| Feature | gitcraft | aicommits | gptcommit |
|---------|----------|-----------|-----------|
| Language | Go | TypeScript | Rust |
| Zero Dependencies | Yes | No (Node.js) | No (Rust) |
| OpenAI | Yes | Yes | Yes |
| Claude | Yes | No | No |
| Ollama | Yes | Yes | No |
| Gemini | Yes | No | No |
| Git Hooks | Yes | Yes | Yes |
| Conventional Commits | Yes | Optional | No |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

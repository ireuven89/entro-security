
# Entro Security – GitHub Secrets Scanner API

Entro is a Go-based HTTP API server that scans all commits in a GitHub repository for secrets like API keys, access tokens, and credentials. It helps teams detect and respond to secret leaks early in development.

## 🔍 How It Works

- Accepts GitHub repository owner, repo name, and access token as input
- Calls the GitHub API to retrieve commit history
- Scans each commit for secret patterns
- Returns a list of suspected secrets via a RESTful API

## 🚀 Getting Started

### Prerequisites

- Go 1.18+
- GitHub Personal Access Token (with `repo` scope)

### Clone & Build

```bash
git clone https://github.com/ireuven89/entro-security.git
cd entro
go build -o entro .

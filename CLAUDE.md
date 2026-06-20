# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`did-this` is a CLI tool for tracking daily completed tasks, designed for remote workers who need to report progress for standups or status updates. Tasks are stored in BoltDB buckets organized by date (YYYY-MM-DD).

## Build & Development Commands

```bash
# Build the binary
go build -o did-this ./cmd/did-this

# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run the CLI directly (without building)
go run ./cmd/did-this [command]
```

## Architecture

### Core Components

- **Entry point**: `cmd/did-this/main.go` - Minimal main that delegates to commands package
- **Commands**: `internal/commands/` - All CLI commands and business logic
  - Uses `spf13/cobra` for command structure
  - Each command is registered in its `init()` function

### Database Model

- **Storage**: BoltDB (embedded key-value store)
- **Location**: `~/.did-this/did-this.db` by default
- **Bucket structure**: One bucket per date (format: `YYYY-MM-DD`)
  - Keys: Auto-incrementing sequence IDs (uint64, BigEndian encoded)
  - Values: Task description strings

### Configuration

- **Config file**: `~/.did-this/config.json`
- **Config struct** (`internal/commands/config.go`):
  - `Version`: Schema version (currently 1)
  - `DbPath`: Directory containing the database
  - `Db`: BoltDB connection handle (not serialized)

### Command Lifecycle

All commands follow this pattern via `rootCmd.PersistentPreRun` and `PersistentPostRun`:

1. Parse log level from `--loglevel` flag
2. Load config (creates config dir/file if missing)
3. Open BoltDB connection
4. Create bucket for current date if needed
5. Execute command
6. Save config and close DB

### Testing

Test files use a custom helper (`setupTestDB` in `internal/commands/testing_test.go`) that:
- Creates temporary database for test isolation
- Returns cleanup function to remove test artifacts
- Pattern: `defer cleanup()`

## Key Dependencies

- **cobra**: CLI framework
- **bbolt**: Embedded key-value database
- **logrus**: Structured logging
- **go-homedir**: Cross-platform home directory detection

## Notes

- No "backfilling" - tasks are always added to today's bucket
- Database timeout is 1 second to prevent lock contention
- Config is saved on every command execution (even if unchanged)

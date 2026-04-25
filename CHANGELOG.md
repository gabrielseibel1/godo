# Changelog

## v0.5.1

### Fixed
- **Overlapping periods** — `auto-work` and `auto-list` now correctly merge overlapping time periods instead of creating duplicates.

### Changed
- **Code cleanup** — removed unused `SplitN` assignment in `parseTimeRange`.

## v0.5.0

### Added
- **Cobra CLI framework** — replaces hand-rolled parser. All commands now have `--help`, proper error messages, and shell completions (`godo completion zsh|bash|fish|powershell`).
- **Shell completions** — activity ID autocompletion on commands that accept an ID (`get`, `delete`, `do`, `undo`, `work`, `tag`, `untag`).
- **Time period tracking** — `work` command accepts `HH:MM-HH:MM` time ranges in addition to durations. Periods are stored with start/end timestamps in `godo.json`.
- **`auto-list` command** — view logged hours from a given date. Supports `month` keyword to list from the 1st of the current month.
- **`auto-work` command** — combines `create` + `work` for quick logging. Auto-creates the activity from the date. Supports `--date`, `--yesterday`, and smart default (before 5am = yesterday).
- **`24:00` support** — time ranges like `18:00-24:00` are handled correctly.
- **Midnight crossing** — time ranges where end < start (e.g. `22:00-02:00`) automatically resolve to the next day.
- **Comprehensive test suite** — 56 tests covering every CLI command, verifying `godo.json` state after each operation.

### Changed
- **`work` command** — now accepts both duration (`3h`) and time range (`18:00-21:00`) formats. Optional `--date` flag for explicit date on time ranges.
- **JSON serialization** — `godo.json` now includes a `periods` array on each activity. Backward compatible: activities without periods still work.
- **File writer** — fixed truncation bug where shrinking JSON would leave stale bytes at the end of the file.

### Fixed
- **`slog.Error` calls** — fixed improper argument format in repository logger (bare `err` values replaced with proper key-value pairs).

### Removed
- Hand-rolled command parser (`parser.go`, `command.go`, `executor.go`, `errors.go`).
- Custom `help` command (replaced by cobra's built-in help).

## v0.4.0

Initial release with TUI and basic CLI commands: `create`, `list`, `get`, `delete`, `do`, `undo`, `work`, `tag`, `untag`, `cat`, `sublist`, `version`, `init`.

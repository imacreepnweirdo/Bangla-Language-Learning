# Bangla Language Learning Portal - Backend API

A Go-based REST API for a Bangla language learning portal that serves as an inventory of vocabulary, learning record store, and unified launchpad for learning apps.

## Features

- **Vocabulary Management**: Store and manage Bangla words with English translations
- **Thematic Groups**: Organize words into learning categories
- **Study Sessions**: Track learning activities and progress
- **Learning Analytics**: Record correct/incorrect answers and provide statistics
- **Dashboard**: Real-time progress tracking and quick stats

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: Gin (HTTP router)
- **Database**: SQLite3
- **Task Runner**: Mage (for database operations)

## Getting Started

### Prerequisites

- Go 1.25+ installed
- Mage task runner installed

### Installation

1. Clone the repository and navigate to the backend_go directory
2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Install Mage:
   ```bash
   go install github.com/magefile/mage@latest
   ```

### Database Setup

1. Initialize and migrate the database:
   ```bash
   mage build
   ```

   This runs the full setup:
   - Creates `words.db` SQLite database
   - Runs all migrations
   - Seeds initial data (if seed files exist in `db/seeds/`)

2. Individual operations:
   ```bash
   mage init      # Create database
   mage migrate   # Run migrations
   mage seed      # Seed data
   mage clean     # Remove database
   ```

### Running the Server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Health Check

```bash
curl http://localhost:8080/health
```

## API Endpoints

### Dashboard
- `GET /api/dashboard/last_study_session` - Get last study session info
- `GET /api/dashboard/study_progress` - Get study progress statistics
- `GET /api/dashboard/quick-stats` - Get quick dashboard stats

### Study Activities
- `GET /api/study_activities/:id` - Get study activity details
- `GET /api/study_activities/:id/study_sessions` - Get paginated study sessions
- `POST /api/study_activities` - Create new study activity

### Words
- `GET /api/words` - Get paginated list of words with stats
- `GET /api/words/:id` - Get word details with groups and stats

### Groups
- `GET /api/groups` - Get paginated list of groups with word counts
- `GET /api/groups/:id` - Get group details
- `GET /api/groups/:id/words` - Get paginated words in a group

### Study Sessions
- `POST /api/study_sessions/:id/words/:word_id/review` - Record word review

### System Management
- `POST /api/reset_history` - Reset study history (keeps words/groups)
- `POST /api/full_reset` - Full system reset (deletes everything)

## Database Schema

The application uses 6 main tables:

- `words` - Vocabulary words with Bengali/English translations
- `groups` - Thematic word groups
- `words_groups` - Many-to-many relationship between words and groups
- `study_sessions` - Study session records
- `study_activities` - Activity records linked to sessions
- `word_review_items` - Individual word practice records

## Project Structure

```
backend_go/
├── cmd/
│   └── server/           # Server entry point
├── internal/
│   ├── handlers/         # HTTP handlers by feature
│   ├── models/           # Data models and structs
│   └── services/         # Business logic (if needed)
├── db/
│   ├── migrations/       # SQL migration files
│   └── seeds/            # JSON seed files
├── magefile.go          # Mage task runner
├── magefile.mod         # Mage module file
├── go.mod               # Go modules
└── words.db             # SQLite database (created by mage)
```

## Development

### Adding New Endpoints

1. Define models in `internal/models/`
2. Create handlers in `internal/handlers/`
3. Register routes in `cmd/server/main.go`

### Database Migrations

1. Create new SQL file in `db/migrations/` with numbered prefix (e.g., `007_new_table.sql`)
2. Run `mage migrate` to apply

### Adding Seed Data

1. Create JSON files in `db/seeds/`
2. Run `mage seed` to import

## Testing

```bash
go test ./...
```

## License

This project is part of the Bangla Language Learning Portal prototype.

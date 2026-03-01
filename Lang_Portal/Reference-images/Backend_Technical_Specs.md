# Backend Server Technical Specs

# Business Goal:

A language learning school wants to build a prototype of
learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned
- Act as a Learning record store(LRS), providing correct and wrong
- A unified launchpad to launch different learning apps

## Technical Requirements

- The backend will be built using Go
- The database will be SQLite3
- The API will be built using Gin

## Directory Structure
backend-go/
|── cmd/
|   └── server/           # Server setup
├── internal/
│   ├── models/           # Data models
│   ├── handlers/         # HTTP handlers organized by feature (dashboard, words, groups, etc.)
│   └── services/         # Business logic
├── db/                   # Database-related files
│   ├── migrations/       # SQL migration files 
│   └── seeds/            # JSON seed files
├── magefile.mod          # Mage task runner 
├── go.mod                # Go modules file
└── words.db              # SQLite database (created by mage) 


## Database Schema

Our database will be a single SQLite database called 'words.db' that will be in the root of the project folder of 'backend-go'

We have the following tables:
- words - stored vocabulary words
  - id integer
  - bengali string
  - parts json
- words_groups - join table for words and groups (many-to-many)
  - id integer
  - word_id integer
  - group_id integer
- groups - thematic groups of words
  - id integer
  - name string
- study_sessions - records of study sessions grouping word_review_items
  - id integer
  - group_id integer
  - created_at datetime
  - study_activities_id integer
- study_activities - a specific study activity, linking a study session to group
  - id integer
  - study_session_id integer
  - group_id integer
  - created_at datetime
- word_review_items - a record of word practice, determining if the word was correct or not
  - id integer
  - word_id integer
  - study_session_id integer
  - correct boolean
  - created_at datetime

## API Endpoints

### GET /api/dashboard/last_study_session
Returns information about the last study session.

#### JSON Response
```json
{
  "id": 123,
  "group_id": 456,
  "created_at": "2026-02-28T10:30:00Z",
  "study_activities_id": 789,
  "group_name": "Basic Greetings"
}
```

### GET /api/dashboard/study_progress
Returns study progress statistics.
Please note that the frontend will determine progress bar based on total words studied and total available words.

#### JSON Response
```json
{
  "total_words_studied": 4,
  "total_available_words": 150
}
```

### GET /api/dashboard/quick-stats
Returns quick statistics for dashboard display.

#### JSON Response
```json
{
  "total_words": 500,
  "total_groups": 12,
  "words_learned": 150,
  "sessions_completed": 15,
  "current_streak": 5,
  "accuracy_rate": 80.0
}
```

### GET /api/study_activities/:id

#### JSON Response
```json
{
  "id": 1,
  "name": "Vocabulary Practice",
  "thumbnail_url": "https://example.com/thumbnail.jpg",
  "description": "Practice your vocabulary with flashcards"
}
```

### /api/study_activities/:id/study_sessions
- pagination with 100 items per page.

#### JSON Response
```json
{
  "items": [
    {
      "id": 5,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-28T10:25:00Z",
      "end_time": "2025-02-28T10:30:00Z",
      "review_items_count": 10
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 6,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### POST /api/study_activities

#### Request Params
- group_id: integer
- study_activity_id: integer

#### JSON Response
```json
{
  "id": 124,
  "group_id": 123
}
```

### GET /api/words
- pagination with 100 items per page.

#### JSON Response
```json
{
  "items": [
    {
      "bengali": "স্বাগতম",
      "parts_of_speech": "অব্যয়",
      "english": "welcome",
      "correct_count": "5",
      "wrong_count": "2"
    }
  ],
  "pagination": {
    "page": 1,
    "total_pages": 5,
    "total_items": 500,
    "items_per_page": 100
  }
}
```

### GET /api/words/:id

#### JSON Response
```json
{
  "id": 1,
  "bengali": "স্বাগতম",
  "english": "welcome",
  "stats": {
    "correct_count": 5,
    "wrong_count": 2
  },
  "items": [
    {
      "id": 1,
      "name": "Basic Greetings"
    }
  ]
}
```

### GET /api/groups
- pagination with 100 items per page.

#### JSON Response
```json
{
  "items": [
    {
      "id": 1,
      "name": "Basic Greetings",
      "words_count": 25
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 12,
    "items_per_page": 100
  }
}
```

### GET /api/groups/:id

#### JSON Response
```json
{
  "id": 1,
  "name": "Basic Greetings",
  "words": [],
  "stats": {
    "total_word_count": 25
  }
}
```

### GET /api/groups/:id/words
- pagination with 100 items per page.

#### JSON Response
```json
{
  "words": [
    {
      "id": 1,
      "bengali": "স্বাগতম",
      "english": "welcome",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 25,
    "items_per_page": 100
  }
}
```

### POST /api/reset_history

#### JSON Response
```json
{
  "success": true,
  "message": "Study history reset successfully"
}
```

### POST /api/full_reset

#### JSON Response
```json
{
  "success": true,
  "message": "System has been fully reset"
}
```

### POST /api/study_sessions/:id/words/:word_id/review

#### Request Params
- id (study_session_id) integer
- word_id integer
- correct boolean

#### Request Payload
```json
{
  "correct": true
}
```

#### JSON Response
```json
{
  "success": true,
  "word_id": 1,
  "study_session_id": 123,
  "correct": true,
  "created_at": "2026-02-28T10:31:00Z"
}
```

## Task Runner Tasks

Mage is a task runner for Go.
Let's list our possible tasks we need for our lang portal.


### Initialize Database

This task will initialize the sqlite database called `words.db`


### Migrate Database

This task will run a series of migrations sql files on the database

Migrations live in 'migrations' folder.
The migration files will be run in order of their file name.
The file names should look like this:

```sql
001_init.sql
002_create_words_table.sql
```

### Seed Database

This task will import json files and transform them into target data for our database

All seed files live in the 'seeds' folder
All seed files should be loaded.

In our task we should have DSL to specify each seed file and its expected group word name.

```json
[
  {
    "bengali": "বসন্ত",
    "part_of_speech": "বিশেষ্য",
    "spelling": "boshonto",
    "english": "spring"
  }
]
```
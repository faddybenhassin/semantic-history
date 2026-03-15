# Semantic History

A small semantic command history search system using vector embeddings and Qdrant.

- Python FastAPI service for indexing shell history commands and semantic search.
- Go CLI client that reads `~/.bash_history`, indexes commands, and queries for the best matching command.
- Designed as a minimal, end-to-end demo of semantic retrieval over command history.

## Architecture

1. `api/main.py`
   - FastAPI app with endpoints:
     - `POST /index`: index a command text into Qdrant as an embedding.
     - `GET /search`: embed query text, nearest-neighbor search in Qdrant.
   - Uses `SentenceTransformer('all-MiniLM-L6-v2')` for text embeddings.
2. `cli/main.go`
   - Reads local bash history at `~/.bash_history`.
   - Posts each command to `/index`.
   - Accepts one query argument, calls `/search`, and prints top semantic match if threshold passed.
3. `qdrant_storage/`
   - Existing Qdrant data layout (local persisted collection data).

## Prerequisites

- Docker (for Qdrant)
- Python 3.10+ and pip (for API)
- Go 1.20+ (for CLI)
- Network connectivity for model download and CLI/API communication

## Quick start

### 1) Start Qdrant

```bash
docker run -p 6333:6333 -v qdrant_storage:/qdrant/storage qdrant/qdrant:latest
```

### 2) Create Python venv and install dependencies

```bash
python -m venv .venv
source .venv/bin/activate   # Linux/Mac
.venv\Scripts\activate     # Windows PowerShell
pip install fastapi uvicorn sentence-transformers qdrant-client
```

### 3) Run API server

```bash
uvicorn api.main:app --host 0.0.0.0 --port 8000 --reload
```

The service will initialize `history` collection automatically on startup.

### 4) Run CLI indexing + search

```bash
go run ./cli/main.go "your semantic query"
```

Example:

```bash
go run ./cli/main.go "git push with tags"
```

If the best result is above the hardcoded threshold (0.5), it prints the matched command.

## Improvements to consider

- Add support for `~/.zsh_history`, cross-platform history paths, and path config flag.
- Add pipeline to avoid re-indexing duplicate entries.
- Store full command metadata (timestamp, shell, cwd) in payload.
- Add auth to API endpoints.
- Add tests for search ranking and endpoint error handling.

## API examples

### Index command

```bash
curl -X POST http://localhost:8000/index \
  -H 'Content-Type: application/json' \
  -d '{"command":"git commit -m \"fix bug\""}'
```

### Search query

```bash
curl "http://localhost:8000/search?query=commit+fix&limit=5"
```

## Notes

- This is a prototype; collection storage under `qdrant_storage/` may already contain data state for local testing.
- If using on Windows, adjust the history path lookup and shell history acquisition accordingly.


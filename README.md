# Go Moderation API

A Go API that integrates with OpenAI's moderation API to check content and cache results in MongoDB.

## Features

- Content moderation using OpenAI's moderation API
- MongoDB caching to reduce API calls and improve performance
- Source system tracking for each moderation request
- Environment-based configuration (.env for local development, OS env vars for production)

## Requirements

- Go 1.16+
- MongoDB
- OpenAI API key

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and update the values:
   ```
   cp .env.example .env
   ```
3. Update the MongoDB connection string and OpenAI API key in the `.env` file
4. Install dependencies:
   ```
   go mod tidy
   ```

## Running the API

### Local Development

```powershell
go run main.go
```

The API will be available at `http://localhost:8080`.

### Production Deployment (Railway)

This application is configured for deployment on Railway. Railway will automatically use the environment variables configured in your project settings.

## API Endpoints

### Health Check

```
GET /api/health
```

Returns a simple health check response.

### Moderate Content

```
POST /api/moderate
```

Request body:
```json
{
  "source_system": "your-app-name",
  "content": "Text to be moderated"
}
```

Response (200 OK - Content allowed):
```json
{
  "allowed": true,
  "message": "Content allowed"
}
```

Response (403 Forbidden - Content denied):
```json
{
  "allowed": false,
  "message": "Content violates content policy"
}
```

## Environment Variables

- `MONGO_URI`: MongoDB connection string
- `MONGO_DATABASE`: MongoDB database name (default: "moderation")
- `MONGO_COLLECTION`: MongoDB collection name (default: "results")
- `OPENAI_API_KEY`: OpenAI API key
- `PORT`: Server port (default: "8080")
- `GO_ENV`: Environment name ("production" for production mode)

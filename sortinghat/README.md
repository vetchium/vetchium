# Sorting Hat

A resume scoring service that evaluates the compatibility of resumes against job descriptions using AI models.

## Features

- HTTP endpoint for batch scoring of resumes against job descriptions
- Uses multiple AI models for comparison:
  - Sentence Transformers (all-MiniLM-L6-v2)
  - TF-IDF vectorizer
- Reads resume PDFs from S3/Minio storage
- Returns compatibility scores on a scale of 0-100

## API Endpoint

### `POST /score-batch`

Scores multiple resumes against a job description in a single batch request.

**Request:**

```json
{
  "job_description": "Full job description text...",
  "resume_paths": [
    "s3://bucket/path/to/resume1.pdf",
    "s3://bucket/path/to/resume2.pdf",
    "s3://bucket/path/to/resume3.pdf"
  ]
}
```

**Response:**

```json
{
  "scores": [
    {
      "application_id": "path/to/resume1.pdf",
      "model_scores": [
        {
          "model_name": "sentence-transformers-all-MiniLM-L6-v2",
          "score": 78
        },
        {
          "model_name": "tfidf-1.0",
          "score": 82
        }
      ]
    },
    {
      "application_id": "path/to/resume2.pdf",
      "model_scores": [
        {
          "model_name": "sentence-transformers-all-MiniLM-L6-v2",
          "score": 65
        },
        {
          "model_name": "tfidf-1.0",
          "score": 71
        }
      ]
    }
  ]
}
```

### `GET /health`

Health check endpoint that returns status of the service.

## Environment Variables

The service requires the following environment variables:

- `S3_ENDPOINT`: S3/Minio endpoint URL
- `S3_ACCESS_KEY`: S3/Minio access key
- `S3_SECRET_KEY`: S3/Minio secret key
- `S3_REGION`: S3 region (default: us-east-1)
- `PORT`: Port for the HTTP server (default: 8080)

## Development

1. Install dependencies: `pip install -r requirements.txt`
2. Run the server: `python main.py`

## Docker

Build the Docker image:

```bash
docker build -t vetchium/sortinghat -f sortinghat/Dockerfile .
```

## Kubernetes

The service is deployed to Kubernetes using configuration in `tilt-env/sortinghat.yaml`. 
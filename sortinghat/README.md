# Sorting Hat

A resume scoring service that evaluates the compatibility of resumes against job descriptions using AI models.

## Features

- HTTP endpoint for scoring resumes against job descriptions
- Uses multiple AI models for comparison:
  - Sentence Transformers (all-MiniLM-L6-v2)
  - spaCy (en_core_web_md)
- Reads resume PDFs from S3/Minio storage
- Returns compatibility scores on a scale of 0-100

## API Endpoint

### `GET /score-resumes-jd`

Scores a resume against a job description.

**Parameters:**

- `fileurl` (required): S3 URL to the resume PDF (format: `s3://bucket/path/to/file.pdf`)
- `job_description` (required): Job description text to compare against

**Response:**

```json
[
  {
    "resume": "s3://bucket/path/to/file.pdf",
    "compatibility_scores": {
      "sentence-transformers": 78.45,
      "spacy": 82.31
    }
  }
]
```

### `GET /health`

Health check endpoint that returns status of the service.

## Environment Variables

The service requires the following environment variables:

- `S3_ENDPOINT`: S3/Minio endpoint URL
- `S3_ACCESS_KEY`: S3/Minio access key
- `S3_SECRET_KEY`: S3/Minio secret key
- `S3_REGION`: S3 region (default: us-east-1)
- `S3_BUCKET`: Bucket name
- `PORT`: Port for the HTTP server (default: 8080)

## Development

1. Install dependencies: `pip install -r requirements.txt`
2. Download spaCy model: `python -m spacy download en_core_web_md`
3. Run the server: `python main.py`

## Docker

Build the Docker image:

```bash
docker build -t psankar/vetchi-sortinghat -f sortinghat/Dockerfile .
```

## Kubernetes

The service is deployed to Kubernetes using configuration in `tilt-env/sortinghat.yaml`. 
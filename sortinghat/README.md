# Sortinghat

AI-powered resume scoring service that evaluates resume-job description compatibility using sentence transformer models.

## What it does

- Accepts batch requests with resumes (S3 paths) and job descriptions
- Downloads PDFs from S3/Minio, extracts text using PyMuPDF
- Scores compatibility using 2 AI models from different vendors
- Returns scores 0-100 for each model

## Current Models

- **Microsoft E5-Base-v2**: `intfloat/e5-base-v2` (1.41GB)
- **Beijing Academy BGE-Base-v1.5**: `BAAI/bge-base-en-v1.5` (1.41GB)

## Docker Architecture

The service uses separate images:
- **Model images**: Each AI model in its own Docker image (~1.4GB each)
- **Runtime image**: Python app that loads models from shared volume (~800MB)

In Kubernetes, init containers copy models to shared storage before the main container starts.

## Making Changes

1. **Add/remove models**: Update `main.py` model loading and scoring
2. **Change scoring logic**: Modify `score_resume()` function
3. **Add endpoints**: Add FastAPI routes in `main.py`
4. **Dependencies**: Update `requirements.txt`

The service automatically rebuilds when you run `make dev` from the project root.

## Resource Requirements

- **Memory**: 4GB limit, 2GB request (for 2 models + runtime)
- **Storage**: 4GB volume for model storage
- **CPU**: 2 CPU limit, 200m request 
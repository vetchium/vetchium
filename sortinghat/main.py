#!/usr/bin/env python3

import os
import json
import logging
import sys
from typing import Dict, List, Any
import uvicorn
from fastapi import FastAPI, HTTPException, Query, Body
import boto3
from botocore.client import Config
import fitz  # PyMuPDF
from sentence_transformers import SentenceTransformer
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np

# Import TypeSpec-generated models
# Add the typespec directory to Python path - look for it relative to current file
current_dir = os.path.dirname(os.path.abspath(__file__))
typespec_path = os.path.join(os.path.dirname(current_dir), 'typespec')
if os.path.exists(typespec_path) and typespec_path not in sys.path:
    sys.path.insert(0, typespec_path)

try:
    from sortinghat import (
        ApplicationSortRequest,
        SortingHatRequest,
        ModelScore,
        SortingHatScore,
        SortingHatResponse,
    )
except ImportError as e:
    logger.error(f"Failed to import TypeSpec models: {e}")
    logger.error(f"Looked for typespec package at: {typespec_path}")
    raise

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger("sortinghat")

app = FastAPI(title="Sorting Hat", description="Resume scoring against job descriptions")

# Initialize model at startup
logger.info("Initializing sentence transformer model")
sentence_model = SentenceTransformer('all-MiniLM-L6-v2')
logger.info("Model initialization complete")

# S3/Minio client initialization
def get_s3_client():
    logger.debug(f"Connecting to S3 endpoint: {os.environ.get('S3_ENDPOINT')}")
    return boto3.client(
        's3',
        endpoint_url=os.environ.get('S3_ENDPOINT'),
        aws_access_key_id=os.environ.get('S3_ACCESS_KEY'),
        aws_secret_access_key=os.environ.get('S3_SECRET_KEY'),
        region_name=os.environ.get('S3_REGION', 'us-east-1'),
        config=Config(signature_version='s3v4')
    )

def extract_text_from_pdf(pdf_content: bytes) -> str:
    """Extract text from a PDF file"""
    text = ""
    try:
        logger.debug(f"Opening PDF document, size: {len(pdf_content)} bytes")
        with fitz.open(stream=pdf_content, filetype="pdf") as doc:
            logger.debug(f"PDF has {len(doc)} pages")
            for page in doc:
                text += page.get_text()
        logger.debug(f"Extracted {len(text)} characters of text from PDF")
    except Exception as e:
        logger.error(f"Error extracting text from PDF: {str(e)}")
        raise HTTPException(status_code=400, detail=f"Error extracting text from PDF: {str(e)}")
    return text

def score_with_sentence_transformers(resume_text: str, job_description: str) -> float:
    """Score resume against job description using Sentence Transformer"""
    logger.debug("Generating embeddings with sentence-transformers model")
    resume_embedding = sentence_model.encode(resume_text)
    jd_embedding = sentence_model.encode(job_description)

    # Calculate cosine similarity
    logger.debug("Calculating cosine similarity between embeddings")
    similarity = cosine_similarity(
        [resume_embedding],
        [jd_embedding]
    )[0][0]

    # Convert similarity score (typically -1 to 1) to 0-100 scale
    score = max(0, min(100, (similarity + 1) * 50))
    logger.debug(f"Sentence-transformer raw similarity: {similarity}, scaled score: {score}")
    return score

def score_with_tfidf(resume_text: str, job_description: str) -> float:
    """Score resume against job description using TF-IDF vectorization"""
    logger.debug("Generating TF-IDF vectors")
    # Create TF-IDF vectorizer
    vectorizer = TfidfVectorizer(stop_words='english')

    # Create document-term matrix
    texts = [resume_text, job_description]
    tfidf_matrix = vectorizer.fit_transform(texts)

    # Calculate cosine similarity
    logger.debug("Calculating cosine similarity between TF-IDF vectors")
    similarity = cosine_similarity(tfidf_matrix[0:1], tfidf_matrix[1:2])[0][0]

    # Convert to 0-100 scale (same transformation as transformer model)
    score = max(0, min(100, (similarity + 1) * 50))
    logger.debug(f"TF-IDF raw similarity: {similarity}, scaled score: {score}")
    return score

def score_single_resume(application_id: str, fileurl: str, job_description: str) -> SortingHatScore:
    """Score a single resume against a job description"""
    logger.info(f"Processing resume: {fileurl} for application {application_id}")
    
    try:
        # Parse S3 URL to extract bucket and key
        if not fileurl.startswith("s3://"):
            logger.error(f"Invalid S3 URI format: {fileurl}")
            raise HTTPException(status_code=400, detail="Invalid S3 URI format. Must start with s3://")

        parts = fileurl[5:].split('/', 1)
        if len(parts) < 2:
            logger.error(f"Invalid S3 URI format (missing key): {fileurl}")
            raise HTTPException(status_code=400, detail="Invalid S3 URI format. Must be s3://bucket/key")

        bucket = parts[0]
        key = parts[1]
        logger.info(f"Parsed S3 location - Bucket: {bucket}, Key: {key}")

        # Get S3 client
        s3_client = get_s3_client()

        # Download PDF from S3
        try:
            logger.info(f"Downloading PDF from S3: {bucket}/{key}")
            response = s3_client.get_object(Bucket=bucket, Key=key)
            pdf_content = response['Body'].read()
            logger.info(f"Successfully downloaded PDF, size: {len(pdf_content)} bytes")
        except Exception as e:
            logger.error(f"Error retrieving file from S3: {str(e)}")
            raise HTTPException(status_code=404, detail=f"Error retrieving file from S3: {str(e)}")

        # Extract text from PDF
        logger.info("Extracting text from PDF")
        resume_text = extract_text_from_pdf(pdf_content)
        logger.info(f"Text extraction complete. Extracted {len(resume_text)} characters")

        # Score resume with different models
        logger.info("Starting resume scoring")
        sbert_score = score_with_sentence_transformers(resume_text, job_description)
        tfidf_score = score_with_tfidf(resume_text, job_description)
        logger.info(f"Scoring complete - Sentence Transformer: {sbert_score:.2f}, TF-IDF: {tfidf_score:.2f}")

        # Prepare model scores using TypeSpec models
        model_scores = [
            ModelScore(model_name="sentence-transformers-all-MiniLM-L6-v2", score=round(sbert_score)),
            ModelScore(model_name="tfidf-1.0", score=round(tfidf_score))
        ]
        
        return SortingHatScore(
            application_id=application_id,
            model_scores=model_scores
        )
        
    except HTTPException:
        raise
    except Exception as e:
        logger.exception(f"Unexpected error while processing resume: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")

@app.post("/score-batch", response_model=SortingHatResponse)
async def score_batch(request: SortingHatRequest = Body(...)):
    """Score multiple resumes against a job description in a single batch"""
    logger.info(f"Processing batch of {len(request.application_sort_requests)} resumes")
    logger.info(f"Job description length: {len(request.job_description)} characters")
    
    scores = []
    
    # Process each resume in the batch
    for app_request in request.application_sort_requests:
        try:
            score_result = score_single_resume(
                app_request.application_id,
                app_request.resume_path, 
                request.job_description
            )
            scores.append(score_result)
        except HTTPException as e:
            logger.error(f"Error processing resume {app_request.resume_path}: {e.detail}")
            # Continue with other resumes even if one fails
            continue
        except Exception as e:
            logger.exception(f"Unexpected error processing resume {app_request.resume_path}: {str(e)}")
            continue
    
    logger.info(f"Completed batch processing. Scored {len(scores)} out of {len(request.application_sort_requests)} resumes")
    return SortingHatResponse(scores=scores)

@app.get("/health")
async def health_check():
    logger.debug("Health check endpoint called")
    return {"status": "healthy"}

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8080))
    logger.info(f"Starting Sorting Hat server on port {port}")
    uvicorn.run(app, host="0.0.0.0", port=port)
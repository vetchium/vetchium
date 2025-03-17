#!/usr/bin/env python3

import os
import json
from typing import Dict, List, Any
import uvicorn
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel
import boto3
from botocore.client import Config
import fitz  # PyMuPDF
from sentence_transformers import SentenceTransformer
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np

app = FastAPI(title="Sorting Hat", description="Resume scoring against job descriptions")

# Initialize model at startup
sentence_model = SentenceTransformer('all-MiniLM-L6-v2')

# S3/Minio client initialization
def get_s3_client():
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
        with fitz.open(stream=pdf_content, filetype="pdf") as doc:
            for page in doc:
                text += page.get_text()
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Error extracting text from PDF: {str(e)}")
    return text

def score_with_sentence_transformers(resume_text: str, job_description: str) -> float:
    """Score resume against job description using Sentence Transformer"""
    resume_embedding = sentence_model.encode(resume_text)
    jd_embedding = sentence_model.encode(job_description)

    # Calculate cosine similarity
    similarity = cosine_similarity(
        [resume_embedding],
        [jd_embedding]
    )[0][0]

    # Convert similarity score (typically -1 to 1) to 0-100 scale
    score = max(0, min(100, (similarity + 1) * 50))
    return score

def score_with_tfidf(resume_text: str, job_description: str) -> float:
    """Score resume against job description using TF-IDF vectorization"""
    # Create TF-IDF vectorizer
    vectorizer = TfidfVectorizer(stop_words='english')

    # Create document-term matrix
    texts = [resume_text, job_description]
    tfidf_matrix = vectorizer.fit_transform(texts)

    # Calculate cosine similarity
    similarity = cosine_similarity(tfidf_matrix[0:1], tfidf_matrix[1:2])[0][0]

    # Convert to 0-100 scale
    score = similarity * 100
    return score

class ScoringResponse(BaseModel):
    resume: str
    compatibility_scores: Dict[str, float]

@app.get("/score-resumes-jd", response_model=List[ScoringResponse])
async def score_resumes(
    fileurl: str = Query(..., description="S3 URI to the resume PDF"),
    job_description: str = Query(..., description="Job description to compare against")
):
    try:
        # Parse S3 URL to extract bucket and key
        if not fileurl.startswith("s3://"):
            raise HTTPException(status_code=400, detail="Invalid S3 URI format. Must start with s3://")

        parts = fileurl[5:].split('/', 1)
        if len(parts) < 2:
            raise HTTPException(status_code=400, detail="Invalid S3 URI format. Must be s3://bucket/key")

        bucket = parts[0]
        key = parts[1]

        # Get S3 client
        s3_client = get_s3_client()

        # Download PDF from S3
        try:
            response = s3_client.get_object(Bucket=bucket, Key=key)
            pdf_content = response['Body'].read()
        except Exception as e:
            raise HTTPException(status_code=404, detail=f"Error retrieving file from S3: {str(e)}")

        # Extract text from PDF
        resume_text = extract_text_from_pdf(pdf_content)

        # Score resume with different models
        sbert_score = score_with_sentence_transformers(resume_text, job_description)
        tfidf_score = score_with_tfidf(resume_text, job_description)

        # Prepare response
        result = ScoringResponse(
            resume=fileurl,
            compatibility_scores={
                "sentence-transformers": round(sbert_score, 2),
                "tfidf": round(tfidf_score, 2)
            }
        )

        return [result]

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8080))
    uvicorn.run(app, host="0.0.0.0", port=port)
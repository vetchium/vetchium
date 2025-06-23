#!/usr/bin/env python3

import os
import logging
import sys
from typing import List
import uvicorn
from fastapi import FastAPI, HTTPException, Body
import boto3
import fitz  # PyMuPDF
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

# Import TypeSpec-generated models
current_dir = os.path.dirname(os.path.abspath(__file__))
typespec_path = os.path.join(os.path.dirname(current_dir), 'typespec')
if os.path.exists(typespec_path) and typespec_path not in sys.path:
    sys.path.insert(0, typespec_path)

from sortinghat import (
    SortingHatRequest,
    ModelScore,
    SortingHatScore,
    SortingHatResponse,
)

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("sortinghat")

app = FastAPI(title="Sorting Hat", description="AI Resume Scoring Service")

# Two different vendors for comprehensive scoring approaches
logger.info("Loading AI models...")
microsoft_model = SentenceTransformer('intfloat/e5-base-v2')          # Microsoft Research
beijing_model = SentenceTransformer('BAAI/bge-base-en-v1.5')         # Beijing Academy of AI
logger.info("Models loaded successfully")

# S3 client setup
s3_client = boto3.client(
    's3',
    endpoint_url=os.environ.get('S3_ENDPOINT'),
    aws_access_key_id=os.environ.get('S3_ACCESS_KEY'),
    aws_secret_access_key=os.environ.get('S3_SECRET_KEY'),
    region_name=os.environ.get('S3_REGION', 'us-east-1'),
)

def download_resume(fileurl: str) -> str:
    """Download and extract text from resume PDF"""
    if not fileurl.startswith("s3://"):
        raise HTTPException(400, "Invalid S3 URL format")
    
    # Parse S3 URL: s3://bucket/key
    bucket, key = fileurl[5:].split('/', 1)
    
    try:
        # Download PDF from S3
        response = s3_client.get_object(Bucket=bucket, Key=key)
        pdf_content = response['Body'].read()
        
        # Extract text using PyMuPDF
        text = ""
        with fitz.open(stream=pdf_content, filetype="pdf") as doc:
            for page in doc:
                text += page.get_text() + "\n"
        
        return text.strip()
    
    except Exception as e:
        logger.error(f"Failed to process resume {fileurl}: {e}")
        raise HTTPException(404, f"Could not process resume: {str(e)}")

def score_resume(resume_text: str, job_description: str) -> List[ModelScore]:
    """Score resume using two AI models from different vendors"""
    try:
        # Generate embeddings using two different vendor approaches
        resume_microsoft = microsoft_model.encode(resume_text)
        job_microsoft = microsoft_model.encode(job_description)
        
        resume_beijing = beijing_model.encode(resume_text)
        job_beijing = beijing_model.encode(job_description)
        
        # Calculate similarities with different vendor models
        microsoft_similarity = cosine_similarity([resume_microsoft], [job_microsoft])[0][0]
        beijing_similarity = cosine_similarity([resume_beijing], [job_beijing])[0][0]
        
        # Convert to 0-100 scale
        microsoft_score = max(0, min(100, (microsoft_similarity + 1) * 50))
        beijing_score = max(0, min(100, (beijing_similarity + 1) * 50))
        
        return [
            ModelScore(model_name="Microsoft-E5-Research", score=round(microsoft_score)),
            ModelScore(model_name="Beijing-Academy-BGE", score=round(beijing_score))
        ]
    
    except Exception as e:
        logger.error(f"Scoring failed: {e}")
        raise HTTPException(500, f"Scoring error: {str(e)}")

@app.post("/score-batch", response_model=SortingHatResponse)
async def score_batch(request: SortingHatRequest = Body(...)):
    """Score multiple resumes against a job description"""
    logger.info(f"Processing {len(request.application_sort_requests)} resumes")
    
    scores = []
    
    for app_request in request.application_sort_requests:
        try:
            # Download and extract resume text
            resume_text = download_resume(app_request.resume_path)
            
            # Score resume against job description
            model_scores = score_resume(resume_text, request.job_description)
            
            scores.append(SortingHatScore(
                application_id=app_request.application_id,
                model_scores=model_scores
            ))
            
        except HTTPException:
            # Skip failed resumes, continue processing others
            logger.warning(f"Skipping failed resume: {app_request.resume_path}")
            continue
    
    logger.info(f"Successfully scored {len(scores)} resumes")
    return SortingHatResponse(scores=scores)

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "models": 2}

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8080))
    logger.info(f"Starting Sorting Hat on port {port}")
    uvicorn.run(app, host="0.0.0.0", port=port)
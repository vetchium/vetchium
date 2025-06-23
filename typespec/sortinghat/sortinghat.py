from typing import List
from pydantic import BaseModel, Field


class ApplicationSortRequest(BaseModel):
    """Request to score a single application's resume"""
    application_id: str = Field(..., description="The unique identifier for the application")
    resume_path: str = Field(..., description="S3 path to the resume file in format s3://bucket/key")


class SortingHatRequest(BaseModel):
    """Request to score multiple resumes against a job description in a batch"""
    job_description: str = Field(..., description="The job description to score resumes against")
    application_sort_requests: List[ApplicationSortRequest] = Field(..., description="List of applications to score")


class ModelScore(BaseModel):
    """Score from a specific model"""
    model_name: str = Field(..., description="Name of the model that generated the score")
    score: int = Field(..., ge=0, le=100, description="Score value from 0 to 100")


class SortingHatScore(BaseModel):
    """Scores for a single application from all models"""
    application_id: str = Field(..., description="The application ID this score relates to")
    model_scores: List[ModelScore] = Field(..., description="Scores from different models")


class SortingHatResponse(BaseModel):
    """Response containing scores for all applications in the batch"""
    scores: List[SortingHatScore] = Field(..., description="List of application scores")


# Define what should be exported when using 'from sortinghat import *'
__all__ = [
    "ApplicationSortRequest",
    "SortingHatRequest",
    "ModelScore", 
    "SortingHatScore",
    "SortingHatResponse",
]

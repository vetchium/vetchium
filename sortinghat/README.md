# Sorting Hat

A resume scoring service that evaluates the compatibility of resumes against job descriptions using multiple AI models from different vendors.

## Features

- HTTP endpoint for batch scoring of resumes against job descriptions
- Uses multiple AI models from different vendors for diverse scoring:
  - **Microsoft Research**: E5-Base-v2 model (intfloat/e5-base-v2)
  - **Beijing Academy of AI**: BGE-Base-v1.5 model (BAAI/bge-base-en-v1.5)
  - **Meta AI Research**: Contriever-MSMARCO model (facebook/contriever-msmarco)
- Modular architecture with individual model images for scalability
- Reads resume PDFs from S3/Minio storage
- Returns compatibility scores on a scale of 0-100

## API Endpoint

### `POST /score-batch`

Scores multiple resumes against a job description in a single batch request.

**Request:**

```json
{
  "job_description": "Full job description text...",
  "application_sort_requests": [
    {
      "application_id": "app_123",
      "resume_path": "s3://bucket/path/to/resume1.pdf"
    },
    {
      "application_id": "app_456", 
      "resume_path": "s3://bucket/path/to/resume2.pdf"
    }
  ]
}
```

**Response:**

```json
{
  "scores": [
    {
      "application_id": "app_123",
      "model_scores": [
        {
          "model_name": "Microsoft-E5-Base-v2",
          "score": 78
        },
        {
          "model_name": "Beijing-BGE-Base-v1.5",
          "score": 82
        },
        {
          "model_name": "Meta-Contriever-MSMARCO",
          "score": 75
        }
      ]
    },
    {
      "application_id": "app_456",
      "model_scores": [
        {
          "model_name": "Microsoft-E5-Base-v2",
          "score": 65
        },
        {
          "model_name": "Beijing-BGE-Base-v1.5",
          "score": 71
        },
        {
          "model_name": "Meta-Contriever-MSMARCO",
          "score": 68
        }
      ]
    }
  ]
}
```

### `GET /health`

Health check endpoint that returns service status and number of loaded models.

## Architecture

The service uses a modular architecture with separate Docker images for each AI model:

```
┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐
│ E5-Base-v2 Model    │  │ BGE-Base-v1.5 Model │  │ Contriever-MSMARCO  │
│ Image (440MB)       │  │ Image (440MB)       │  │ Model Image (440MB) │
└─────────────────────┘  └─────────────────────┘  └─────────────────────┘
           │                        │                        │
           └──────────┬─────────────┼────────────────────────┘
                      │             │
           ┌─────────────────────────────────┐
           │ Runtime Image (800MB)           │
           └─────────────────────────────────┘
```

This design enables:
- **Scalability**: Add new models without image size explosion
- **Flexibility**: Mix and match models as needed
- **Performance**: Parallel model loading and better caching

## Adding New Models

To add a new AI model to the service:

### 1. Create Model Dockerfile

Create `sortinghat/Dockerfile.model-{model-name}`:

```dockerfile
# Example: sortinghat/Dockerfile.model-jina-embeddings-v2
FROM python:3.10-slim as model-downloader

ENV HF_HOME=/models/huggingface \
    TRANSFORMERS_CACHE=/models/huggingface \
    SENTENCE_TRANSFORMERS_HOME=/models/sentence-transformers

RUN pip install --no-cache-dir sentence-transformers>=2.2.0
RUN mkdir -p /models/huggingface /models/sentence-transformers

# Download the new model
RUN python -c "from sentence_transformers import SentenceTransformer; SentenceTransformer('jinaai/jina-embeddings-v2-base-en')"

FROM scratch
COPY --from=model-downloader /models /models
COPY --from=model-downloader /root/.cache /root/.cache
```

### 2. Update Build Configuration

**Tiltfile:**
```python
docker_build('vetchium/sortinghat-model-jina-embeddings-v2', '.', dockerfile='sortinghat/Dockerfile.model-jina-embeddings-v2')
```

**Makefile:**
```makefile
# Add to docker target:
docker buildx build -f sortinghat/Dockerfile.model-jina-embeddings-v2 \
    -t ghcr.io/vetchium/sortinghat-model-jina-embeddings-v2:$(GIT_SHA) \
    --platform=linux/amd64,linux/arm64 .

# Add to publish target with --push flag
```

### 3. Update Kubernetes Deployment

**Add init container to YAML files:**
```yaml
- name: jina-embeddings-v2-model-downloader
  image: vetchium/sortinghat-model-jina-embeddings-v2
  command: ['sh', '-c', 'cp -r /models/* /shared-models/']
  volumeMounts:
    - name: model-storage
      mountPath: /shared-models
  resources:
    requests:
      memory: "256Mi"
      cpu: "100m"
    limits:
      memory: "512Mi"
      cpu: "500m"
```

**Update deployment metadata:**
```yaml
# Add to deployment annotations in both tilt-env/ and helm templates
annotations:
  model-name/jina-embeddings-v2: "jinaai/jina-embeddings-v2-base-en"
  image-name/jina-embeddings-v2: "vetchium/sortinghat-model-jina-embeddings-v2:latest"
  deployment/model-count: "4"  # Update count
```

### 4. Update Application Code

**Add model to `main.py`:**
```python
# Load the new model
jina_model = SentenceTransformer('jinaai/jina-embeddings-v2-base-en')

# Add scoring in score_resume function
jina_embedding = jina_model.encode(resume_text)
# ... calculate similarity and add to results
```

## Recommended Models for Resume-JD Matching

### Currently Deployed
- **Microsoft E5-Base-v2**: Fast, optimized for text retrieval (intfloat/e5-base-v2)
- **Beijing Academy BGE-Base-v1.5**: High accuracy, general-purpose embeddings (BAAI/bge-base-en-v1.5)
- **Meta AI Contriever-MSMARCO**: Trained on MS MARCO dataset for document retrieval (facebook/contriever-msmarco)

### Excellent Additions (Commercial-friendly)

#### 1. Jina AI (Recommended)
```python
# High-quality embeddings from AI search company
SentenceTransformer('jinaai/jina-embeddings-v2-base-en')
```
- **License**: Apache 2.0 ✅
- **Size**: ~560MB
- **Specialty**: Document retrieval and search
- **Company**: Jina AI (Neural search specialists)

#### 2. Salesforce Research
```python
# Enterprise-focused embeddings
SentenceTransformer('Salesforce/SFR-Embedding-Mistral')
```
- **License**: Apache 2.0 ✅  
- **Size**: ~400MB
- **Specialty**: Business document understanding
- **Company**: Salesforce Research

#### 3. Cohere AI
```python
# Multilingual capabilities
SentenceTransformer('Cohere/Cohere-embed-english-v3.0')
```
- **License**: Apache 2.0 ✅
- **Size**: ~500MB
- **Specialty**: Professional text embeddings
- **Company**: Cohere AI

#### 4. Alibaba DAMO Academy
```python
# Research-grade performance
SentenceTransformer('BAAI/bge-large-en-v1.5')
```
- **License**: Apache 2.0 ✅
- **Size**: ~1.2GB
- **Specialty**: High-accuracy embeddings
- **Company**: Alibaba Research

## Environment Variables

- `S3_ENDPOINT`: S3/Minio endpoint URL
- `S3_ACCESS_KEY`: S3/Minio access key
- `S3_SECRET_KEY`: S3/Minio secret key
- `S3_REGION`: S3 region (default: us-east-1)
- `S3_BUCKET`: S3 bucket name
- `PORT`: Port for the HTTP server (default: 8080)
- `HF_HOME`: HuggingFace cache directory (set automatically)

## Resource Requirements

### Memory Requirements
The sortinghat service loads 3 AI models simultaneously, requiring significant memory:

- **E5-Base-v2**: ~1.4GB in memory
- **BGE-Base-v1.5**: ~1.4GB in memory  
- **Contriever-MSMARCO**: ~6.8GB in memory
- **Total Models**: ~9.6GB
- **Python Runtime**: ~200-300MB
- **Processing Buffer**: ~1-2GB during inference
- **Recommended**: 12GB memory limit, 4GB request

### Storage Requirements
- **Model Storage Volume**: 12GB (EmptyDir for 3 models + copy buffer)
- **Individual Model Images**: 
  - E5-Base-v2: 1.41GB
  - BGE-Base-v1.5: 1.41GB
  - Contriever-MSMARCO: 6.76GB
- **Runtime Image**: ~800MB (Python + dependencies)

### CPU Requirements
- **Model Loading**: CPU-intensive during startup
- **Inference**: Moderate CPU usage
- **Recommended**: 3 CPU limit, 200m request

## Development

1. Install dependencies: `pip install -r requirements.txt`
2. Set environment variables for S3 access
3. Run the server: `python main.py`

## Docker & Kubernetes

The service automatically builds and deploys using:
- **Development**: `make dev` (Tilt with live reload)
- **Production**: `make publish` (Multi-platform images to registry)
- **Testing**: `make devtest VMUSER=test VMADDR=<ip>` (Remote deployment)

All model images are built automatically and deployed as init containers in Kubernetes.

### Version Tracking & Metadata

Each deployment includes comprehensive metadata for operational visibility:

**Check deployed model versions:**
```bash
# View all model versions and images
kubectl get deployment sortinghat -o jsonpath='{.metadata.annotations}' | jq

# Get specific model version
kubectl get deployment sortinghat -n vetchium-dev \
  -o jsonpath='{.metadata.annotations.model-name/e5-base-v2}'

# Check all image tags
kubectl get deployment sortinghat -n vetchium-dev \
  -o jsonpath='{.metadata.annotations}' | jq 'to_entries | map(select(.key | startswith("image-name/")))'
```

**Available metadata annotations:**
- `model-name/*`: HuggingFace model identifiers (e.g., `intfloat/e5-base-v2`)
- `image-name/*`: Full Docker image names with tags
- `deployment/model-count`: Number of AI models deployed
- `deployment/created-by`: Deployment method (helm/tilt-dev)
- `deployment/deployment-time`: ISO timestamp of deployment

**Production deployment example:**
```bash
# Check production deployment metadata
kubectl get deployment sortinghat -n production \
  -o jsonpath='{.metadata.annotations.deployment/deployment-time}'
# Output: 2024-01-15T14:23:45Z

# Verify model versions match expectations  
kubectl get deployment sortinghat -n production \
  -o jsonpath='{.metadata.annotations.model-name/bge-base-v1-5}'
# Output: BAAI/bge-base-en-v1.5
```

## Model Licensing

All included models use commercial-friendly licenses:
- **MIT License**: Microsoft E5 models
- **Apache 2.0**: All other recommended models

✅ **Safe for commercial use without licensing fees** 
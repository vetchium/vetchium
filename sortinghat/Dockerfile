# Optimized Dockerfile for amd64 (Linux x86) and arm64 (Apple Silicon)
# Uses pre-built wheels for maximum speed on these common architectures
FROM python:3.10-slim

# Set environment variables for optimal pip behavior and model caching
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    HF_HOME=/root/.cache/huggingface

WORKDIR /app

# Install minimal runtime dependencies in one layer
RUN apt-get update && apt-get install -y --no-install-recommends \
    libgl1-mesa-glx \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

# Update pip to latest version for better wheel support
RUN pip install --upgrade pip

# Copy requirements first for optimal layer caching
COPY sortinghat/requirements.txt .

# Install packages with aggressive wheel preference for amd64/arm64
# These architectures have excellent pre-built wheel support
RUN pip install \
    --prefer-binary \
    --only-binary=torch,torchvision,torchaudio,scikit-learn,numpy,scipy \
    -r requirements.txt

# Pre-download sentence-transformers model in separate cached layer
# This ensures model is available at runtime without download delays
RUN python -c "from sentence_transformers import SentenceTransformer; SentenceTransformer('all-MiniLM-L6-v2')"

# Copy TypeSpec-generated Python models
COPY typespec/sortinghat/ ./sortinghat/

# Copy application code last to maximize cache hits during development
COPY sortinghat/main.py .

# Set runtime configuration
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["python", "main.py"] 
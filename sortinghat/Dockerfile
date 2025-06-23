FROM python:3.10-slim

WORKDIR /app

# Update pip
RUN pip install --no-cache-dir --upgrade pip

# Install build dependencies needed for PyMuPDF and other packages with native extensions
RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Install dependencies
COPY sortinghat/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Remove build dependencies to reduce image size
RUN apt-get update && apt-get remove -y \
    build-essential \
    gcc \
    g++ \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Pre-download sentence-transformers model
# This step caches the model to avoid downloading at runtime
RUN python -c "from sentence_transformers import SentenceTransformer; SentenceTransformer('all-MiniLM-L6-v2')"

# Copy TypeSpec-generated Python models
COPY typespec/sortinghat/ ./sortinghat/

# Copy application code
COPY sortinghat/main.py .

# Set environment variables
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["python", "main.py"] 
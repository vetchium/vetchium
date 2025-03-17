FROM python:3.10-slim

WORKDIR /app

# Install minimal dependencies for model downloads
RUN pip install --no-cache-dir sentence-transformers spacy

# Download models first - these rarely change and can be cached
RUN python -c "from sentence_transformers import SentenceTransformer; SentenceTransformer('all-MiniLM-L6-v2')"
RUN python -m spacy download en_core_web_md

# Install all other dependencies - this layer will change more frequently
COPY sorting-hat/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY sorting-hat/main.py .

# Set environment variables
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["python", "main.py"] 
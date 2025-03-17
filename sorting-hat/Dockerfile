FROM python:3.10-slim

WORKDIR /app

# Update pip
RUN pip install --no-cache-dir --upgrade pip

# Install dependencies
COPY sorting-hat/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Pre-download sentence-transformers model
# This step caches the model to avoid downloading at runtime
RUN python -c "from sentence_transformers import SentenceTransformer; SentenceTransformer('all-MiniLM-L6-v2')"

# Copy application code
COPY sorting-hat/main.py .

# Set environment variables
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["python", "main.py"] 
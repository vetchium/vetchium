FROM sqitch/sqitch:v1.4

# Install psql for the vetchium-tags-creat job
USER root
RUN apt-get update && \
    apt-get install -y postgresql-client && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /vetchi/migrations

# Copy the sqitch migrations to the container
COPY . /vetchi/migrations

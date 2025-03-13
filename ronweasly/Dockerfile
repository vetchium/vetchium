FROM node:22.11.0-alpine AS typespec-builder
WORKDIR /app/typespec
# Copy only package files first for better caching
COPY typespec/package*.json ./
# Install dependencies without running prepare script
RUN npm install --ignore-scripts
# Then copy source files
COPY typespec/*.json typespec/*.ts typespec/*.tsp typespec/*.yaml ./
COPY typespec/common ./common
COPY typespec/employer ./employer
COPY typespec/hub ./hub
# Now run the build explicitly
RUN npm run build && npm pack

# Development build stage
FROM node:22.11.0-alpine AS dev-builder
WORKDIR /app

ARG API_ENDPOINT
ENV API_ENDPOINT=$API_ENDPOINT
ENV NODE_ENV=development
ENV NEXT_TELEMETRY_DISABLED=1

# Copy package files first
COPY ronweasly/package*.json ./
# Copy the typespec archive
COPY --from=typespec-builder /app/typespec/psankar-vetchi-typespec-0.1.0.tgz ./
# Install dependencies with all devDependencies for better development experience
RUN npm install ./psankar-vetchi-typespec-0.1.0.tgz && \
    npm install

# Copy source files
COPY ronweasly/ .

# Development runner stage
FROM dev-builder AS dev-runner
WORKDIR /app
ENV NODE_ENV=development
ENV NEXT_TELEMETRY_DISABLED=1

EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Enable source maps and debugging
ENV NEXT_SHARP_PATH=/app/node_modules/sharp
ENV NODE_OPTIONS='--inspect=0.0.0.0:9229'

CMD ["npm", "run", "dev"]

# Production build stage
FROM node:22.11.0-alpine AS prod-builder
WORKDIR /app

ARG API_ENDPOINT
ENV API_ENDPOINT=$API_ENDPOINT
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Copy package files first
COPY ronweasly/package*.json ./
# Copy the typespec archive
COPY --from=typespec-builder /app/typespec/psankar-vetchi-typespec-0.1.0.tgz ./
# Install dependencies without devDependencies for production
RUN npm install ./psankar-vetchi-typespec-0.1.0.tgz && \
    npm install --omit=dev

# Copy source files
COPY ronweasly/ .
# Build for production with optimizations
RUN npm run build

# Production runner stage
FROM node:22.11.0-alpine AS prod-runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

COPY --from=prod-builder /app/public ./public
COPY --from=prod-builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=prod-builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# The runner stage inherits environment variables from the builder stage
ENV API_ENDPOINT=$API_ENDPOINT

CMD ["node", "server.js"]

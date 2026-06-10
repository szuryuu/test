# Stage 1: Build
FROM node:22-alpine AS builder

WORKDIR /app

# Cache deps first
COPY package.json package-lock.json* ./
RUN npm ci

# Copy source and build
COPY . .
RUN npm run build

# Stage 2: Serve with nginx
FROM nginx:1.27-alpine

# Custom nginx config
COPY --from=builder /app/nginx.conf /etc/nginx/conf.d/default.conf

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:80/ || exit 1

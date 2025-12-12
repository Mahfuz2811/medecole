FROM --platform=linux/arm64 node:22-slim

RUN apt-get update && apt-get install -y \
    python3 \
    make \
    g++ \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
# Remove any node_modules from host that might have wrong binaries
RUN rm -rf node_modules/.cache && npm rebuild

EXPOSE 3000
CMD ["npm", "run", "dev"]
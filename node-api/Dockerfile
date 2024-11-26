FROM docker.io/node:21.2.0-alpine

COPY . .
RUN npm ci

CMD ["node", "index.mjs"]

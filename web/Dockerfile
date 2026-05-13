# Use a Node.js Alpine image for the builder stage
FROM node:26-alpine AS builder
WORKDIR /app
COPY ./balembala/package*.json ./
RUN npm ci
COPY ./balembala/ .
RUN npm run build
RUN npm prune --production

# Use another Node.js Alpine image for the final stage
FROM node:26-alpine
WORKDIR /app
COPY --from=builder /app/build build/
COPY --from=builder /app/node_modules node_modules/
COPY ./balembala/package.json .
EXPOSE 3000
ENV NODE_ENV=production
ENV ORIGIN=http://localhost:3000
CMD [ "node", "build" ]

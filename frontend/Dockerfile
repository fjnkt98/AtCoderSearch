# docker build -t atcodersearch-frontend:0.1.0 .
# docker container run --rm --name atcodersearch-frontend -p 3000:3000 atcodersearch-frontend:0.1.0

FROM node:18-slim AS stage
WORKDIR /app

COPY . .
ENV API_HOST='http://localhost:8000'
RUN npm install && npm run build

FROM node:18-slim

COPY --from=stage /app/package.json .
COPY --from=stage /app/build /app/build

RUN npm install --omit=dev

CMD ["node", "/app/build"]

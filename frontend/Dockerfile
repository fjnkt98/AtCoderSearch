# docker build -t atcodersearch-frontend:0.1.0 .
# docker container run --rm --name atcodersearch-frontend -p 3000:3000 atcodersearch-frontend:0.1.0

FROM node:18-bullseye-slim AS stage
WORKDIR /app

COPY . .
ARG HOST='http://localhost:8000'
ENV API_HOST=${HOST}
RUN npm install && npm run build

FROM node:18-bullseye-slim

COPY --from=stage /app/package.json .
COPY --from=stage /app/build /app/build

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV TZ="Asia/Tokyo"

RUN npm install --omit=dev

USER node

CMD ["node", "/app/build"]

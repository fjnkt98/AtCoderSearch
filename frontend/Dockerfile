# docker build -t atcodersearch-frontend:latest .
# docker container run --rm --name atcodersearch-frontend -p 3000:3000 atcodersearch-frontend:latest

FROM node:18-bullseye-slim AS stage
WORKDIR /app

COPY . .
RUN npm install && npm run build

FROM node:18-bullseye-slim

COPY --from=stage /app/package.json .
COPY --from=stage /app/build /app/build

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV TZ="Asia/Tokyo"

RUN npm install --omit=dev

USER node

CMD ["node", "/app/build"]

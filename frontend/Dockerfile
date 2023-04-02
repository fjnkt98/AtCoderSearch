# docker build -t atcoder-search-frontend --build-arg API_HOST=http://localhost:8000 .
# docker run -p 80:80 --rm -d --name atcoder-search-frontend atcoder-search-frontend

FROM node:18 AS builder

WORKDIR /app

ARG API_HOST
ENV VITE_API_HOST ${API_HOST}

COPY ./src ./src
COPY ./package.json ./package.json
COPY ./package-lock.json ./package-lock.json
COPY ./index.html ./index.html
COPY ./postcss.config.cjs ./postcss.config.cjs
COPY ./tailwind.config.cjs ./tailwind.config.cjs
COPY ./tsconfig.json ./tsconfig.json
COPY ./tsconfig.node.json ./tsconfig.node.json
COPY ./vite.config.ts ./vite.config.ts

RUN npm install

RUN npm run build

FROM nginx:1.23

COPY --from=builder /app/dist/ /usr/share/nginx/html/
COPY ./default.conf.template /etc/nginx/templates/default.conf.template

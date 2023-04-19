# docker build -t atcoder-search-frontend -f middlewares/nginx/Dockerfile --build-arg API_HOST=http://localhost:8000 .
# docker run -p 80:80 --rm -d --name atcoder-search-frontend -e API_HOST=http://localhost:8000 atcoder-search-frontend

FROM node:18 AS builder

WORKDIR /app

ARG API_HOST
ENV VITE_API_HOST ${API_HOST}
ENV API_HOST ${API_HOST}

COPY ./frontend/src ./src
COPY ./frontend/package.json ./package.json
COPY ./frontend/package-lock.json ./package-lock.json
COPY ./frontend/index.html ./index.html
COPY ./frontend/postcss.config.cjs ./postcss.config.cjs
COPY ./frontend/tailwind.config.cjs ./tailwind.config.cjs
COPY ./frontend/tsconfig.json ./tsconfig.json
COPY ./frontend/tsconfig.node.json ./tsconfig.node.json
COPY ./frontend/vite.config.ts ./vite.config.ts
COPY ./frontend/AtCoderSearchLogo.svg ./AtCoderSearchLogo.svg

RUN npm install

RUN npm run build

FROM nginx:1.23

COPY --from=builder /app/dist/ /usr/share/nginx/html/
COPY ./middlewares/nginx/default.conf.template /etc/nginx/templates/default.conf.template

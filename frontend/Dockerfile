FROM node:18-bullseye-slim AS stage

WORKDIR /app

COPY ./src /app/src
COPY ./static /app/static
COPY ./package-lock.json /app/package-lock.json
COPY ./package.json /app/package.json
COPY ./postcss.config.js /app/postcss.config.js
COPY ./svelte.config.js /app/svelte.config.js
COPY ./tailwind.config.js /app/tailwind.config.js
COPY ./tsconfig.json /app/tsconfig.json
COPY ./vite.config.ts /app/vite.config.ts

RUN npm install && npm run build

FROM node:18-bullseye-slim

COPY --from=stage /app/package.json .
COPY --from=stage /app/build /app/build

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV TZ="Asia/Tokyo"

RUN npm install --omit=dev

USER node

CMD ["node", "/app/build"]

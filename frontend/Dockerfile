FROM node:22.11-bullseye-slim AS stage

WORKDIR /app

COPY ./.eslintrc.cjs /app/.eslintrc.cjs
COPY ./app /app/app 
COPY ./package-lock.json /app/package-lock.json
COPY ./package.json /app/package.json
COPY ./postcss.config.js /app/postcss.config.js
COPY ./public /app/public
COPY ./tailwind.config.ts /app/tailwind.config.ts
COPY ./tsconfig.json /app/tsconfig.json
COPY ./vite.config.ts /app/vite.config.ts

RUN npm install && npm run build

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV TZ="Asia/Tokyo"

USER node

CMD ["npm", "start"]

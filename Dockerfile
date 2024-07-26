FROM golang:alpine AS builder

WORKDIR /app

COPY server/go.mod server/go.sum ./

RUN go mod download

COPY server/main.go server/config.go ./
COPY server/pkg ./pkg

RUN go build -o server .

FROM node:20-alpine AS ui-base

RUN npm i -g npm@latest

FROM ui-base AS ui-deps

RUN mkdir /app
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

FROM ui-base as ui-builder

RUN mkdir /app
WORKDIR /app

COPY --from=ui-deps /app/node_modules ./node_modules
COPY src ./src
COPY public ./public
COPY package.json package-lock.json ./
COPY tsconfig.app.json tsconfig.json tsconfig.node.json ./
COPY vite.config.ts tailwind.config.ts postcss.config.js index.html .env.production ./

RUN npm run build

FROM alpine

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=ui-builder /app/dist ./dist

EXPOSE 8080

CMD ["./server"]
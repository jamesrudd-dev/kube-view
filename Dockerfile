### Build React Frontend ###
FROM node:14 AS frontend-builder

ENV WEB_SERVER_PATH="/kube-view" 
ENV PUBLIC_URL=${WEB_SERVER_PATH}

WORKDIR /kube-view-frontend
COPY ./frontend .

RUN npm install && npm run build

### Build Go Backend ###
FROM golang:1.17.11-alpine as go-builder

WORKDIR /kube-view
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/kube-view cmd/kube-view/*.go

### Build Redis Final Image ###
FROM redis:7-alpine as final

ENV IN_PRODUCTION="true"
ENV KUBE_CONFIG_LOCATION="/etc/config/kube-config"
ENV WEB_SERVER_PATH="/kube-view"
ENV ImageTagFilter=""
ENV NamespaceFilter=""

WORKDIR /app
COPY --from=go-builder ./kube-view/build/kube-view .
COPY --from=frontend-builder ./kube-view-frontend/build ./frontend/build
COPY docker-entrypoint.sh /bin

EXPOSE 8080

ENTRYPOINT [ "/bin/docker-entrypoint.sh" ]
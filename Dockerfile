FROM node:14 AS frontend-builder
WORKDIR /kube-view-frontend
COPY ./frontend .
RUN npm install && npm run build

FROM golang:1.17.11-alpine as go-builder
WORKDIR /kube-view
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/kube-view cmd/kube-view/*.go

FROM redis:7-alpine as final
ENV IN_PRODUCTION="false"
ENV KUBE_CONFIG_LOCATION="/etc/config/kube-config"
WORKDIR /app
COPY --from=go-builder ./kube-view/build/kube-view .
COPY --from=frontend-builder ./kube-view-frontend/build ./frontend/build
COPY startup.sh /bin
EXPOSE 8080
ENTRYPOINT [ "/bin/startup.sh" ]
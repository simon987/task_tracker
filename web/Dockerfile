FROM node:10-alpine as npm_build
COPY ./ ./
RUN cd ./angular/ && npm install
RUN cd ./angular/ && ./node_modules/@angular/cli/bin/ng build --prod --optimization --output-path "/webroot"

FROM nginx:alpine
WORKDIR /root

COPY nginx.conf /etc/nginx/
COPY --from=npm_build ["/webroot", "/webroot"]

EXPOSE 80

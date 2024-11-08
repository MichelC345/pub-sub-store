FROM node:14.16.1-alpine3.10 AS base
WORKDIR /var/www/

FROM base AS contact-service
ADD  services/contact/ .
RUN npm install --only=production 
CMD [ "node", "app.js" ]

FROM base AS order-service
ADD  services/order/ .
RUN npm install --only=production 
CMD [ "node", "app.js" ]

FROM base AS shipping-service
ADD  services/shipping/ .
RUN npm install --only=production 
CMD [ "node", "app.js" ]

# Build stage for Go report-service
FROM golang:1.23-alpine AS report-build
WORKDIR /app
ADD services/report/ .
RUN go build -o report-service .

# Final stage for Go report-service
FROM alpine:3.10 AS report-service
WORKDIR /var/www/
COPY --from=report-build /app/report-service .
CMD ["./report-service"]
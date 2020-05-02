# Fetch the dependencies
FROM golang:1.13-alpine AS builder

RUN apk add --update ca-certificates git gcc g++ libc-dev
WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor

# Build the final image
FROM hashicorp/terraform:0.12.1
COPY --from=builder /src/terraform-provider-statuscake /root/.terraform.d/plugins/

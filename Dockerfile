FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags lambda.norpc -o /main ./cmd/leaderboard-lambda

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /main /main
ENTRYPOINT ["/main"]

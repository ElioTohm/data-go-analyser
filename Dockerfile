FROM golang:1.16 as builder
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./app

FROM public.ecr.aws/lambda/go:latest
# Copy function code
COPY --from=builder /workspace/app ${LAMBDA_TASK_ROOT}
# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "app" ]

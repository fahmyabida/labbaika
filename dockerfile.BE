FROM golang:alpine as builder

WORKDIR /app

COPY ./backend .

RUN go build -o labbaika-payslip .

RUN chmod +x labbaika-payslip

FROM scratch

# Copy our static executable.
COPY --from=builder /app/labbaika-payslip /go/bin/labbaika-payslip
# Run the hello binary.
ENTRYPOINT ["/go/bin/labbaika-payslip"]


# CMD ["/bin/bash"]
# CMD ["labbaika-payslip rest"]

# FROM golang:1.19 AS build-stage

# WORKDIR /app

# COPY ./backend .

# RUN go build -o labbaika-payslip

# FROM gcr.io/distroless/base-debian11 AS build-release-stage

# WORKDIR /

# COPY --from=build-stage /docker-gs-ping /docker-gs-ping

# EXPOSE 8080

# USER nonroot:nonroot

# ENTRYPOINT ["/docker-gs-ping"]
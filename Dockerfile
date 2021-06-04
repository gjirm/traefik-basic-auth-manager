ARG  BUILDER_IMAGE=golang:alpine
############################
# STEP 1 build executable binary
############################
FROM ${BUILDER_IMAGE} as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=11111

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

#WORKDIR $GOPATH/src/tbam/

# COPY go.mod .
# COPY go.sum .
# COPY app/ app/

#COPY internal/ internal/
COPY templates/* /templates/

RUN mkdir /tbam

COPY tbam /tbam/tbam-server

# ENV GO111MODULE=on
# RUN go mod download
# RUN go mod verify

# Build the binary
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
#     -ldflags='-w -s' -a \
#     -o /tbam/tbam-server app/main.go

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy static executable
COPY --from=builder --chown=appuser:appuser /tbam/tbam-server /tbam/tbam-server

# Copy templates
COPY --from=builder --chown=appuser:appuser /templates/* /tbam/templates/

# Use an unprivileged user.
USER appuser:appuser

WORKDIR /tbam

EXPOSE 8080

# Run the tbam binary.
ENTRYPOINT ["/tbam/tbam-server"]

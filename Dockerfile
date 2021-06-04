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
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY templates/* /templates/

RUN mkdir /tbam

COPY tbam /tbam/tbam-server

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

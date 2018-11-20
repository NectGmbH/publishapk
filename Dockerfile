FROM alpine:3.8

# Path to the apk which should be published
ENV APK ""

# Package name of the application
ENV PKG ""

# Service accounts email address
ENV EMAIL ""

# Service accounts private key, base64 encoded OR file path to private key
ENV KEY ""

# Track to which the apk should be pushed (e.g.: alpha, beta, internal, production)
ENV TRACK ""

COPY ./publishapk /bin/publishapk

CMD /bin/publishapk
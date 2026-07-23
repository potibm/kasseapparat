FROM node:22-alpine

RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

RUN npm install --global md-to-pdf

RUN addgroup -S pptruser && adduser -S pptruser -G pptruser \
    && mkdir -p /app \
    && chown -R pptruser:pptruser /app

USER pptruser
WORKDIR /app

ENTRYPOINT ["md-to-pdf"]
CMD ["--help"]
FROM ghcr.io/puppeteer/puppeteer:latest

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=1

USER root

RUN npm install --global md-to-pdf

USER pptruser

WORKDIR /app

ENTRYPOINT ["md-to-pdf"]

CMD ["--help"]
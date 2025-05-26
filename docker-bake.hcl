group "default" {
  targets = ["kasseapparat"]
}

target "kasseapparat" {
  context    = "."
  dockerfile = "Dockerfile"
  platforms  = ["linux/amd64", "linux/arm64"]

  args = {
    VERSION    = "dev-unknown"
    BUILD_DATE = "1970-01-01T00:00:00Z"
  }

  output = ["type=image,push=true"]

  labels = {
    "org.opencontainers.image.title"         = "kasseapparat"
    "org.opencontainers.image.description"   = "A POS system for demoparties"
    "org.opencontainers.image.url"           = "https://github.com/potibm/kasseapparat"
    "org.opencontainers.image.source"        = "https://github.com/potibm/kasseapparat"
    "org.opencontainers.image.documentation" = "https://github.com/potibm/kasseapparat/tree/main/doc"
    "org.opencontainers.image.licenses"      = "MIT"
    "org.opencontainers.image.authors"       = "potibm"
    "org.opencontainers.image.version"       = "dev-unknown"
    "org.opencontainers.image.created"       = "1970-01-01T00:00:00Z"
  }

  annotations = [
    "index,manifest:org.opencontainers.image.title=kasseapparat",
    "index,manifest:org.opencontainers.image.description=A POS system for demoparties",
    "index,manifest:org.opencontainers.image.url=https://github.com/potibm/kasseapparat",
    "index,manifest:org.opencontainers.image.source=https://github.com/potibm/kasseapparat",
    "index,manifest:org.opencontainers.image.documentation=https://github.com/potibm/kasseapparat/tree/main/doc",
    "index,manifest:org.opencontainers.image.licenses=MIT",
    "index,manifest:org.opencontainers.image.authors=potibm",
    "index,manifest:org.opencontainers.image.version=dev-unknown",
    "index,manifest:org.opencontainers.image.created=1970-01-01T00:00:00Z"
  ]
}

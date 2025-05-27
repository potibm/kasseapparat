group "default" {
  targets = ["kasseapparat"]
}

target "kasseapparat" {
  context    = "."
  dockerfile = "Dockerfile"
  platforms  = ["linux/amd64", "linux/arm64"]

  labels = {
    "org.opencontainers.image.url" = "https://github.com/potibm/kasseapparat"
    "org.opencontainers.image.source" = "https://github.com/potibm/kasseapparat"
    "org.opencontainers.image.documentation" = "https://github.com/potibm/kasseapparat/tree/main/doc"
    "org.opencontainers.image.authors" = "potibm"
  }
  
  annotations = [
    "index,manifest:org.opencontainers.image.title=Kasseapparat",
    "index,manifest:org.opencontainers.image.description=A POS system for demoparties",
    "index,manifest:org.opencontainers.image.url=https://github.com/potibm/kasseapparat",
    "index,manifest:org.opencontainers.image.source=https://github.com/potibm/kasseapparat",
    "index,manifest:org.opencontainers.image.documentation=https://github.com/potibm/kasseapparat/tree/main/doc",
    "index,manifest:org.opencontainers.image.licenses=MIT",
    "index,manifest:org.opencontainers.image.authors=potibm"
  ]
}

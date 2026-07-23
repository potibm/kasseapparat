# Supply Chain Security

This container image is built and secured automatically via [GitHub Actions](https://github.com/features/actions).

- Signed with [Sigstore Cosign](https://docs.sigstore.dev/)
- SBOM (Software Bill of Materials) included in [SPDX format](https://spdx.dev/)
- Attestations for metadata and SBOM
- Verifiable using `cosign verify-attestation`

## Verify the Attestation

You can verify the SBOM attestation with Cosign:

```bash
cosign verify-attestation \
  --type spdx \
  --certificate-identity-regexp "github.com/potibm/.*" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ghcr.io/potibm/kasseapparat@sha256:<digest>
```

Get the digest for a given version tag:

```bash
docker buildx imagetools inspect ghcr.io/potibm/kasseapparat:<version> \
  | grep -m1 Digest | awk '{print $2}'
```

## Example Image Tag

```bash
ghcr.io/potibm/kasseapparat:2.3.7
```

For full details on how the image is built and signed, see the [release workflow](.github/workflows/release.yml).

import React, { useEffect, useState } from "react";
import { useConfig } from "@core/config/hooks/useConfig";
import { Link } from "react-router-dom";
import { z } from "zod";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Api");

const VERSION_URL =
  "https://api.github.com/repos/potibm/kasseapparat/releases/latest";
const LS_KEY_LATEST_VERSION = "kasseapparat.latest_version";
const VersionStringSchema = z
  .string()
  .regex(/^v?\d+\.\d+\.\d+$/)
  .transform((v) => v.replace(/^v/, ""));

const GitHubReleaseSchema = z.object({
  tag_name: VersionStringSchema,
});

interface VersionProps {
  url?: string;
}

export const Version: React.FC<VersionProps> = ({ url = VERSION_URL }) => {
  const { version } = useConfig();

  const [latestVersion, setLatestVersion] = useState(() => {
    const stored = sessionStorage.getItem(LS_KEY_LATEST_VERSION);
    const result = VersionStringSchema.safeParse(stored);
    return result.success ? result.data : null;
  });

  useEffect(() => {
    if (latestVersion) {
      return;
    }

    const controller = new AbortController();

    const fetchVersion = async () => {
      try {
        const res = await fetch(url, { signal: controller.signal });
        const data = await res.json();

        // Validierung
        const result = GitHubReleaseSchema.safeParse(data);

        if (result.success && !controller.signal.aborted) {
          sessionStorage.setItem(LS_KEY_LATEST_VERSION, result.data.tag_name);
          setLatestVersion(result.data.tag_name);
        }
      } catch (err: unknown) {
        if (err instanceof Error && err.name !== "AbortError") {
          log.error("Error loading version from GitHub:", err);
        }
      }
    };

    fetchVersion();

    return () => controller.abort();
  }, [url, latestVersion]); // latestVersion im Guard oben reicht aus

  const isDevelopment =
    !version || version.startsWith("0") || version === "dev";
  const versionLink = isDevelopment
    ? null
    : `https://github.com/potibm/kasseapparat/releases/tag/v${version.replace(/^v/, "")}`;

  const isOutdated =
    latestVersion && latestVersion !== version && !isDevelopment;

  if (!versionLink) {
    return <span className="text-gray-500">Version {version || "dev"}</span>;
  }

  return versionLink ? (
    <Link
      to={versionLink}
      target="_blank"
      rel="noopener noreferrer"
      reloadDocument
      className={isOutdated ? "text-red-600 font-semibold" : ""}
      title={
        isOutdated
          ? `A newer version (${latestVersion}) is available.`
          : "You are using the current version"
      }
    >
      Version {version}
    </Link>
  ) : (
    <>Version {version}</>
  );
};

export default Version;

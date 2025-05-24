import React, { useEffect, useState } from "react";
import { useConfig } from "../provider/ConfigProvider";
import { Link } from "react-router-dom";

export const Version = () => {
  const version = useConfig().version;

  const [latestVersion, setLatestVersion] = useState(null);

  useEffect(() => {
    const cached = sessionStorage.getItem("kasseapparat_latest_version");
    if (cached) {
      setLatestVersion(cached);
      return;
    }

    const controller = new AbortController();

    fetch("https://api.github.com/repos/potibm/kasseapparat/releases/latest", {
      signal: controller.signal,
    })
      .then((res) => res.json())
      .then((data) => {
        if (data?.tag_name) {
          const cleanVersion = data.tag_name.replace(/^v/, "");
          sessionStorage.setItem("kasseapparat_latest_version", cleanVersion);
          setLatestVersion(cleanVersion);
        }
      })
      .catch((err) => {
        if (err.name !== "AbortError") {
          console.error("Error loading the version from GitHub:", err);
        }
      });

    return () => controller.abort();
  }, []);

  let versionLink = null;
  if (version && !version.startsWith("0")) {
    versionLink = `https://github.com/potibm/kasseapparat/releases/tag/v${version}`;
  }

  const isOutdated = latestVersion && latestVersion !== version;

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

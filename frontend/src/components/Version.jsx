import React from "react";
import { useConfig } from "../provider/ConfigProvider";
import { Link } from "react-router-dom";
import { useEffect, useState } from "react";

export const Version = () => {
  const version = useConfig().version;

  const [latestVersion, setLatestVersion] = useState(null);

  useEffect(() => {
    fetch("https://api.github.com/repos/potibm/kasseapparat/releases/latest")
      .then((res) => res.json())
      .then((data) => {
        if (data?.tag_name) {
          setLatestVersion(data.tag_name.replace(/^v/, ""));
        }
      })
      .catch((err) => {
        console.error("Error loading the version from GitHub:", err);
      });
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

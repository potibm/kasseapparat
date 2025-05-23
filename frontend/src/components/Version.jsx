import React from "react";
import { useConfig } from "../provider/ConfigProvider";
import { Link } from "react-router-dom";

export const Version = () => {
  const version = useConfig().version;

  let versionLink = null;
  if (version && !version.startsWith("0")) {
    versionLink = `https://github.com/potibm/kasseapparat/releases/tag/v${version}`;
  }

  return versionLink ? (
    <Link to={versionLink} target="_blank" reloadDocument>
      Version {version}
    </Link>
  ) : (
    <>Version {version}</>
  );
};

export default Version;

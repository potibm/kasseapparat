import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "./index.css";
import App from "./App.tsx";

import { ConfigContext } from "@core/config/context/ConfigContext.ts";
import { createLogger } from "@core/logger/logger.ts";
import { ConfigSchema } from "@core/config/schemas/config.schemas.ts";
import * as Sentry from "@sentry/react";
import { transformConfig } from "@core/config/utils/config.transform.ts";
import { buildApiBaseUrl } from "@core/config/constants.ts";
import { CriticalError } from "@core/components/CriticalError.tsx";

const log = createLogger("Bootstrapper");
const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3100";

export async function bootstrapApp() {
  const rootElement = document.getElementById("root");
  if (!rootElement) {
    const root = createRoot(document.body);
    root.render(
      <CriticalError
        title="System Configuration Error"
        message="Failed to find the root element in index.html"
      />,
    );
    log.error("Bootstrap failed: Root element missing");
    return;
  }

  const root = createRoot(rootElement);

  try {
    const controller = new AbortController();
    const timeoutId = globalThis.setTimeout(() => controller.abort(), 10000);

    let data: unknown;

    // 1. Fetch config and validate
    try {
      const res = await fetch(`${buildApiBaseUrl(API_HOST)}/config`, {
        signal: controller.signal,
      });
      if (!res.ok) throw new Error(`Config error: ${res.statusText}`);
      data = await res.json();
    } finally {
      globalThis.clearTimeout(timeoutId);
    }

    const parsedData = ConfigSchema.parse(data);

    const config = transformConfig(parsedData, API_HOST);
    log.debug("Config loaded:", config);

    // 2. Initialize Sentry
    if (config.sentryDSN && !Sentry.isInitialized()) {
      log.debug("Configuring Sentry");
      Sentry.init({
        dsn: config.sentryDSN,
        environment: config.sentryEnvironment,
        release: config.version,
        replaysSessionSampleRate: config.sentryReplaySessionSampleRate,
        replaysOnErrorSampleRate: config.sentryReplayErrorSampleRate,
        integrations: [
          Sentry.replayIntegration(),
          Sentry.browserTracingIntegration(),
        ],
      });
    }

    // 3. Start React
    root.render(
      <StrictMode>
        <ConfigContext value={config}>
          <App />
        </ConfigContext>
      </StrictMode>,
    );
  } catch (err) {
    log.error("Bootstrap failed:", err);

    root.render(
      <CriticalError
        title="System Configuration Error"
        message="Failed to initialize application"
        details={err instanceof Error ? err.message : "Unknown error"}
      />,
    );
  }
}

if (!import.meta.env.TEST) {
  await bootstrapApp();
}

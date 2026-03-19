import React, { useEffect } from "react";
import * as Sentry from "@sentry/react";
import { useConfig } from "../config/providers/ConfigProvider";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Core");

const SentryInitializer: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const {
    sentryDSN,
    version,
    sentryReplayErrorSampleRate,
    sentryReplaySessionSampleRate,
    sentryTraceSampleRate,
  } = useConfig();

  useEffect(() => {
    if (sentryDSN) {
      try {
        Sentry.init({
          dsn: sentryDSN,
          environment: import.meta.env.MODE,
          release: `kasseapparat@${version ?? "unknown"}`,
          integrations: [
            Sentry.browserTracingIntegration(),
            Sentry.replayIntegration(),
          ],
          tracesSampleRate: sentryTraceSampleRate ?? 1,
          replaysSessionSampleRate: sentryReplaySessionSampleRate ?? 0.1,
          replaysOnErrorSampleRate: sentryReplayErrorSampleRate ?? 1,
        });
      } catch (error: unknown) {
        log.error("Error initializing Sentry", error);
      }
    }
  }, [
    sentryDSN,
    version,
    sentryReplayErrorSampleRate,
    sentryReplaySessionSampleRate,
    sentryTraceSampleRate,
  ]);

  return (
    <Sentry.ErrorBoundary fallback={<p>Critical Error in Kasseapparat.</p>}>
      {children}
    </Sentry.ErrorBoundary>
  );
};

export default SentryInitializer;

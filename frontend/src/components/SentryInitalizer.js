import { useEffect } from "react";
import * as Sentry from "@sentry/react";
import { useConfig } from "../provider/ConfigProvider";

const SentryInitializer = ({ children }) => {
  const config = useConfig();

  useEffect(() => {
    if (config && config.sentryDSN) {
      try {
        Sentry.init({
          dsn: config.sentryDSN,
          environment: process.env.NODE_ENV,
          integrations: [
            Sentry.browserTracingIntegration(),
            Sentry.replayIntegration(),
          ],
          tracesSampleRate: config.sentryTraceSampleRate ?? 1.0,
          replaysSessionSampleRate: config.sentryReplaySessionSampleRate ?? 0.1,
          replaysOnErrorSampleRate: config.sentryReplayErrorSampleRate ?? 1.0,
        });
      } catch (error) {
        console.error("Error initializing Sentry:", error);
      }
    }
  }, [config]);

  return children;
};

export default SentryInitializer;

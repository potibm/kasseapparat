import { useEffect } from "react";
import * as Sentry from "@sentry/react";
import { useConfig } from "../provider/ConfigProvider";
import { useAuth } from "../Auth/provider/AuthProvider";

const SentryInitializer = ({ children }) => {
  const config = useConfig();
  const { id, username } = useAuth();

  useEffect(() => {
    if (config?.sentryDSN) {
      try {
        Sentry.init({
          dsn: config.sentryDSN,
          environment: import.meta.env.MODE,
          release: `kasseapparat@${config?.version ?? "unknown"}`,
          integrations: [
            Sentry.browserTracingIntegration(),
            Sentry.replayIntegration(),
          ],
          tracesSampleRate: config.sentryTraceSampleRate ?? 1.0,
          replaysSessionSampleRate: config.sentryReplaySessionSampleRate ?? 0.1,
          replaysOnErrorSampleRate: config.sentryReplayErrorSampleRate ?? 1.0,
        });

        Sentry.setUser({
          id,
          username,
        });
      } catch (error) {
        console.error("Error initializing Sentry:", error);
      }
    }
  }, [config, id, username]);

  return children;
};

export default SentryInitializer;

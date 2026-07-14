import { useEffect } from "react";
import * as Sentry from "@sentry/react";
import { useAuth } from "../../apps/pos/features/auth/hooks/useAuth";

const SentryUserWatcher = () => {
  const { username } = useAuth();

  useEffect(() => {
    if (username) {
      Sentry.setUser({ username });
    } else {
      Sentry.setUser(null);
    }
  }, [username]);

  return null;
};

export default SentryUserWatcher;

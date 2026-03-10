import { useEffect } from "react";
import * as Sentry from "@sentry/react";
import { useAuth } from "../../apps/pos/features/auth/providers/auth-provider";

const SentryUserWatcher = () => {
  const { id, username } = useAuth();

  useEffect(() => {
    if (id) {
      Sentry.setUser({ id: String(id), username });
    } else {
      Sentry.setUser(null);
    }
  }, [id, username]);

  return null;
};

export default SentryUserWatcher;

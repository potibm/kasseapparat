//import SentryUserWatcher from "./core/monitoring/SentryUserWatcher";
import * as Sentry from "@sentry/react";
import Routes from "./routes";
import { AuthProvider } from "@pos/features/auth/provider/AuthProvider";
import { CriticalError } from "@core/components/CriticalError";

function App() {
  return (
    <Sentry.ErrorBoundary
      fallback={
        <CriticalError
          title="Application Error"
          message="A serious error has occurred. Please restart the Kasseapparat."
        />
      }
    >
      {/* <SentryUserWatcher /> */}
      <AuthProvider>
        <Routes />
      </AuthProvider>
    </Sentry.ErrorBoundary>
  );
}

export default App;

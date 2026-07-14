//import SentryUserWatcher from "./core/monitoring/SentryUserWatcher";
import * as Sentry from "@sentry/react";
import Routes from "./routes";
import { AuthProvider } from "@pos/features/auth/provider/AuthProvider";

function App() {
  return (
    <Sentry.ErrorBoundary
      fallback={
        <p>A serious error has occurred. Please restart the Kasseapparat.</p>
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

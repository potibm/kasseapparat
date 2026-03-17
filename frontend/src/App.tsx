import ConfigProvider from "./core/config/providers/ConfigProvider";
import AuthProvider from "./apps/pos/features/auth/providers/AuthProvider";
import SentryInitializer from "./core/monitoring/SentryInitializer";
import SentryUserWatcher from "./core/monitoring/SentryUserWatcher";
import Routes from "./routes";

function App() {
  return (
    <ConfigProvider>
      <SentryInitializer>
        <AuthProvider>
          <SentryUserWatcher />
          <Routes />
        </AuthProvider>
      </SentryInitializer>
    </ConfigProvider>
  );
}

export default App;

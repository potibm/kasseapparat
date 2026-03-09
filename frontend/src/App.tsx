import ConfigProvider from "./core/config/providers/ConfigProvider";
import AuthProvider from "./apps/pos/features/auth/providers/auth-provider";
import SentryProvider from "./core/monitoring/SentryProvider";
import Routes from "./routes";

function App() {
  return (
    <ConfigProvider>
      <SentryProvider>
        <AuthProvider>
          <Routes />
        </AuthProvider>
      </SentryProvider>
    </ConfigProvider>
  );
}

export default App;

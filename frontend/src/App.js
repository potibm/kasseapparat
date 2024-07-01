import React from "react";
import "./App.css";
import AuthProvider from "./Auth/provider/AuthProvider";
import Routes from "./routes";
import SentryInitializer from "./components/SentryInitalizer";
import ConfigProvider from "./provider/ConfigProvider";

function App() {
  return (
    <ConfigProvider>
      <SentryInitializer>
        <AuthProvider>
          <Routes />
        </AuthProvider>
      </SentryInitializer>
    </ConfigProvider>
  );
}

export default App;

import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { refreshJwtToken } from "../Main/hooks/Api";
import PropTypes from "prop-types";

const AuthContext = createContext();

const API_HOST = process.env.REACT_APP_API_HOST ?? "http://localhost:3001";

const AuthProvider = ({ children }) => {
  const [auth, setAuth] = useState({
    token: localStorage.getItem("token"),
    expiryDate: localStorage.getItem("expiryDate"),
    username: localStorage.getItem("username"),
  });

  const performTokenRefresh = useCallback(() => {
    if (auth.token == null) {
      return;
    }
    refreshJwtToken(API_HOST, auth.token)
      .then((response) => {
        const newToken = response.token;
        const newExpiryDate = response.expire;
        const username = auth.username;
        console.log("Refreshed token, Expiry: " + newExpiryDate);

        setAuth({ token: newToken, expiryDate: newExpiryDate, username });
      })
      .catch((error) => {
        console.error("Error refreshing the token", error);
      });

    // check if token is expired
    const now = new Date();
    const expiryDate = new Date(auth.expiryDate);
    if (now > expiryDate) {
      // @todo notify user
      setAuth({ token: null, expiryDate: null, username: null });
      window.location = "/logout";
    }
  }, [auth.token, auth.expiryDate, auth.username]);

  useEffect(() => {
    if (auth.token) {
      localStorage.setItem("token", auth.token);
    } else {
      localStorage.removeItem("token");
    }

    if (auth.expiryDate) {
      localStorage.setItem("expiryDate", auth.expiryDate);
    } else {
      localStorage.removeItem("expiryDate");
    }

    if (auth.username) {
      localStorage.setItem("username", auth.username);
    } else {
      localStorage.removeItem("username");
    }
  }, [auth]);

  // ensure token validity is checked when app is focused or revived
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.visibilityState === "visible") {
        handleReviviedApp();
      }
    };

    const handleReviviedApp = () => {
      if (auth.token == null) {
        return;
      }

      const now = new Date();
      const expiryDate = new Date(auth.expiryDate);
      if (now > expiryDate) {
        setAuth({ token: null, expiryDate: null, username: null });
        window.location = "/logout";
      }
      // Refresh the token if it will expire in the next two minutes
      if (expiryDate - now < 2 * 60 * 1000) {
        performTokenRefresh();
      }
    };

    window.addEventListener("focus", handleReviviedApp);
    document.addEventListener("visibilitychange", handleVisibilityChange);

    return () => {
      window.removeEventListener("focus", handleReviviedApp);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    };
  }, [auth.expiryDate, auth.token, performTokenRefresh]);

  useEffect(() => {
    const tokenRefreshInterval = setInterval(() => {
      performTokenRefresh();
    }, 60 * 1000); // Refresh the token every 60 seconds

    return () => {
      clearInterval(tokenRefreshInterval);
    };
  }, [auth.expiryDate, auth.token, auth.username, performTokenRefresh]);

  const contextValue = useMemo(
    () => ({
      token: auth.token,
      setToken: (token) => setAuth((prev) => ({ ...prev, token })),
      expiryDate: auth.expiryDate,
      setExpiryDate: (expiryDate) =>
        setAuth((prev) => ({ ...prev, expiryDate })),
      username: auth.username,
      setUsername: (username) => setAuth((prev) => ({ ...prev, username })),
    }),
    [auth],
  );

  // Provide the authentication context to the children components
  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
};

export const useAuth = () => {
  return useContext(AuthContext);
};

AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export default AuthProvider;

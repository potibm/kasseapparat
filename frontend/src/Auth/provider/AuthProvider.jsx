import React, {
  createContext,
  useContext,
  useMemo,
  useRef,
  useState,
} from "react";
import { refreshJwtToken } from "../hooks/Api";
import PropTypes from "prop-types";
import { useConfig } from "../../provider/ConfigProvider";

const AuthContext = createContext();

const AuthProvider = ({ children }) => {
  const apiHost = useConfig().apiHost;
  const LOCALSTORAGE_PREFIX = "kasseapparat.auth.";
  const refreshingPromise = useRef(null);

  const updateSession = (token, expiresIn) => {
    const expiryDate = new Date(
      // eslint-disable-next-line react-hooks/purity
      Date.now() + (expiresIn - 10) * 1000,
    ).toISOString();
    console.log("Updating session with new expiry date: " + expiryDate);

    setSession({ token: token, expiryDate: expiryDate });
    localStorage.setItem(LOCALSTORAGE_PREFIX + "token", token);
    localStorage.setItem(LOCALSTORAGE_PREFIX + "expiryDate", expiryDate);

    return token;
  };

  const removeSession = () => {
    setSession({ token: null, expiryDate: null });
    localStorage.removeItem(LOCALSTORAGE_PREFIX + "token");
    localStorage.removeItem(LOCALSTORAGE_PREFIX + "expiryDate");
  };

  const updateUser = (userdata) => {
    setUser(userdata);
    localStorage.setItem(
      LOCALSTORAGE_PREFIX + "userdata",
      JSON.stringify(userdata),
    );
  };

  const getToken = async () => {
    const currentDate = new Date();
    const expiryDate = new Date(session.expiryDate);

    if (session.token && expiryDate > currentDate) {
      return session.token;
    } else if (!session.token) {
      console.log("No token found in session");
      return null;
    }

    if (refreshingPromise.current) {
      return refreshingPromise.current;
    }

    console.log("Token expired or missing, starting refresh...");

    refreshingPromise.current = refreshJwtToken(apiHost)
      .then((response) => {
        const newToken = response.token || response.access_token;
        const expiresIn = response.expire_in || 60;

        updateSession(newToken, expiresIn);

        return newToken;
      })
      .catch((error) => {
        console.error("Critical error during token refresh:", error);
        removeSession();

        throw error;
      })
      .finally(() => {
        refreshingPromise.current = null;
      });

    return refreshingPromise.current;
  };

  const isLoggedIn = async () => {
    const token = await getToken();
    return !!token;
  };

  const getSessionFromLocalStorage = () => {
    console.log("Getting session from local storage");
    const token = localStorage.getItem(LOCALSTORAGE_PREFIX + "token");
    const expiryDate = localStorage.getItem(LOCALSTORAGE_PREFIX + "expiryDate");

    if (!token || !expiryDate) {
      return { token: null, expiryDate: null };
    }

    if (new Date(expiryDate) < new Date()) {
      localStorage.removeItem(LOCALSTORAGE_PREFIX + "token");
      localStorage.removeItem(LOCALSTORAGE_PREFIX + "expiryDate");
      return { token: null, expiryDate: null };
    }

    return { token: token, expiryDate: expiryDate };
  };

  const getUserFromLocalStorage = () => {
    const userdata = localStorage.getItem(LOCALSTORAGE_PREFIX + "userdata");

    return userdata ? JSON.parse(userdata) : null;
  };

  const [session, setSession] = useState(getSessionFromLocalStorage);
  const [user, setUser] = useState(getUserFromLocalStorage);

  const contextValue = useMemo(
    () => ({
      getToken: async () => await getToken(),
      isLoggedIn: async () => await isLoggedIn(),
      setSession: (token, expiresIn) => updateSession(token, expiresIn),
      removeSession: () => removeSession(),
      userdata: user,
      setUserdata: (userdata) => updateUser(userdata),
      gravatarUrl: user?.gravatarUrl ?? "",
      role: user?.role ?? "user",
      username: user?.username ?? "unknown",
      id: user?.id ?? 0,
    }),
    [session, user],
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

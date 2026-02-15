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
  const LOCALSTORAGE_TOKEN_KEY = LOCALSTORAGE_PREFIX + "token";
  const LOCALSTORAGE_EXPIRY_KEY = LOCALSTORAGE_PREFIX + "expiryDate";
  const LOCALSTORAGE_USERDATA_KEY = LOCALSTORAGE_PREFIX + "userdata";
  const refreshingPromise = useRef(null);

  const updateSession = (token, expiresIn) => {
    const expiryDate = new Date(
      // eslint-disable-next-line react-hooks/purity
      Date.now() + (expiresIn - 30) * 1000,
    ).toISOString();
    console.log("Updating session with new expiry date: " + expiryDate);

    setSession({ token: token, expiryDate: expiryDate });
    localStorage.setItem(LOCALSTORAGE_TOKEN_KEY, token);
    localStorage.setItem(LOCALSTORAGE_EXPIRY_KEY, expiryDate);

    return token;
  };

  const removeSession = () => {
    setSession({ token: null, expiryDate: null });
    localStorage.removeItem(LOCALSTORAGE_TOKEN_KEY);
    localStorage.removeItem(LOCALSTORAGE_EXPIRY_KEY);
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
        const newToken = response.access_token;
        const expiresIn = response.expires_in || 60;

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
    const token = localStorage.getItem(LOCALSTORAGE_TOKEN_KEY);
    const expiryDate = localStorage.getItem(LOCALSTORAGE_EXPIRY_KEY);

    if (!token || !expiryDate) {
      return { token: null, expiryDate: null };
    }

    if (new Date(expiryDate) < new Date()) {
      localStorage.removeItem(LOCALSTORAGE_TOKEN_KEY);
      localStorage.removeItem(LOCALSTORAGE_EXPIRY_KEY);
      return { token: null, expiryDate: null };
    }

    return { token: token, expiryDate: expiryDate };
  };

  const getUserFromLocalStorage = () => {
    const userdata = localStorage.getItem(LOCALSTORAGE_USERDATA_KEY);

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

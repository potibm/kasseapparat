import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { refreshJwtToken } from "../Main/hooks/Api";

const AuthContext = createContext();

const API_HOST = process.env.REACT_APP_API_HOST ?? 'http://localhost:3001'

const AuthProvider = ({ children }) => {
  const [auth, setAuth] = useState({
    token: localStorage.getItem("token"),
    expiryDate: localStorage.getItem("expiryDate"),
    username: localStorage.getItem("username"),
  });
  
  useEffect(() => {
    if (auth.token) {
      localStorage.setItem('token', auth.token);
    } else {
      localStorage.removeItem('token');
    }
  
    if (auth.expiryDate) {
      localStorage.setItem('expiryDate', auth.expiryDate);
    } else {
      localStorage.removeItem('expiryDate');
    }

    if (auth.username) {
      localStorage.setItem('username', auth.username);
    } else {
      localStorage.removeItem('username');
    }
  }, [auth]);
  
  useEffect(() => {
    
    const tokenRefreshInterval = setInterval(() => {
      refreshJwtToken(API_HOST, auth.token).then((response) => {
        const newToken = response.token;
        const newExpiryDate = response.expire;
        const username = auth.username;
        console.log("Refreshed token, Expiry: " + newExpiryDate)
  
        setAuth({token: newToken, expiryDate: newExpiryDate, username: username})
      }).catch((error) => {
        console.error("Fehler beim Aktualisieren des Tokens: ", error);
      });
  
      // check if token is expired
      const now = new Date();
      const expiryDate = new Date(auth.expiryDate);
      if (now > expiryDate) {
        // @todo notify user
        setAuth({token: null, expiryDate: null, username: null});
        window.location = '/logout';
      }
    }, 60 * 1000); // Refresh the token every 60 seconds
  
    return () => {
      clearInterval(tokenRefreshInterval);
    };

  }, [auth.expiryDate]);
  
  const contextValue = useMemo(
    () => ({
      token: auth.token,
      setToken: (token) => setAuth((prev) => ({ ...prev, token })),
      expiryDate: auth.expiryDate,
      setExpiryDate: (expiryDate) => setAuth((prev) => ({ ...prev, expiryDate })),
      username: auth.username,
      setUsername: (username) => setAuth((prev) => ({ ...prev, username })),
    }),
    [auth]
  );

  // Provide the authentication context to the children components
  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
};

export const useAuth = () => {
  return useContext(AuthContext);
};

export default AuthProvider;
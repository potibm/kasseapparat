import inMemoryJWT from './inMemoryJWT';

const API_HOST = process.env.REACT_APP_API_HOST ?? 'http://localhost:3001'

const authProvider = {
    login: ({ username, password }) => {
        const request = new Request(API_HOST + '/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
            headers: new Headers({ 'Content-Type': 'application/json' }),
            credentials: 'include',
        });
        inMemoryJWT.setRefreshTokenEndpoint(API_HOST + '/auth/refresh_token');
        return fetch(request)
            .then((response) => {
                if (response.status < 200 || response.status >= 300) {
                    throw new Error(response.statusText);
                }
                return response.json();
            })
            .then(({ token, expire }) => {
                const expireDate = new Date(expire);
                const delay = expireDate.getTime() - Date.now();
                return inMemoryJWT.setToken(token, delay);
            });
    },

    logout: () => {
        /*const request = new Request('http://localhost:8001/logout', {
            method: 'GET',
            headers: new Headers({ 'Content-Type': 'application/json' }),
            credentials: 'include',
        });*/
        inMemoryJWT.ereaseToken();

        //return fetch(request).then(() => '/login');
        return Promise.resolve();
    },

    checkAuth: () => {
        return inMemoryJWT.waitForTokenRefresh().then(() => {
            return inMemoryJWT.getToken() ? Promise.resolve() : Promise.reject();
        });
    },

    checkError: (error) => {
        const status = error.status;
        if (status === 401 || status === 403) {
            inMemoryJWT.ereaseToken();
            return Promise.reject();
        }
        return Promise.resolve();
    },

    getPermissions: () => {
        return inMemoryJWT.waitForTokenRefresh().then(() => {
            return inMemoryJWT.getToken() ? Promise.resolve() : Promise.reject();
        });
    },
};

export default authProvider; 
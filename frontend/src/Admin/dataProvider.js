import jsonServerProvider from 'ra-data-json-server'
import { fetchUtils } from 'react-admin'
import inMemoryJWT from './inMemoryJWT'

const API_HOST = process.env.REACT_APP_API_HOST ?? 'http://localhost:3001'

const httpClient = (url, options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: 'application/json' })
  }
  // add your own headers here
  const token = inMemoryJWT.getToken()
  console.log('Token', token)
  options.headers.set('Authorization', `Bearer ${token}`)
  return fetchUtils.fetchJson(url, options)
}

const dataProvider = jsonServerProvider(API_HOST + '/api/v1', httpClient)

export default dataProvider

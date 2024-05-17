import jsonServerProvider from 'ra-data-json-server';

const API_HOST = process.env.REACT_APP_API_HOST ?? 'http://localhost:3001'

const dataProvider = jsonServerProvider(API_HOST + "/api/v1");

export default dataProvider;
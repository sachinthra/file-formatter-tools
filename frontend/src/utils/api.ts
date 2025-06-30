import { API_KEY } from '../config';

export async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const headers = new Headers(options.headers);
  headers.set('X-API-Key', API_KEY);

  console.log(`[DEBUG] Fetching URL: ${url} with API_KEY: ${API_KEY} and with options:`, options);

  return fetch(url, {
    ...options,
    headers,
    credentials: 'include'
  });
}
export const __PROD__ = false;
export const version = 1;
export const API_URL = __PROD__
  ? "https://kjhlsdfgbnihjodsfbgdf.com/api/v" + version + "/"
  : "http://localhost:8000/api/v" + version;
export const AUTH_URL = __PROD__
  ? "https://kjhlsdfgbnihjodsfbgdf.com/api/auth"
  : "http://localhost:8000/api/auth";

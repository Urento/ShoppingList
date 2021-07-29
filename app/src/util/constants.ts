export const __PROD__ = process.env.NODE_ENV === "production" ? true : false;
export const API_URL = __PROD__ ? "" : "http://localhost:8000/api/v1";
export const AUTH_API_URL = __PROD__ ? "" : "http://localhost:8000/api/auth";
export const AUTH_REGISTER_API_URL = __PROD__
  ? ""
  : "http://localhost:8000/api/auth/register";

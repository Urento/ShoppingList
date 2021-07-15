export const __PROD__ = __DEV__ ? false : true;

export const API_URL = __PROD__
  ? "https://skdhjfbgkhjsdf.ödfgk/api/v1"
  : "http://localhost:8000/api/v1/";

export const AUTH_API_URL = __PROD__
  ? "https://skdhjfbgkhjsdf.ödfgk/api/auth"
  : "http://localhost:8000/api/auth";

export const AUTH_REGISTER_API_URL = __PROD__
  ? "https://skdhjfbgkhjsdf.ödfgk/api/auth/register"
  : "http://localhost:8000/api/auth/register";

export type HomeStackParametersList = {
  Login: undefined;
  Register: undefined;
  Home: undefined;
};

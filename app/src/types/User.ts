export interface UserInfo {
  created_on: number;
  modified_on: number;
  deleted_at: number | null;
  id: number;
  e_mail: string;
  email_verified: boolean;
  username: string;
  rank: string;
  two_factor_authentication: boolean;
}

export interface UserInfoResponse {
  data: UserInfo;
  message: string;
  code: string;
}

export interface UpdateUserData {
  success: "true" | "false";
  message: string;
  error: string;
}

export interface UpdateUserResponse {
  message: string;
  data: UpdateUserData;
  code: number;
}

export interface LoginJSONResponse {
  code: string;
  message: "fail" | "ok";
  data: DataResponse;
}

export interface DataResponse {
  token: string;
  success: "true" | "false";
  otp: boolean;
  error: string;
}

export interface AuthDataResponse {
  success: "true" | "false";
  error: string;
}

export interface AuthCheckResponse {
  code: string;
  message: string;
  data: AuthDataResponse;
}

export interface LogoutDataResponse {
  success: "true" | "false";
  error: string;
}

export interface LogoutResponse {
  message: string;
  code: string;
  data: LogoutDataResponse;
}

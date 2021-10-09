export interface TwoFactorAuthenticationData {
  success: "true" | "false";
  message: string;
  status: "true" | "false";
}

export interface TwoFactorAuthenticationResponse {
  message: string;
  data: TwoFactorAuthenticationData;
  code: number;
}

export interface VerifyResponseData {
  success: "true" | "false";
  message: string;
  verified: "true" | "false";
}

export interface VerifyResponse {
  code: number;
  data: VerifyResponseData;
  message: string;
}

export interface EnableTOTPResponseData {
  success: "true" | "false";
  img: string;
}

export interface EnableTOTPResponse {
  code: number;
  data: EnableTOTPResponseData;
  message: string;
}

export interface TOTPDataResponse {
  token: string;
  success: "true" | "false";
  otp: boolean;
  error: string;
  verified: "true" | "false";
}

export interface TOTPJSONResponse {
  code: string;
  message: "fail" | "ok";
  data: TOTPDataResponse;
}

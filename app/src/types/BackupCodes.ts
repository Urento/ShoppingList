export interface BackupCodesResponse {
  code: string;
  message: string;
  data: BackupCodesResponseData;
}

interface BackupCodesResponseData {
  codes: string;
  success: "true" | "false";
  has: "true" | "false";
}

export interface GenerateBackupCodesResponse {
  code: string;
  message: string;
  data: GenerateBackupCodesResponseData;
}

interface GenerateBackupCodesResponseData {
  codes: string;
  success: "true" | "false";
}

export interface VerifyBackupCode {
  code: number;
  message: string;
  data: VerifyBackupCodeData;
}

interface VerifyBackupCodeData {
  success: "true" | "false";
  ok: "true" | "false";
  error: string;
}

export interface ResetPasswordBackupCode {
  code: number;
  message: string;
  data: ResetPasswordBackupCodeData;
}

interface ResetPasswordBackupCodeData {
  error: string;
  success: "true" | "false";
  ok: "true" | "false";
}

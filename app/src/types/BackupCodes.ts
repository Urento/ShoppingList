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

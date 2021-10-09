export interface ListResponseData {
  error: string;
  success: "true" | "false";
  created_on: number;
  modified_on: number;
  deleted_at: number | null;
  id: number;
  title: string;
  items: string[];
  owner: string;
  participants: string;
}

export interface ListResponse {
  message: string;
  data: ListResponseData;
  code: number;
}

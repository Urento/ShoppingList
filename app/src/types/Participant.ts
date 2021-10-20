export interface Participant {
  created_on?: number;
  id?: number;
  participantId?: number;
  parentListId?: number;
  status?: string;
  email: string;
  request_from?: string;
}

export interface RequestsFromList {
  message: string;
  code: number;
  data: Participant[];
}

export interface AddParticipantResponse {
  message: string;
  code: number;
  data: AddParticipantData;
}

interface AddParticipantData {
  success: "true" | "false";
  error: string;
}

export interface Participant {
  id?: number;
  participantId?: number;
  parentListId?: number;
  status?: string;
  email: string;
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

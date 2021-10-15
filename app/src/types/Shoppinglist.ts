import { Participant } from "./Participant";

export interface CreateItemResponse {
  message: string;
  code: number;
  data: Item;
}

export interface ListResponseData {
  //  error: string;
  //  success: "true" | "false";
  created_on: number;
  modified_on: number;
  deleted_at: number | null;
  id: number;
  title: string;
  items: Item[];
  owner: string;
  participants: Participant[];
}

export interface ListResponse {
  message: string;
  data: ListResponseData;
  code: number;
}

export interface Item {
  id: number;
  parentListId: number;
  itemId: number;
  title: string;
  position: number;
  bought: boolean;
}

export interface Shoppinglist {
  id: number;
  title: string;
  items: Item[];
  owner: string;
  participants: Participant[];
}

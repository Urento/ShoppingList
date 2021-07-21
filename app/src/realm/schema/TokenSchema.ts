export const TokenSchema: Realm.ObjectSchema = {
  name: "Token",
  primaryKey: "_id",
  properties: {
    token: "string",
  },
};

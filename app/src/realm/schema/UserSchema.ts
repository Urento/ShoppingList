export const UserSchema: Realm.ObjectSchema = {
  name: "User",
  primaryKey: "_id",
  properties: {
    email: "string",
    password: "string",
    email_verified: "boolean",
    rank: "string",
  },
};

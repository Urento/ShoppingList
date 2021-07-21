import { realm } from "../App";
import { SchemaTypes } from "./schema/SchemaTypes";

export const getJWTToken = () => {
  console.log(SchemaTypes.TOKEN.toString());
  console.log(realm.objects(SchemaTypes.TOKEN.toString())[0].toJSON().token);
  return realm.objects(SchemaTypes.TOKEN.toString())[0].toJSON().token;
};

export const removeJWTToken = () => {
  const token = realm.objects(SchemaTypes.TOKEN.toString());
  realm.write(() => {
    realm.delete(token);
  });
};

export const setJWTToken = (newToken: string) => {
  const tokens: Realm.Results<Object> = realm.objects(
    SchemaTypes.TOKEN.toString()
  );
  realm.write(() => {
    //delete old token
    realm.delete(tokens);

    //create new token
    realm.create(SchemaTypes.TOKEN.toString(), { token: newToken });
  });
};

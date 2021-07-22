import { NavigationContainer } from "@react-navigation/native";
import React from "react";
import { useEffect } from "react";
import { Provider } from "react-redux";
import Realm from "realm";
import { LoginNavigation } from "./navigation/TabNavigation";
import { ShoppinglistSchema } from "./realm/schema/ShoppinglistSchema";
import { TokenSchema } from "./realm/schema/TokenSchema";
import { UserSchema } from "./realm/schema/UserSchema";
import store from "./redux/store";
import { CheckAuthentication } from "./util/CheckAuthentication";

const DarkTheme = {
  dark: true,
  colors: {
    primary: "rgb(255, 45, 85)",
    background: "rgb(242, 242, 242)",
    card: "rgb(20,29,38)",
    text: "rgb(28, 28, 30)",
    border: "rgb(199, 199, 204)",
    notification: "rgb(255, 69, 58)",
  },
};

/**
 * TODO: Implement React Query
 * TODO: Implement Authenticated Only Screens
 * TODO: Add Schemas more elegant
 * TODO: Add Tests
 * https://realm.io/
 */

export let realm: Realm;

export const App = () => {
  //TODO: Dont create a new one everytime just open
  realm = new Realm({ schema: [TokenSchema, ShoppinglistSchema, UserSchema] });

  useEffect(() => {
    CheckAuthentication();
  }, []);

  return (
    <Provider store={store}>
      <NavigationContainer theme={DarkTheme}>
        <LoginNavigation />
      </NavigationContainer>
    </Provider>
  );
};

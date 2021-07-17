import { NavigationContainer } from "@react-navigation/native";
import React from "react";
import { Provider } from "react-redux";
import { LoginNavigation } from "./navigation/TabNavigation";
import store from "./redux/store";

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
 */

export const App = () => {
  return (
    <Provider store={store}>
      <NavigationContainer theme={DarkTheme}>
        <LoginNavigation />
      </NavigationContainer>
    </Provider>
  );
};

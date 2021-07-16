import { NavigationContainer } from "@react-navigation/native";
import { createStackNavigator } from "@react-navigation/stack";
import React from "react";
import TabNavigation from "./navigation/TabNavigation";

export const App = () => {
  return (
    <NavigationContainer>
      <TabNavigation loggedIn={false} />
    </NavigationContainer>
  );
};

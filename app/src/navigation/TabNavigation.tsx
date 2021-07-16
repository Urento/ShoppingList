import React from "react";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import HomeScreen from "../HomeScreen";
import { LoginScreen } from "../LoginScreen";
import RegisterScreen from "../RegisterScreen";

const Tab = createBottomTabNavigator();

interface IProps {
  loggedIn: boolean;
}

//https://reactnavigation.org/docs/tab-based-navigation

/**
 * distinguish if logged in or not
 */

const TabNavigation: React.FC<IProps> = ({ loggedIn }) => {
  return (
    <Tab.Navigator>
      {loggedIn ? (
        <>
          <Tab.Screen name="Home" component={HomeScreen} />
          <Tab.Screen name="Settings" component={LoginScreen} />
        </>
      ) : (
        <>
          <Tab.Screen name="Login" component={LoginScreen} />
          <Tab.Screen name="Register" component={RegisterScreen} />
        </>
      )}
    </Tab.Navigator>
  );
};

export default TabNavigation;

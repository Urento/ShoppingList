import React from "react";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import HomeScreen from "../HomeScreen";
import { LoginScreen } from "../LoginScreen";
import RegisterScreen from "../RegisterScreen";
import Ionicons from "react-native-vector-icons/Ionicons";
import SettingsScreeen from "../SettingsScreen";

const Tab = createBottomTabNavigator();

interface IProps {}

//https://reactnavigation.org/docs/tab-based-navigation

export const LoginNavigation: React.FC<IProps> = () => {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: string | undefined;

          //TODO: differentiate between ios and android
          switch (route.name) {
            case "Login":
              iconName = focused
                ? "ios-person-circle"
                : "ios-person-circle-outline";
              break;
            case "Register":
              iconName = focused
                ? "ios-person-circle"
                : "ios-person-circle-outline";
              break;
            default:
              break;
          }

          //@ts-ignore
          return <Ionicons name={iconName} size={size} color={color} />;
        },
      })}
      tabBarOptions={{
        activeTintColor: "tomato",
        inactiveTintColor: "gray",
      }}
    >
      <>
        <Tab.Screen name="_" component={HomeScreen} />
        <Tab.Screen name="Login" component={LoginScreen} />
        <Tab.Screen name="Register" component={RegisterScreen} />
      </>
    </Tab.Navigator>
  );
};

export const LoggedInNavigation: React.FC<IProps> = ({}) => {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: string | undefined;

          //TODO: differentiate between ios and android
          switch (route.name) {
            case "Home":
              iconName = focused ? "ios-home-outline" : "ios-home-outline";
              break;
            case "Settings":
              iconName = focused
                ? "ios-settings-outline"
                : "ios-settings-outline";
              break;
            default:
              break;
          }

          //@ts-ignore
          return <Ionicons name={iconName} size={size} color={color} />;
        },
      })}
      tabBarOptions={{
        activeTintColor: "tomato",
        inactiveTintColor: "gray",
      }}
    >
      <>
        <Tab.Screen name="Home" component={HomeScreen} />
        <Tab.Screen name="Login" component={LoginScreen} />
        <Tab.Screen name="Register" component={RegisterScreen} />
      </>
    </Tab.Navigator>
  );
};

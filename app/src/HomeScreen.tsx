import { StackNavigationProp } from "@react-navigation/stack";
import { createDrawerNavigator } from "@react-navigation/drawer";
import React from "react";
import { HomeStackParametersList } from "./util/constants";
import { View } from "react-native";
import { Text } from "react-native-elements";

interface IProps {
  navigation: StackNavigationProp<HomeStackParametersList>;
}

/**
 * TODO: Add Sidemenu
 * TODO: Add View of all Shoppinglists
 * TODO: Maybe save the shoppinglist you where last on and instant redirect to that
 */

const HomeScreen: React.FC<IProps> = ({ navigation }) => {
  return <Text>Home Screen</Text>;
};

export default HomeScreen;

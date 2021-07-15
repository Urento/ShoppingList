import { StackNavigationProp } from "@react-navigation/stack";
import React from "react";
import { Text } from "react-native-elements";
import { HomeStackParametersList } from "./util/constants";

interface IProps {
  navigation: StackNavigationProp<HomeStackParametersList>;
}

/**
 * TODO: Add Sidemenu
 * TODO: Add View of all Shoppinglists
 * TODO: Maybe save the shoppinglist you where last on and instant redirect to that
 */

const HomeScreen: React.FC<IProps> = ({ navigation }) => {
  return <Text>Test</Text>;
};

export default HomeScreen;

import { StackNavigationProp } from "@react-navigation/stack";
import React from "react";
import { Text } from "react-native-elements";
import { HomeStackParametersList } from "./util/constants";

interface IProps {
  navigation: StackNavigationProp<HomeStackParametersList>;
}

const SettingsScreeen: React.FC<IProps> = ({ navigation }) => {
  return <Text>Settings</Text>;
};

export default SettingsScreeen;

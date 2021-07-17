import { StackNavigationProp } from "@react-navigation/stack";
import React from "react";
import { HomeStackParametersList } from "./util/constants";
import { Text } from "react-native-elements";
import { SafeAreaView } from "react-native-safe-area-context";

interface IProps {
  navigation: StackNavigationProp<HomeStackParametersList>;
}

/**
 * TODO: Add View of all Shoppinglists
 * TODO: Maybe save the shoppinglist you where last on and instant redirect to that
 */

const HomeScreen: React.FC<IProps> = ({ navigation }) => {
  return (
    <SafeAreaView>
      <Text>This is top text.</Text>
      <Text>This is bottom text.</Text>
    </SafeAreaView>
  );
};

export default HomeScreen;

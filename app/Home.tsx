import React from "react";
import { Text, StyleSheet, View, TextInput, Button } from "react-native";

export type Props = {
  loggedIn: boolean;
};

const Home: React.FC<Props> = ({ loggedIn = true }) => {
  return (
    <View style={styles.container}>
      <Text style={styles.formLabel}> Home </Text>
      <View>
        <Text>ösydkjhfgbnhdsflkhgdgf</Text>
      </View>
    </View>
  );
};
const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#141d26",
    color: "#fff",
    alignItems: "center",
    justifyContent: "center",
    height: 50,
  },

  formLabel: {
    fontSize: 20,
    color: "#fff",
  },
  inputStyle: {
    marginTop: 20,
    width: 300,
    height: 40,
    paddingHorizontal: 10,
    borderRadius: 50,
    backgroundColor: "#DCDCDC",
  },
  formText: {
    alignItems: "center",
    justifyContent: "center",
    color: "#fff",
    fontSize: 20,
  },
  text: {
    color: "#fff",
    fontSize: 20,
  },
  button: {
    marginTop: 20,
    width: 300,
    height: 40,
    paddingHorizontal: 10,
    borderRadius: 50,
    backgroundColor: "#000",
  },
});
export default Home;

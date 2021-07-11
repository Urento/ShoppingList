import React, { useState } from "react";
import { Text, StyleSheet, View, TextInput } from "react-native";
import { AUTH_URL } from "../../util/constants";
import { Button } from "react-native-elements";

export type Props = {
  loggedIn: boolean;
};

//https://reactnativeelements.com/docs/button/

const Login: React.FC<Props> = ({ loggedIn = false }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const loginRequest = async () => {
    return fetch(AUTH_URL, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });
  };

  const login = async (e: any) => {
    e.preventDefault();
    try {
      const resJson = await (await loginRequest()).json();
      console.log(resJson);
    } catch (err) {
      console.log("err: " + err);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.formLabel}>Anmelden</Text>
      <View>
        <TextInput
          onChangeText={(e) => setEmail(e)}
          placeholder="Email oder Benutzername"
          style={styles.inputStyle}
          keyboardType="email-address"
          value={email}
        />
        <TextInput
          onChangeText={(password) => setPassword(password)}
          secureTextEntry={true}
          placeholder="Passwort"
          style={styles.inputStyle}
          value={password}
        />
        <View style={{ paddingTop: "5%" }}>
          <Button
            onPress={login}
            title="LOGIN"
            accessibilityLabel="Login Button"
          />
        </View>
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
    fontSize: 16,
    lineHeight: 21,
    fontWeight: "bold",
    letterSpacing: 0.25,
    color: "black",
  },
});

export default Login;

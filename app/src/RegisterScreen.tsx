import React from "react";
import { StyleSheet, View } from "react-native";
import Icon from "react-native-vector-icons/FontAwesome";
import { Input, Text } from "react-native-elements";
import { Button } from "react-native-elements/dist/buttons/Button";
import { useState } from "react";
import { StackNavigationProp } from "@react-navigation/stack";
import { validateEmail } from "./util/validate";
import {
  AUTH_REGISTER_API_URL,
  HomeStackParametersList,
} from "./util/constants";

interface IProps {
  navigation: StackNavigationProp<HomeStackParametersList>;
}

const RegisterScreen: React.FC<IProps> = ({ navigation }) => {
  const [emailError, setEmailError] = useState(false);
  const [pwdError, setPwdError] = useState(false);
  const [usernameError, setUsernameError] = useState(false);
  const [error, setError] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [loading, setLoading] = useState(false);

  const register = async () => {
    setLoading(true);
    const validEmail = validateEmail(email);

    if (!validEmail) {
      setEmailError(true);
      setLoading(false);
      return setError("Email not valid!");
    }

    if (password.length < 8) {
      setPwdError(true);
      setLoading(false);
      return setError("Password has to be atleast 8 characters long!");
    }

    if (username.length > 32) {
      setUsernameError(true);
      setLoading(false);
      return setError("Username can't be longer than 32 characters!");
    }

    console.log(AUTH_REGISTER_API_URL);
    const loginRequest = await fetch(AUTH_REGISTER_API_URL, {
      method: "POST",
      body: JSON.stringify({
        email: email,
        username: username,
        password: password,
      }),
    });

    const resJson = await loginRequest.json();
    console.log(loginRequest.status);
    console.log(resJson);
    setLoading(false);
  };

  const redirectToLogin = () => {
    const navigate = navigation.navigate;
    navigate("Login");
  };

  return (
    <View style={styles.container}>
      <Input
        placeholder="Email"
        onChangeText={(email) => setEmail(email)}
        leftIcon={<Icon name="envelope" size={24} color="white" />}
        style={{ color: "#fff" }}
        errorStyle={emailError ? { color: "red" } : {}}
        errorMessage={emailError ? error : ""}
        keyboardType="email-address"
      />
      <Input
        placeholder="Username"
        onChangeText={(username) => setUsername(username)}
        leftIcon={<Icon name="user" size={24} color="white" />}
        style={{ color: "#fff" }}
        errorStyle={usernameError ? { color: "red" } : {}}
        errorMessage={usernameError ? error : ""}
      />
      <Input
        placeholder="Password"
        onChangeText={(password) => setPassword(password)}
        leftIcon={<Icon name="lock" size={24} color="white" />}
        secureTextEntry={true}
        style={{ color: "#fff" }}
        errorStyle={pwdError ? { color: "red" } : {}}
        errorMessage={pwdError ? error : ""}
      />
      {loading ? (
        <Button onPress={register} title="REGISTER" loading />
      ) : (
        <Button onPress={register} title="REGISTER" />
      )}
      <Text style={styles.registerText} onPress={redirectToLogin}>
        Already have an account? Login!
      </Text>
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
  registerText: {
    color: "#fff",
  },
});

export default RegisterScreen;

import { StackNavigationProp } from "@react-navigation/stack";
import React, { Component } from "react";
import { Button, StyleSheet, View } from "react-native";
import Icon from "react-native-vector-icons/FontAwesome";
import { Input } from "react-native-elements";
import { AUTH_API_URL } from "./util/constants";
import { validateEmail } from "./util/validate";

type HomeStackParametersList = {
  Login: undefined;
};

interface Props {
  navigation: StackNavigationProp<HomeStackParametersList>;
}

interface IState {
  email: string;
  password: string;
  emailError: boolean;
  pwdError: boolean;
  error: string;
}

export class LoginScreen extends Component<Props, IState> {
  constructor(props: any) {
    super(props);
    this.state = {
      email: "",
      password: "",
      emailError: false,
      pwdError: false,
      error: "",
    };
  }

  private login = async (e: any) => {
    e.preventDefault();
    const { email } = this.state;
    const { password } = this.state;
    const validEmail = validateEmail(email);
    console.log(email);
    console.log(password);
    console.log(validEmail);

    if (!validEmail) {
      this.setState({
        emailError: true,
        error: "Email not valid",
      });
      return;
    }

    const navigate = this.props.navigation.navigate;
    const loginRequest = await fetch(AUTH_API_URL, {
      method: "POST",
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });
    const resJson = await loginRequest.json();
  };

  public render() {
    const { emailError, error } = this.state;
    return (
      <View style={styles.container}>
        {!emailError ? (
          <Input
            placeholder="Email"
            onChangeText={(email) => this.setState({ email: email })}
            leftIcon={<Icon name="envelope" size={24} color="white" />}
            style={{ color: "#fff" }}
          />
        ) : (
          <Input
            placeholder="Email"
            onChangeText={(email) => this.setState({ email: email })}
            leftIcon={<Icon name="envelope" size={24} color="white" />}
            style={{ color: "#fff" }}
            errorStyle={{ color: "red" }}
            errorMessage={error}
          />
        )}
        <Input
          placeholder="Password"
          onChangeText={(password) => this.setState({ password: password })}
          leftIcon={<Icon name="lock" size={24} color="white" />}
          secureTextEntry={true}
          style={{ color: "#fff" }}
        />
        <Button onPress={this.login} title={"Login"} />
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#141d26",
    color: "#fff",
    alignItems: "center",
    justifyContent: "center",
    height: 50,
  },
});

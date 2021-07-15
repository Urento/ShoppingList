import { StackNavigationProp } from "@react-navigation/stack";
import React, { Component } from "react";
import { StyleSheet, View } from "react-native";
import Icon from "react-native-vector-icons/FontAwesome";
import { Input, Text } from "react-native-elements";
import { AUTH_API_URL, HomeStackParametersList } from "./util/constants";
import { validateEmail } from "./util/validate";
import { Button } from "react-native-elements/dist/buttons/Button";

interface Props {
  navigation: StackNavigationProp<HomeStackParametersList>;
}

interface IState {
  email: string;
  password: string;
  emailError: boolean;
  pwdError: boolean;
  error: string;
  loading: boolean;
}

export class LoginScreen extends Component<Props, IState> {
  constructor(props: Props) {
    super(props);
    this.state = {
      email: "",
      password: "",
      emailError: false,
      pwdError: false,
      error: "",
      loading: false,
    };
  }

  private redirectToRegister = () => {
    const navigate = this.props.navigation.navigate;
    navigate("Register");
  };

  private login = async () => {
    this.setState({ loading: true });
    const { email } = this.state;
    const { password } = this.state;
    const validEmail = validateEmail(email);

    if (!validEmail) {
      return this.setState({
        emailError: true,
        pwdError: false,
        error: "Email not valid",
        loading: false,
      });
    }

    if (password.length < 8) {
      return this.setState({
        pwdError: true,
        emailError: false,
        error: "Password has to be atleast 8 characters long!",
        loading: false,
      });
    }

    //const navigate = this.props.navigation.navigate;
    console.log(AUTH_API_URL);
    const loginRequest = await fetch(AUTH_API_URL, {
      method: "POST",
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });
    const resJson = await loginRequest.json();
    console.log(loginRequest.status);
    console.log(resJson);
    this.setState({ loading: false });
  };

  public render() {
    const { emailError, error, pwdError, loading } = this.state;
    return (
      <View style={styles.container}>
        <Input
          placeholder="Email"
          onChangeText={(email) => this.setState({ email: email })}
          leftIcon={<Icon name="envelope" size={24} color="white" />}
          style={{ color: "#fff" }}
          errorStyle={emailError ? { color: "red" } : {}}
          errorMessage={emailError ? error : ""}
          keyboardType="email-address"
        />
        <Input
          placeholder="Password"
          onChangeText={(password) => this.setState({ password: password })}
          leftIcon={<Icon name="lock" size={24} color="white" />}
          secureTextEntry={true}
          style={{ color: "#fff" }}
          errorStyle={pwdError ? { color: "red" } : {}}
          errorMessage={pwdError ? error : ""}
        />
        {loading ? (
          <Button onPress={this.login} title="LOGIN" loading />
        ) : (
          <Button onPress={this.login} title="LOGIN" />
        )}
        <Text style={styles.registerText} onPress={this.redirectToRegister}>
          Dont have an Account yet? Make one!
        </Text>
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
  registerText: {
    color: "#fff",
  },
});

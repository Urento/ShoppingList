import React, { useState } from "react";
import {
  Text,
  StyleSheet,
  View,
  TextInput,
  Modal,
  Alert,
  Pressable,
} from "react-native";
import { AUTH_URL } from "../../util/constants";
import { Button } from "react-native-elements";

export type Props = {
  loggedIn: boolean;
};

//https://reactnativeelements.com/docs/button/

const Login: React.FC<Props> = ({ loggedIn = false }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [error, setError] = useState("");

  const login = async (e: any) => {
    e.preventDefault();

    setLoading(true);
    try {
      let response = await fetch(AUTH_URL, {
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
      let resJson = await response.json();
      if (response.status != 200) {
        setError("Error while logging in [" + response.status + "]");
        setModalVisible(true);
      }
      console.log(resJson);
      if (resJson != null) setLoading(false);
      if (resJson.code != 200) {
        setError("Error while logging in [" + resJson.code + "]");
        setModalVisible(true);
      }
    } catch (err) {
      setError(err.toString());
      setModalVisible(true);
      console.log("err: " + err);
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.formLabel}>Anmelden</Text>
      {modalVisible && (
        <Modal
          animationType="slide"
          transparent={true}
          visible={modalVisible}
          onRequestClose={() => setModalVisible(!modalVisible)}
        >
          <View style={styles.centeredView}>
            <View style={styles.modalView}>
              <Text style={styles.modalText}>
                Error while logging in: {error}!
              </Text>
              <Pressable
                style={[styles.button, styles.buttonClose]}
                onPress={() => setModalVisible(!modalVisible)}
              >
                <Text style={styles.textStyle}>Try Again</Text>
              </Pressable>
            </View>
          </View>
        </Modal>
      )}
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
          {!loading ? (
            <Button
              onPress={login}
              title="LOGIN"
              accessibilityLabel="Login Button"
            />
          ) : (
            <Button title="LOGIN" accessibilityLabel="Login Button" loading />
          )}
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
  centeredView: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    marginTop: 22,
  },
  modalView: {
    margin: 20,
    backgroundColor: "white",
    borderRadius: 20,
    padding: 35,
    alignItems: "center",
    shadowColor: "#000",
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 4,
    elevation: 5,
  },
  button: {
    borderRadius: 20,
    padding: 10,
    elevation: 2,
  },
  buttonOpen: {
    backgroundColor: "#F194FF",
  },
  buttonClose: {
    backgroundColor: "#2196F3",
  },
  textStyle: {
    color: "white",
    fontWeight: "bold",
    textAlign: "center",
  },
  modalText: {
    marginBottom: 15,
    textAlign: "center",
  },
});

export default Login;

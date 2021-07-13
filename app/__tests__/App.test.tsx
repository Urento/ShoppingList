import "react-native";
import React from "react";
import Login from "../components/form/Login";

// Note: test renderer must be required after react-native.
import renderer from "react-test-renderer";

describe("Login Screen", () => {
  test("Render Login Screen", () => {
    renderer.create(<Login loggedIn={false} />);
  });

  test("renders correctly", () => {
    const tree = renderer.create(<Login loggedIn={false} />).toJSON();
    expect(tree).toMatchSnapshot();
  });
});

import "react-native";
import React from "react";
import App from "../App";

// Note: test renderer must be required after react-native.
import renderer from "react-test-renderer";

describe("Login Screen", () => {
  test("Render Login Screen", () => {
    renderer.create(<App loggedIn={false} />);
  });

  test("renders correctly", () => {
    const tree = renderer.create(<App loggedIn={false} />).toJSON();
    expect(tree).toMatchSnapshot();
  });
});

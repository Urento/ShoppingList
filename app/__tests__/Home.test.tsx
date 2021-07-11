import "react-native";
import React from "react";
import Home from "../Home";

// Note: test renderer must be required after react-native.
import renderer from "react-test-renderer";

describe("Home Screen", () => {
  test("Render Home Screen", () => {
    renderer.create(<Home loggedIn={true} />);
  });

  test("renders correctly", () => {
    const tree = renderer.create(<Home loggedIn={true} />).toJSON();
    expect(tree).toMatchSnapshot();
  });
});

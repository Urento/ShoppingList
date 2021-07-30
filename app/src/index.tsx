import ReactDOM from "react-dom";
import "./index.css";
import { BrowserRouter, Switch, Route } from "react-router-dom";
import Login from "./App";
import Register from "./Register";

//TODO: add auth stuff and check exp time on jwt token

ReactDOM.render(
  <BrowserRouter>
    <Switch>
      <Route exact path="/register" component={Register} />
      <Route exact path="/" component={Login} />
    </Switch>
  </BrowserRouter>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
// serviceWorker.unregister();

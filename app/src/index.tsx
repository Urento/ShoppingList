import ReactDOM from "react-dom";
import "./index.css";
import { BrowserRouter, Switch, Route, Redirect } from "react-router-dom";
import Login from "./App";
import Register from "./Register";
import { Dashboard } from "./screens/Dashboard";
import { QueryClient, QueryClientProvider } from "react-query";
import { Settings } from "./screens/Settings";
import { NotFound } from "./screens/NotFound";

//TODO: Store User Info in Redux

const queryClient = new QueryClient();
ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <Switch>
        <Route exact path="/register" component={Register} />
        <Route exact path="/" component={Login} />

        <Route exact path="/dashboard" component={Dashboard} />
        <Route exact path="/settings" component={Settings} />

        <Route exact path="/404" component={NotFound} />
        <Redirect to="/404" />
      </Switch>
    </BrowserRouter>
  </QueryClientProvider>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
// serviceWorker.unregister();

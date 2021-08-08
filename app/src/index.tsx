import ReactDOM from "react-dom";
import "./index.css";
import { BrowserRouter, Switch, Route, Redirect } from "react-router-dom";
import Login from "./App";
import Register from "./Register";
import { Dashboard } from "./screens/Dashboard";
import { QueryClient, QueryClientProvider } from "react-query";
import { Settings } from "./screens/Settings";
import { NotFound } from "./screens/NotFound";
import { NewShoppinglist } from "./screens/shoppinglist/NewShoppinglist";
import { ViewShoppinglist } from "./screens/shoppinglist/ViewShoppinglist";
import { TwoFactorAuthentication } from "./TwoFactorAuthentication";
import { UpdateTwoFactorAuthentication } from "./screens/settings/UpdateTwoFactorAuthentication";

export const queryClient = new QueryClient();
ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <Switch>
        <Route exact path="/register" component={Register} />
        <Route exact path="/" component={Login} />
        <Route
          exact
          path="/twofactorauthentication"
          component={TwoFactorAuthentication}
        />

        <Route exact path="/dashboard" component={Dashboard} />
        <Route exact path="/settings" component={Settings} />
        <Route
          exact
          path="/settings/totp"
          component={UpdateTwoFactorAuthentication}
        />
        <Route exact path="/lists/create" component={NewShoppinglist} />
        <Route path="/list/:id" component={ViewShoppinglist} />

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

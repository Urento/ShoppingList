import ReactDOM from "react-dom";
import "./index.css";
import { BrowserRouter, Switch, Route, Redirect } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "react-query";
import { lazy, Suspense } from "react";
import { Loading } from "./components/Loading";

const RegisterComponent = lazy(() => import("./Register"));
const LoginComponent = lazy(() => import("./App"));
const TwoFactorAuthenticationComponent = lazy(
  () => import("./TwoFactorAuthentication")
);
const UpdateTwoFactorAuthenticationComponent = lazy(
  () => import("./screens/settings/UpdateTwoFactorAuthentication")
);
const DashboardComponent = lazy(() => import("./screens/Dashboard"));
const SettingsComponent = lazy(() => import("./screens/Settings"));
const NewShoppinglistComponent = lazy(
  () => import("./screens/shoppinglist/NewShoppinglist")
);
const ViewShoppinglistComponent = lazy(
  () => import("./screens/shoppinglist/ViewShoppinglist")
);
const NotFoundComponent = lazy(() => import("./screens/NotFound"));
const NotificationsComponent = lazy(
  () => import("./screens/notifications/Notifications")
);
const BackupCodesComponent = lazy(
  () => import("./screens/backupcodes/BackupCodes")
);

//TODO: Cache

export const queryClient = new QueryClient();
ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <Suspense fallback={<Loading />}>
        <Switch>
          <Route exact path="/register" component={RegisterComponent} />
          <Route exact path="/" component={LoginComponent} />
          <Route
            exact
            path="/twofactorauthentication"
            component={TwoFactorAuthenticationComponent}
          />

          <Route exact path="/dashboard" component={DashboardComponent} />
          <Route exact path="/settings" component={SettingsComponent} />
          <Route
            exact
            path="/settings/totp"
            component={UpdateTwoFactorAuthenticationComponent}
          />
          <Route
            exact
            path="/lists/create"
            component={NewShoppinglistComponent}
          />
          <Route path="/list/:id" component={ViewShoppinglistComponent} />

          <Route path="/notifications" component={NotificationsComponent} />

          <Route path="/backupcodes" component={BackupCodesComponent} />

          <Route exact path="/404" component={NotFoundComponent} />
          <Redirect to="/404" />
        </Switch>
      </Suspense>
    </BrowserRouter>
  </QueryClientProvider>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
// serviceWorker.unregister();

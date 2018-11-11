import Layout from './layout.js';
import IndexPage from './pages/index.js';
import FormsPage from './pages/forms/index.js';

const e = React.createElement;
const Router = ReactRouter.Router;
const Switch = ReactRouter.Switch;
const Route = ReactRouter.Route;
const Redirect = ReactRouter.Redirect;
const createBrowserHistory = History.createBrowserHistory;

export default function AppComponent({ redirect = null}) {
    return e(Router, { history: createBrowserHistory() },
        e(Layout, null, [
            e(Switch, null, [
                e(Route, {
                    path: '/',
                    exact: true,
                    component: IndexPage
                }),
                e(Route, {
                    path: '/forms',
                    component: FormsPage
                })
            ]),
            redirect ? e(Redirect, { to: redirect }) : null
        ])
    )
}
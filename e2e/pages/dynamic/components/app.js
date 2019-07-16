import Layout from './layout.js';
import IndexPage from './pages/index.js';
import FormsPage from './pages/forms/index.js';
import EventsPage from './pages/events/index.js';
import IframePage from './pages/iframes/index.js';
import MediaPage from './pages/media/index.js';

const e = React.createElement;
const Router = ReactRouter.Router;
const Switch = ReactRouter.Switch;
const Route = ReactRouter.Route;
const Redirect = ReactRouter.Redirect;
const createBrowserHistory = History.createBrowserHistory;

export default React.memo(function AppComponent(params = {}) {
    let redirectTo;

    if (params.redirect) {
        let search = '';

        Object.keys(params).forEach((key) => {
            if (key !== 'redirect') {
                search += `${key}=${params[key]}`;
            }
        });

        const to = {
            pathname: params.redirect,
            search: search ? `?${search}` : '',
        };

        redirectTo = e(Redirect, { to });
    }

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
                }),
                e(Route, {
                    path: '/events',
                    component: EventsPage
                }),
                e(Route, {
                    path: '/iframe',
                    component: IframePage
                }),
                e(Route, {
                    path: '/media',
                    component: MediaPage
                }),
            ]),
            redirectTo
        ])
    )
})
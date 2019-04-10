const e = React.createElement;
const NavLink = ReactRouterDOM.NavLink;

export default function Layout({ children }) {
    return e("div", { id: "layout"}, [
        e("nav", { className: "navbar navbar-expand-md navbar-dark bg-dark mb-4", id: "navbar" }, [
            e(NavLink, { className: "navbar-brand", to: "/"}, "Ferret"),
            e("button", { className: "navbar-toggler", type: "button"}, [
                e("span", { className: "navbar-toggler-icon" })
            ]),
            e("div", { className: "collapse navbar-collapse" }, [
                e("ul", { className: "navbar-nav mr-auto" }, [
                    e("li", { className: "nav-item"}, [
                        e(NavLink, { className: "nav-link", to: "/forms" }, "Forms")
                    ]),
                    e("li", { className: "nav-item"}, [
                        e(NavLink, { className: "nav-link", to: "/navigation" }, "Navigation")
                    ]),
                    e("li", { className: "nav-item"}, [
                        e(NavLink, { className: "nav-link", to: "/events" }, "Events")
                    ])
                ])
            ])
        ]),
        e("main", { className: "container"}, children)
    ])
}
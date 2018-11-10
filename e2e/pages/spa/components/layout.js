const e = React.createElement;
const NavLink = ReactRouterDOM.NavLink;

// <nav class="navbar navbar-expand-md navbar-dark bg-dark mb-4">
//     <a class="navbar-brand" href="#">Top navbar</a>
// <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
//     <span class="navbar-toggler-icon"></span>
//     </button>
//     <div class="collapse navbar-collapse" id="navbarCollapse">
//     <ul class="navbar-nav mr-auto">
//     <li class="nav-item active">
//     <a class="nav-link" href="#">Home <span class="sr-only">(current)</span></a>
// </li>
// <li class="nav-item">
//     <a class="nav-link" href="#">Link</a>
//     </li>
//     <li class="nav-item">
//     <a class="nav-link disabled" href="#">Disabled</a>
//     </li>
//     </ul>
//     <form class="form-inline mt-2 mt-md-0">
//     <input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search">
//     <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
//     </form>
//     </div>
//     </nav>

export default function Layout({ children }) {
    return e("div", { className: "cover-container d-flex w-100 h-100 p-3 mx-auto flex-column"}, [
        e("header", null, [
            e("div", { className: "inner" }, [
                e("h3", { className: "masthead-brand"}, "Ferret"),
                e("nav", {className: "nav nav-masthead justify-content-center"}, [
                    e(NavLink, { className: "nav-link", to: "/forms"}, "Forms"),
                    e(NavLink, { className: "nav-link", to: "/navigation"}, "Navigation"),
                    e(NavLink, { className: "nav-link", to: "/events"}, "Events")
                ])
            ])
        ]),
        e("main", { className: "inner"}, children)
    ])
}
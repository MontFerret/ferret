const e = React.createElement;

export default function IndexPage() {
    return e("div", { className: "jumbotron" }, [
        e("div", null, e("h1", null, "Welcome to Ferret E2E test page!")),
        e("div", null, e("p", { className: "lead" }, "It has several pages for testing different possibilities of the library"))
    ])
}
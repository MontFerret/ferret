const e = React.createElement;
const Fragment = React.Fragment;

export default function IndexPage() {
    return e(Fragment, null, [
        e("h1", { className: "cover-heading" }, "Welcome to Ferret E2E test page!"),
        e("p", { className: "lead" }, "It has several pages for testing different possibilities of the library")
    ])
}
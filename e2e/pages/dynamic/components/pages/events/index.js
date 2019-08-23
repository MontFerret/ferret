import Hoverable from "./hoverable.js";
import Clickable from "./clickable.js";
import Appearable from "./appearable.js";
import Focusable from "./focusable.js";

const e = React.createElement;

export default class EventsPage extends React.Component {
    render() {
        return e("div", { id: "page-events" }, [
            e("div", { className: "row" }, [
                e("div", { className: "col-lg-4"}, [
                    e(Hoverable),
                ]),
                e("div", { id: "wait-class", className: "col-lg-4"}, [
                    e(Clickable, {
                        id: "wait-class",
                        title: "Add class"
                    })
                ]),
                e("div", { id: "wait-class-random", className: "col-lg-4"}, [
                    e(Clickable, {
                        id: "wait-class-random",
                        title: "Add class 2",
                        randomTimeout: true
                    })
                ])
            ]),
            e("div", { className: "row" }, [
                e("div", { id: "wait-no-class", className: "col-lg-4"}, [
                    e(Clickable, {
                        id: "wait-no-class",
                        title: "Remove class",
                        show: true
                    })
                ]),
                e("div", { id: "wait-no-class-random", className: "col-lg-4"}, [
                    e(Clickable, {
                        id: "wait-no-class-random",
                        title: "Remove class 2",
                        show: true,
                        randomTimeout: true
                    })
                ]),
                e("div", { id: "wait-element", className: "col-lg-4"}, [
                    e(Appearable, {
                        id: "wait-element",
                        appear: true,
                        title: "Appearable"
                    })
                ]),
            ]),
            e("div", { className: "row" }, [
                e("div", { id: "wait-no-element", className: "col-lg-4"}, [
                    e(Appearable, {
                        id: "wait-no-element",
                        appear: false,
                        title: "Disappearable"
                    })
                ]),
                e("div", { id: "wait-style", className: "col-lg-4"}, [
                    e(Appearable, {
                        id: "wait-style",
                        appear: true,
                        title: "Appearable with style",
                        useStyle: true,
                    })
                ]),
                e("div", { id: "wait-no-style", className: "col-lg-4"}, [
                    e(Appearable, {
                        id: "wait-no-style",
                        appear: false,
                        title: "Disappearable",
                        useStyle: true,
                    })
                ]),
            ]),
            e("div", { className: "row" }, [
                e("div", { id: "focus", className: "col-lg-4"}, [
                    e(Focusable, {
                        id: "focus",
                        appear: false,
                        title: "Focusable"
                    })
                ]),
            ])
        ])
    }
}
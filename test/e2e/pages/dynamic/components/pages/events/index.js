import Hoverable from "./hoverable.js";
import Clickable from "./clickable.js";
import Appearable from "./appearable.js";
import Focusable from "./focusable.js";
import Pressable from "./pressable.js";
import Ajax from "./ajax.js";

const e = React.createElement;

export default class EventsPage extends React.Component {
    render() {
        return e("div", { id: "page-events" }, [
            e("div", { className: "row" }, [
                e("div", { className: "col-lg-4" }, [
                    e(Hoverable),
                ]),
                e("div", { className: "col-lg-4" }, [
                    e(Clickable, {
                        id: "wait-class",
                        title: "Add class"
                    })
                ]),
                e("div", { className: "col-lg-4" }, [
                    e(Clickable, {
                        id: "wait-class-random",
                        title: "Add class 2",
                        randomTimeout: true
                    })
                ])
            ]),
            e("div", { className: "row" }, [
                e("div", { className: "col-lg-4" }, [
                    e(Clickable, {
                        id: "wait-no-class",
                        title: "Remove class",
                        show: true
                    })
                ]),
                e("div", { className: "col-lg-4" }, [
                    e(Clickable, {
                        id: "wait-no-class-random",
                        title: "Remove class 2",
                        show: true,
                        randomTimeout: true
                    })
                ]),
                e("div", { className: "col-lg-4" }, [
                    e(Appearable, {
                        id: "wait-element",
                        appear: true,
                        title: "Appearable"
                    })
                ]),
            ]),
            e("div", { className: "row" }, [
                e("div", { className: "col-lg-4" }, [
                    e(Appearable, {
                        id: "wait-no-element",
                        appear: false,
                        title: "Disappearable"
                    })
                ]),
                e("div", { className: "col-lg-4" }, [
                    e(Appearable, {
                        id: "wait-style",
                        appear: true,
                        title: "Appearable with style",
                        useStyle: true,
                    })
                ]),
                e("div", { className: "col-lg-4" }, [
                    e(Appearable, {
                        id: "wait-no-style",
                        appear: false,
                        title: "Disappearable",
                        useStyle: true,
                    })
                ]),
            ]),
            e("div", { className: "row" }, [
                e("div", { className: "col-lg-4" }, [
                    e(Focusable, {
                        id: "focus",
                        appear: false,
                        title: "Focusable"
                    })
                ]),

                e("div", { className: "col-lg-4" }, [
                    e(Pressable, {
                        id: "press",
                        appear: false,
                        title: "Pressable"
                    })
                ]),

                e("div", { className: "col-lg-4" }, [
                    e(Ajax, {
                        id: "ajax",
                        appear: false,
                        title: "Requests"
                    })
                ]),
            ]),
        ])
    }
}
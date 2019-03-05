import Hoverable from "./hoverable.js";
import Clickable from "./clickable.js";
import Appearable from "./appearable.js";

const e = React.createElement;

export default class EventsPage extends React.Component {
    render() {
        return e("div", { id: "page-events" }, [
            e("div", { className: "row" }, [
                e("div", { className: "col-lg-4"}, [
                    e(Hoverable),
                ]),
                e("div", { className: "col-lg-4"}, [
                    e(Clickable, {
                        id: "wait-class",
                        title: "Add class"
                    })
                ]),
                e("div", { className: "col-lg-4"}, [
                    e(Clickable, {
                        id: "wait-class-random",
                        title: "Add class 2",
                        randomTimeout: true
                    })
                ])
            ]),
            e("div", { className: "row" }, [
                e("div", { className: "col-lg-4"}, [
                    e(Clickable, {
                        id: "wait-no-class",
                        title: "Remove class",
                        show: true
                    })
                ]),
                e("div", { className: "col-lg-4"}, [
                    e(Clickable, {
                        id: "wait-no-class-random",
                        title: "Remove class 2",
                        show: true,
                        randomTimeout: true
                    })
                ]),
                e("div", { className: "col-lg-4"}, [
                    e(Appearable, {
                        id: "wait-element",
                        appear: true,
                        title: "Appearable"
                    })
                ]),
                e("div", { className: "col-lg-4"}, [
                    e(Appearable, {
                        id: "wait-no-element",
                        appear: false,
                        title: "Disappearable"
                    })
                ])
            ])
        ])
    }
}
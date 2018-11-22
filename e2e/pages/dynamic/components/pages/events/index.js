import Hoverable from "./hoverable.js";
import Clickable from "./clickable.js";

const e = React.createElement;

export default class EventsPage extends React.Component {
    render() {
        return e("div", { className: "row", id: "page-events" }, [
            e("div", { className: "col-lg-4"}, [
                e(Hoverable),
            ]),
            e("div", { className: "col-lg-4"}, [
                e(Clickable, { id: "wait-class" })
            ]),
            e("div", { className: "col-lg-4"}, [
                e(Clickable, { id: "wait-class-random", randomTimeout: true })
            ])
        ])
    }
}
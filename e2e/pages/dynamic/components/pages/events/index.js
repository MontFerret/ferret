import Hoverable from "./hoverable.js";

const e = React.createElement;

export default class EventsPage extends React.Component {
    render() {
        return e("div", { className: "row", id: "page-events" }, [
            e("div", { className: "col-lg-12"}, [
                e(Hoverable)
            ])
        ])
    }
}
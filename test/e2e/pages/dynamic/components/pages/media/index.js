const e = React.createElement;

export default class MediaPage extends React.Component {
    render() {
        return e("div", { id: "media" }, [
            e("div", { className: "row" }, [
                e("div", {className: "col-6"}, [
                    e("h3", null, [
                        "Height"
                    ]),
                    e("span", { id: "screen-height"}, [
                        window.innerHeight
                    ]),
                ]),
                e("div", {className: "col-6"}, [
                    e("h3", null, [
                        "Width"
                    ]),
                    e("span", { id: "screen-width"}, [
                        window.innerWidth
                    ])
                ])
            ])
        ])
    }
}
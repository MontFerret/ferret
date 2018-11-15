const e = React.createElement;

export default class HoverableComponent extends React.PureComponent {
    constructor(props) {
        super(props);

        this.state = {
            hovered: false
        };
    }

    handleMouseEnter() {
        this.setState({
            hovered: true
        });
    }

    handleMouseLeave() {
        this.setState({
            hovered: false
        });
    }

    render() {
        const children = [];
        children.push(
            e("p", null, [
                e("a", {
                    id: "hoverable-btn",
                    className: "btn btn-primary",
                    href: "#",
                    onMouseEnter: this.handleMouseEnter.bind(this),
                    onMouseLeave: this.handleMouseLeave.bind(this)
                }, [
                    "Hoverable link"
                ]),
            ])
        );

        if (this.state.hovered) {
            children.push(
                e("div", null, [
                    e("div", { id: "hoverable-content", className: "card card-body"}, [
                        "Lorem ipsum dolor sit amet."
                    ])
                ])
            )
        }

        return e("div", null, children);
    }
}
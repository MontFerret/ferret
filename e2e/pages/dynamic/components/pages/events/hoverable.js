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
        let content;

        if (this.state.hovered) {
            content =  e("p", { id: "hoverable-content"}, [
                "Lorem ipsum dolor sit amet."
            ]);
        }

        return e("div", { className: "card"}, [
            e("div", {className: "card-header"}, [
                e("button", {
                    id: "hoverable-btn",
                    className: "btn btn-primary",
                    onMouseEnter: this.handleMouseEnter.bind(this),
                    onMouseLeave: this.handleMouseLeave.bind(this)
                }, [
                    "Show content"
                ])
            ]),
            e("div", {className: "card-body"}, content)
        ]);
    }
}
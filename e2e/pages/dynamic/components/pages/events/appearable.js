import random from "../../../utils/random.js";

const e = React.createElement;

function render(id, props = {}) {
    return  e("span", { id: `${id}-content`, ...props }, ["Hello world"]);
}

export default class AppearableComponent extends React.PureComponent {
    constructor(props) {
        super(props);

        let element = null;

        if (props.appear) {
            if (props.useStyle) {
                element = render(props.id, { style: {display: "none"}})
            }
        } else {
            if (props.useStyle) {
                element = render(props.id, { style: {display: "block" }})
            } else {
                element = render(props.id)
            }
        }

        this.state = {
            element
        };
    }

    handleClick() {
        setTimeout(() => {
            const props = this.props;
            let element = null;

            if (props.appear) {
                if (props.useStyle) {
                    element = render(props.id, { style: {display: "block" }})
                } else {
                    element = render(props.id)
                }
            } else {
                if (props.useStyle) {
                    element = render(props.id, { style: {display: "none"}})
                }
            }

            this.setState({
                element,
            })
        }, random())
    }

    render() {
        const btnId = `${this.props.id}-btn`;

        return e("div", {className: "card"}, [
            e("div", { className: "card-header"}, [
                e("button", {
                    id: btnId,
                    className: "btn btn-primary",
                    onClick: this.handleClick.bind(this)
                }, [
                    this.props.title || "Toggle class"
                ])
            ]),
            e("div", { className: "card-body"}, this.state.element)
        ]);
    }
}
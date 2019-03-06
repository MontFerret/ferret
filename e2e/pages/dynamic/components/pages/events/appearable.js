import random from "../../../utils/random.js";

const e = React.createElement;

function render(id) {
    return  e("span", { id: `${id}-content` }, ["Hello world"]);
}

export default class AppearableComponent extends React.PureComponent {
    constructor(props) {
        super(props);

        this.state = {
            element: props.appear === true ? null : render(props.id)
        };
    }

    handleClick() {
        setTimeout(() => {
            this.setState({
                element: this.props.appear === true ? render(this.props.id) : null
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
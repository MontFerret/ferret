import random from "../../../utils/random.js";

const e = React.createElement;

export default class FocusableComponent extends React.PureComponent {
    constructor(props) {
        super(props);

        this.state = {
            show: props.show === true
        };
    }

    toggle(show) {
        let timeout = 500;

        if (this.props.randomTimeout) {
            timeout = random();
        }

        setTimeout(() => {
            this.setState({
                show: show
            })
        }, timeout)
    }

    handleFocus() {
        this.toggle(true);
    }

    handleBlur() {
        this.toggle(false);
    }

    render() {
        const btnId = `${this.props.id}-btn`;
        const contentId = `${this.props.id}-content`;
        const classNames = ["alert"];

        if (this.state.show === true) {
            classNames.push("alert-success");
        }

        return e("div", {className: "card focusable"}, [
            e("div", { className: "card-header"}, [
                e("div", { classNames: "form-group" }, [
                    e("input", {
                        id: btnId,
                        className: "form-control",
                        onFocus: this.handleFocus.bind(this),
                        onBlur: this.handleBlur.bind(this)

                    })
                ])
            ]),
            e("div", { className: "card-body"}, [
                e("div", { id: contentId, className: classNames.join(" ")}, [
                    e("p", null, [
                        "Lorem ipsum dolor sit amet."
                    ])
                ])
            ])
        ]);
    }
}
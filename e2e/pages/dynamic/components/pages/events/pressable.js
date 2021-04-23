import random from "../../../utils/random.js";

const e = React.createElement;

export default class PressableComponent extends React.PureComponent {
    constructor(props) {
        super(props);

        this.state = {
            key: ''
        };
    }

    handleKeyDown(e) {
        if (e.key == 'Unidentified') {
            return;
        }

        this.setState({
            key: this.state.key ? this.state.key + ' + ' + e.key : e.key
        })
    }

    handleReset() {
        this.setState({ key: '' })
    }

    render() {
        const inputId = `${this.props.id}-input`;
        const contentId = `${this.props.id}-content`;
        const classNames = ["alert", "alert-success"];

        return e("div", { id: this.props.id, className: "card clickable"}, [
            e("div", { className: "card-header"}, [
                e("div", { className: "form-group" }, [
                    e("label", null, "Pressable"),
                    e("input", {
                        id: inputId,
                        type: "text",
                        className: "form-control",
                        onKeyDown: this.handleKeyDown.bind(this)
                    }),
                    e("input", {
                        type: "button",
                        className: "btn btn-primary",
                        onClick: this.handleReset.bind(this),
                        value: "Reset"
                    },
                    )
                ]),
            ]),
            e("div", { className: "card-body"}, [
                e("div", { className: classNames.join(" ")}, [
                    e("p", {
                        id: contentId
                    }, [
                        this.state.key
                    ])
                ])
            ])
        ]);
    }
}
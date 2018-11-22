const e = React.createElement;

export default class ClickableComponent extends React.PureComponent {
    constructor(props) {
        super(props);

        this.state = {
            clicked: false
        };
    }

    handleClick() {
        let timeout = 500;

        if (this.props.randomTimeout) {
            timeout = Math.ceil(Math.random() * 1000 * 10);
        }

        setTimeout(() => {
            this.setState({
                clicked: !this.state.clicked
            })
        }, timeout)
    }

    render() {
        const btnId = `${this.props.id}-btn`;
        const contentId = `${this.props.id}-content`;
        const classNames = ["alert"];

        if (this.state.clicked) {
            classNames.push("alert-success");
        }

        return e("div", {className: "card clickable"}, [
            e("div", { className: "card-header"}, [
                e("button", {
                    id: btnId,
                    className: "btn btn-primary",
                    onClick: this.handleClick.bind(this)
                }, [
                    "Toggle class"
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
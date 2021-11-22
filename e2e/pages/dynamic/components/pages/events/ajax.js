import random from "../../../utils/random.js";

const e = React.createElement;

function request(url, body, method = 'GET') {
    fetch(url, {
        method,
        body
    })
        .then((res) => res.text())
        .then(text => console.log(text)).catch(er => console.error(er));
}

export default class AjaxComponent extends React.PureComponent {
    constructor(props) {
        super(props);

        this.state = {
            target: ''
        };
    }

    handleSeq(e) {
        [
            'https://www.montferret.dev/try/',
            'https://www.montferret.dev/docs/',
            'https://www.montferret.dev/blog/',
            'https://www.montferret.dev/cookbook/'
        ].forEach((url) => {
            setTimeout(() => {
                request(url)
            }, random(1000, 2000))
        });
    }

    handleTyping(evt) {
        this.setState({
            target: evt.target.value
        })
    }

    handleTarget(e) {
        setTimeout(() => {
            request(this.state.target)
        }, random())
    }

    render() {
        const inputId = `${this.props.id}-input`;
        const contentId = `${this.props.id}-content`;
        const classNames = ["alert", "alert-success"];

        return e("div", { id: this.props.id, className: "card ajax"}, [
            e("div", { className: "card-header"}, [
                "Ajax requests"
            ]),
            e("div", { className: "card-body"}, [
                e("div", { className: "form-group" }, [
                    e("label", null, "Make Sequential Request"),
                    e("input", {
                            id: inputId + "-seq-buttons",
                            type: "button",
                            className: "btn btn-primary",
                            onClick: this.handleSeq.bind(this),
                            value: "Send"
                        },
                    )
                ]),
                e("div", { className: "form-group" }, [
                    e("label", null, "Make Targeted Request"),
                    e("input", {
                            id: inputId,
                            type: "text",
                            onChange: this.handleTyping.bind(this),
                        },
                    ),
                    e("input", {
                            id: inputId + "-button",
                            type: "button",
                            className: "btn btn-primary",
                            onClick: this.handleTarget.bind(this),
                            value: "Send"
                        },
                    )
                ]),
            ])
        ]);
    }
}
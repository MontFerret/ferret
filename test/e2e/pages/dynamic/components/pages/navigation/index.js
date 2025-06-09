const e = React.createElement;

export default class NavigationPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            url: "",
        };

        this.handleTextInput = (evt) => {
            evt.preventDefault();

            this.setState({
                url: evt.target.value
            });
        };

        this.handleClick = () => {
            window.location.href = this.state.url;
        };
    }

    render() {
        return e("div", { id: "navigation" }, [
            e("div", { className: "form-group" }, [
                e("label", null, "Url"),
                e("input", {
                    id: "url",
                    type: "text",
                    className: "form-control",
                    onChange: this.handleTextInput
                }),
                e("input", {
                    id: "submit",
                    type: "button",
                    value: "Go",
                    className: "form-control",
                    onClick: this.handleClick
                }),
            ])
        ])
    }
}
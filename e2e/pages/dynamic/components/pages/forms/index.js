const e = React.createElement;

export default class FormsPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            textInput: "",
            select: "",
            multiSelect: "",
            textarea: ""
        };

        this.handleTextInput = (evt) => {
            evt.preventDefault();

            this.setState({
                textInput: evt.target.value
            });
        };

        this.handleSelect = (evt) => {
            evt.preventDefault();

            this.setState({
                select: evt.target.value
            });
        };

        this.handleMultiSelect = (evt) => {
            evt.preventDefault();

            this.setState({
                multiSelect: Array.prototype.map.call(evt.target.selectedOptions, i => i.value).join(", ")
            });
        };

        this.handleTextarea = (evt) => {
            evt.preventDefault();

            this.setState({
                textarea: evt.target.value
            });
        }
    }

    render() {
        return e("form", { id: "page-form" }, [
            e("div", { className: "form-group" }, [
                e("label", null, "Text input"),
                e("input", {
                    id: "text_input",
                    type: "text",
                    className: "form-control",
                    onChange: this.handleTextInput
                }),
                e("small", {
                    id: "text_output",
                    className: "form-text text-muted"
                },
                    this.state.textInput
                )
            ]),
            e("div", { className: "form-group" }, [
                e("label", null, "Select"),
                e("select", {
                    id: "select_input",
                    className: "form-control",
                    onChange: this.handleSelect
                    },
                    [
                        e("option", null, 1),
                        e("option", null, 2),
                        e("option", null, 3),
                        e("option", null, 4),
                        e("option", null, 5),
                    ]
                ),
                e("small", {
                        id: "select_output",
                        className: "form-text text-muted"
                    }, this.state.select
                )
            ]),
            e("div", { className: "form-group" }, [
                e("label", null, "Multi select"),
                e("select", {
                        id: "multi_select_input",
                        multiple: true,
                        className: "form-control",
                        onChange: this.handleMultiSelect
                    },
                    [
                        e("option", null, 1),
                        e("option", null, 2),
                        e("option", null, 3),
                        e("option", null, 4),
                        e("option", null, 5),
                    ]
                ),
                e("small", {
                        id: "multi_select_output",
                        className: "form-text text-muted"
                    }, this.state.multiSelect
                )
            ]),
            e("div", { className: "form-group" }, [
                e("label", null, "Textarea"),
                e("textarea", {
                        id: "textarea_input",
                        rows:"5",
                        className: "form-control",
                        onChange: this.handleTextarea
                    }
                ),
                e("small", {
                        id: "textarea_output",
                        className: "form-text text-muted"
                    }, this.state.textarea
                )
            ]),
        ])
    }
}
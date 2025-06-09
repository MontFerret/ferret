import { parse } from '../../../utils/qs.js';

const e = React.createElement;

export default class IFramePage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            url: '/'
        }
    }

    handleUrlInput(evt) {
        this.setState({
            url: evt.target.value
        })
    }

    handleReload() {
        window.location.href = this.state.url;
    }

    render() {
        const search = parse(this.props.location.search);

        let redirect;

        if (search.src) {
            redirect = search.src;
        }

        let navGroup;

        if (window.top !== window) {
            navGroup = [
                e("div", { className: "form-group row" }, [
                    e("input", {
                        id: "url_input",
                        type: "text",
                        className: "form-control",
                        onChange: this.handleUrlInput.bind(this)
                    }),
                    e("button", {
                        id: "submit",
                        className: "btn btn-primary",
                        onClick: this.handleReload.bind(this)
                    }, [
                        "Navigate"
                    ]),
                ]),
            ];
        }

        return e("div", { id: "iframe" }, [
            navGroup,
            e("iframe", {
                name: 'nested',
                style: {
                    width: '100%',
                    height: '800px',
                },
                src: redirect ? `/?redirect=${redirect}` : '/'
            }),
        ])
    }
}
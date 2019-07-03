import { parse } from '../../../utils/qs.js';

const e = React.createElement;

export default class IFramePage extends React.Component {
    render() {
        const search = parse(this.props.location.search);

        let redirect;

        if (search.src) {
            redirect = search.src;
        }

        return e("div", { id: "iframe" }, [
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
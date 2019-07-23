const e = React.createElement;
const Pagination = ReactBootstrap.Pagination;
const Item = Pagination.Item;

export default class PaginationPage extends React.Component {
    min = 1;
    max = 5;

    constructor(props) {
        super(props);

        this.state = {
            active: 1
        }
    }

    handleClick(number) {
        if (number > this.max || number < this.min) {
            return;
        }

        this.setState({
            active: number
        })
    }

    render() {
        const items = [];

        items.push(e(Pagination.Prev, {
            className: "page-item-prev",
            disabled: this.state.active === this.min,
            onClick: this.handleClick.bind(this, this.state.active - 1)
        }));

        for (let number = this.min; number <= this.max; number++) {
            items.push(
                e(Item, {
                    key: number,
                    active: number === this.state.active,
                    onClick: this.handleClick.bind(this, number)
                }, number)
            );
        }

        items.push(e(Pagination.Next, {
            className: "page-item-next",
            disabled: this.state.active === this.max,
            onClick: this.handleClick.bind(this, this.state.active + 1)
        }));

        return e("div", { className: "row"}, [
            e("div", { className: "col-12" }, [
                e(Pagination, {
                }, items)
            ])
        ])
    }
}
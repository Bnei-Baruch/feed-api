import React, { Component } from 'react'
import Feed from './Feed.js'

const MORE_ITEMS = 10;

class FeedContainer extends Component {
    constructor() {
        super();
        this.state = {
            items: [],
        };
		this.moreHandler = this.moreHandler.bind(this);
    }

    componentDidMount() {
		this.moreHandler();
    }

	moreHandler() {
		console.log('More');
		const offset = this.state.items.length;
        fetch(`http://bbdev6.kbb1.com:9590/items?offset=${offset}`).
        then(results => results.json()).
        then(data => {
            this.setState({items: this.state.items.concat(data.items)});
        });
	}

    render() {
        const {items} = this.state;
        return (
            <Feed items={items} more={this.moreHandler} />
        );
    }
}

export default FeedContainer

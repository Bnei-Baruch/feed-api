import React, { PureComponent } from 'react';

import PropTypes from 'prop-types';
import { Button, Container, Feed as SemanticFeed, Icon } from 'semantic-ui-react'

import './Feed.css';

class Feed extends PureComponent {
	static propTypes = {
		items: PropTypes.arrayOf(PropTypes.shape({})),
		more: PropTypes.func,
	};

	static defaultProps = {
		items: [],
	};

	renderItem(item) {
		console.log(item);
		return (
			<SemanticFeed.Event key={item.id}>
  			  <SemanticFeed.Label>
  				{item.created_at} - {item.name}
  			  </SemanticFeed.Label>
			</SemanticFeed.Event>
		);
	}

	render() {
		const {items, more} = this.props;
		return (
			<Container>
				<SemanticFeed>
					{items.map((item) => this.renderItem(item))}
				</SemanticFeed>
				<Button onClick={more}>More</Button>
			</Container>
		);
	}
}

export default Feed;


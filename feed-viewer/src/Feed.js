import React, { PureComponent } from 'react';

import PropTypes from 'prop-types';
import { Button, Container, Icon, Segment } from 'semantic-ui-react'

import './Feed.css';
import { CT_DAILY_LESSON } from './helpers/consts';
import { canonicalLink } from './helpers/links';

class Feed extends PureComponent {
	static propTypes = {
		items: PropTypes.arrayOf(PropTypes.shape({})),
		more: PropTypes.func,
	};

	static defaultProps = {
		items: [],
	};

	renderCU = cu => {
		return (
		  <a
			key={cu.id}
			href={canonicalLink(cu, "he")}
			target="_blank"
		  >
			<Button basic size="tiny" className="link_to_cu">
			  {cu.name}
			</Button>
		  </a>
		);
	};

	renderItem(item) {
		console.log(item);
		const {id: mdbUid} = item;
		const toLink = canonicalLink(item || { id: mdbUid, content_type: item.content_type }, "he");
		if (item.content_type === CT_DAILY_LESSON) {
			return (
				<Segment key={item.id}>
				  <Container>
					<Container as="h3">
					  <a
						className="search__link"
						href={toLink}
						target="_blank"
					  >
						<span>שיעור בוקר {item.film_date}</span>
					  </a>
					</Container>

					{
					false &&
					<Container className="content">
					  <span>{item.content_units.length}{' '}פרקים</span>
					</Container>
					}
					<div className="clear" />
				  </Container>

				  <Container className="content clear margin-top-8">
					{item.content_units.slice(0, 5).map(this.renderCU)}

					<a
					  href={toLink}
					  target="_blank"
					>
					  <Icon name="tasks" size="small" />
					  {`הראה את כל ${item.content_units.length} החלקים`}
					</a>
				  </Container>
			    </Segment>
			);
		} else {
			return (
				<Segment key={item.id}>
				  <Container>
					{item.name} - {item.id} - {item.content_type}
				  </Container>
				</Segment>
			);
		}
	}

	render() {
		const {items, more} = this.props;
		return (
			<Container>
				{items.map((item) => this.renderItem(item))}
				<Button onClick={more}>More</Button>
			</Container>
		);
	}
}

export default Feed;


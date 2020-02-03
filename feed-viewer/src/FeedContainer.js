import React, { Component } from 'react'
import Feed from './Feed.js'

const MORE_ITEMS = 10;

class FeedContainer extends Component {
    constructor() {
        super();
        this.state = {
            feed: [],
			items: [],
			itemsByUid: {},
        };
		this.moreHandler = this.moreHandler.bind(this);
    }

    componentDidMount() {
		this.moreHandler();
    }

	paramsToUrl(params) {
		return Object.entries(params).map(([param, value]) => `${param}=${value}`).join('&')
	}

	fetchContentUnits(contentItems) {
		const contentUnits = contentItems.filter((contentItem) => contentItem.content_type === 'DAILY_LESSON');
		if (contentUnits.length === 0) {
			return Promise.resolve({content_units: [], total: 0});
		}
		
		const params = {
			'page_size': contentUnits.length,
			'language': 'he',
			'with_files': 'true',
			'with_derivations': 'true',
		};
		const ids = contentUnits.map((contentUnit) => `id=${contentUnit.uid}`).join('&')
		return fetch(`https://kabbalahmedia.info/backend/content_units?${ids}&${this.paramsToUrl(params)}`).
			then((results) => results.json());
	}

	fetchCollections(contentItems) {
		const collections = contentItems.filter((contentItem) => contentItem.content_type === 'DAILY_LESSON');
		if (collections.length === 0) {
			return Promise.resolve({collections: [], total: 0});
		}
		
		const params = {
			'page_size': collections.length,
			'language': 'he',
		};
		const ids = collections.map((collection) => `id=${collection.uid}`).join('&')
		return fetch(`https://kabbalahmedia.info/backend/collections?${ids}&${this.paramsToUrl(params)}`).
			then((results) => results.json());
	}

	moreHandler() {
		const {feed, itemsByUid} = this.state;
		console.log('More');
        fetch(`http://bbdev6.kbb1.com:9590/more`, {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({more_items: 10, current_feed: feed}),
		}).
		then(results => results.json()).
		then(data => {
			const newFeed = data.feed;
			const newFeedUids = new Set(feed.map((contentItem) => contentItem.uid));
			const fetchItems = newFeed.filter((contentItem) => !(contentItem.uid in itemsByUid));
			this.fetchCollections(fetchItems).then((data) => {
				data.collections.forEach((collection) => {
					itemsByUid[collection.id] = collection;
				});
				feed.forEach((contentItem) => {
					if (!newFeedUids.has(contentItem)) {
						delete itemsByUid[contentItem.uid];
					}
				});
				const items = newFeed.map((contentItem) => itemsByUid[contentItem.uid]);
				this.setState({feed: newFeed, items, itemsByUid});
			})
		});
	}

    render() {
        const {items} = this.state;
        return (
            <Feed items={items} more={this.moreHandler} />
        );
    }
}

export default FeedContainer;

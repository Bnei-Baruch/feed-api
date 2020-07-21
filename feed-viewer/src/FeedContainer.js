import React, { Component } from 'react'
import merge from 'lodash/merge';
import isEqual from 'lodash/isEqual'
import Feed from './Feed.js'
import {
	CT_VIDEO_PROGRAM,
	CT_CLIPS,
	CT_ARTICLES,
	CONTENT_UNIT_TYPES_SET,
	COLLECTION_TYPES_SET
} from './helpers/consts'

const MORE_ITEMS = 10;

class FeedContainer extends Component {
    constructor() {
        super();
        this.state = {
            feed: [],
			items: [],
			itemsByUid: {},
			fetchingSubscribeCollections: false,
			subscribeCollections: [],
			options: {},
        };
		this.moreHandler = this.moreHandler.bind(this);
		this.resetHandler = this.resetHandler.bind(this);
		this.updateOptions = this.updateOptions.bind(this);
    }

    componentDidMount() {
		this.moreHandler();
		this.fetchSubscribeCollections();
    }

    shouldComponentUpdate(nextProps, nextState) {
        const stateEqual =  isEqual(this.state, nextState);
        //console.log('state equal:', stateEqual, this.state.options, nextState.options);
        return !stateEqual;
    }

	paramsToUrl(params) {
		return Object.entries(params).reduce((pairs, [param, value]) => {
			if (Array.isArray(value)) {
				value.forEach((v) => pairs.push(`${param}=${v}`));
			} else {
				pairs.push(`${param}=${value}`);
			}
			return pairs;
		}, []).join('&')
	}

	fetchSubscribeCollections() {
		this.setState({...this.state, fetchingSubscribeCollections: true})

		const params = {
			'content_type': [CT_VIDEO_PROGRAM, CT_CLIPS, CT_ARTICLES],
			'language': 'he',
			'page_size': 1000000, // TODO: Fix to fetch correct size. 
		};
		fetch(`https://kabbalahmedia.info/backend/collections?${this.paramsToUrl(params)}`)
			.then((results) => results.json())
			.then((json) => {
			//console.log(json);
			this.setState({...this.state, fetchingSubscribeCollections: false, subscribeCollections: json.collections});
		});
	}

	fetchContentUnits(contentItems) {
		const contentUnits = contentItems.filter((contentItem) => CONTENT_UNIT_TYPES_SET.has(contentItem.content_type));
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
		return fetch(`https://kabbalahmedia.info/backend/content_units?${ids}&${this.paramsToUrl(params)}`)
			.then((results) => results.json());
	}

	fetchCollections(contentItems) {
		const collections = contentItems.filter((contentItem) => COLLECTION_TYPES_SET.has(contentItem.content_type));
		if (collections.length === 0) {
			return Promise.resolve({collections: [], total: 0});
		}
		
		const params = {
			'page_size': collections.length,
			'language': 'he',
		};
		const ids = collections.map((collection) => `id=${collection.uid}`).join('&')
		return fetch(`https://kabbalahmedia.info/backend/collections?${ids}&${this.paramsToUrl(params)}`)
			.then((results) => results.json());
	}

	updateOptions(updateOptions) {
        const {options} = this.state;
		this.setState({...this.state, options: merge({}, options, updateOptions)});
	}

	resetHandler() {
		this.setState({...this.state, feed: [], items: [], itemsByUid: {}}, () => this.moreHandler());
	}

	moreHandler() {
		const {feed, itemsByUid, options} = this.state;
		console.log('More');
        fetch(`http://bbdev6.kbb1.com:9590/more`, {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({more_items: MORE_ITEMS, current_feed: feed, options}),
		}).then(results => results.json()).then(data => {
			const newFeed = data.feed;
			const newFeedUids = new Set(newFeed.map((contentItem) => contentItem.uid));
			const fetchItems = newFeed.filter((contentItem) => !(contentItem.uid in itemsByUid));
			const fetchPromises = [
				// Fetch collections.
				this.fetchCollections(fetchItems).then((data) => {
					data.collections.forEach((collection) => {
						if (collection.content_units) {
							collection.content_units.sort((a, b) => {
								return b.film_date.localeCompare(a.film_date);
							});
						}
						itemsByUid[collection.id] = collection;
					});
				}),
				// Fetch content units.
				this.fetchContentUnits(fetchItems).then((data) => {
					data.content_units.forEach((content_unit) => {
						itemsByUid[content_unit.id] = content_unit;
					});
				}),
			];
			Promise.all(fetchPromises).then(() => {
				// Delete data from non required uids.
				feed.forEach((contentItem) => {
					if (!newFeedUids.has(contentItem.uid)) {
						delete itemsByUid[contentItem.uid];
					}
				});
				const items = newFeed.map((contentItem) => itemsByUid[contentItem.uid]);
				this.setState({feed: newFeed, items, itemsByUid});
			});
		});
	}

    render() {
        //console.log('render container');
        const {items, options, fetchingSubscribeCollections, subscribeCollections} = this.state;
        return (
            <Feed
				items={items}
                options={options}
				more={this.moreHandler}
				reset={this.resetHandler}
				updateOptions={this.updateOptions}
				fetchingSubscribeCollections={fetchingSubscribeCollections}
				subscribeCollections={subscribeCollections}
			/>
        );
    }
}

export default FeedContainer;

import React, { Component } from 'react'
import merge from 'lodash/merge';
import isEqual from 'lodash/isEqual'
import Feed from './Feed.js'
import {
	CT_VIDEO_PROGRAM,
	CT_CLIPS,
	CT_ARTICLES,
} from './helpers/consts'
import {
  more,
  paramsToUrl,
} from './api'

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

	fetchSubscribeCollections() {
		this.setState({...this.state, fetchingSubscribeCollections: true})

		const params = {
			'content_type': [CT_VIDEO_PROGRAM, CT_CLIPS, CT_ARTICLES],
			'language': 'he',
			'page_size': 1000000, // TODO: Fix to fetch correct size. 
		};
		fetch(`https://kabbalahmedia.info/backend/collections?${paramsToUrl(params)}`)
			.then((results) => results.json())
			.then((json) => {
			//console.log(json);
			this.setState({...this.state, fetchingSubscribeCollections: false, subscribeCollections: json.collections});
		});
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
    more(feed, itemsByUid, options).then(({feed, items, itemsByUid}) => this.setState({feed, items, itemsByUid}));
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

import React, { Component } from 'react'
import merge from 'lodash/merge';
import isEqual from 'lodash/isEqual'
import Feed from './Feed.js'
/*import {
	CT_VIDEO_PROGRAM,
	CT_CLIPS,
	CT_ARTICLES,
} from './helpers/consts'*/
import {
  more,
  feed as feedFunc,
//  paramsToUrl,
} from './api'

class FeedContainer extends Component {
  constructor() {
    super();
    this.state = {
      feed: [],
      items: [],
      itemsByUid: {},
      /*fetchingSubscribeCollections: false,
      subscribeCollections: [],*/
      options: {},
      error: '',
      spec: new URLSearchParams(window.location.search).get('spec') || '',
      debugTimestamp: new URLSearchParams(window.location.search).get('debug_timestamp') || '',
      numItems: new URLSearchParams(window.location.search).get('num_items') || '20',
      languages: new URLSearchParams(window.location.search).get('languages') || 'he,en',
      skipUids: new URLSearchParams(window.location.search).get('skip_uids') || '',
      feedOrMore: new URLSearchParams(window.location.search).get('feed_or_more') || 'feed',
      withBlogPosts: !!(new URLSearchParams(window.location.search).get('with_blog_posts')),
    };
    this.moreHandler = this.moreHandler.bind(this);
    this.resetHandler = this.resetHandler.bind(this);
    this.updateOptions = this.updateOptions.bind(this);
    this.updateSpec = this.updateSpec.bind(this);
    this.updateDebugTimestamp = this.updateDebugTimestamp.bind(this);
    this.updateNumItems = this.updateNumItems.bind(this);
    this.updateLanguages = this.updateLanguages.bind(this);
    this.updateSkipUids = this.updateSkipUids.bind(this);
    this.updateFeedOrMore = this.updateFeedOrMore.bind(this);
    this.updateWithBlogPosts = this.updateWithBlogPosts.bind(this);
  }

  componentDidMount() {
    this.moreHandler();
    // this.fetchSubscribeCollections();
  }

  shouldComponentUpdate(nextProps, nextState) {
    const stateEqual =  isEqual(this.state, nextState);
    //console.log('state equal:', stateEqual, this.state.options, nextState.options);
    return !stateEqual;
  }

  /*
	fetchSubscribeCollections() {
		this.setState({...this.state, fetchingSubscribeCollections: true})

		const params = {
			'content_type': [CT_VIDEO_PROGRAM, CT_CLIPS, CT_ARTICLES],
			'language': 'he',
			'page_size': 1000000,  // TODO: Fix to fetch correct size. 
		};
		fetch(`https://kabbalahmedia.info/backend/collections?${paramsToUrl(params)}`)
			.then((results) => results.json())
			.then((json) => {
			//console.log(json);
			this.setState({...this.state, fetchingSubscribeCollections: false, subscribeCollections: json.collections});
		});
	}
  */

	updateOptions(updateOptions) {
    const {options} = this.state;
		this.setState({...this.state, options: merge({}, options, updateOptions)});
	}

  updateSpec(spec) {
    this.setState({spec}, () => this.updateUrl());
  }

  updateDebugTimestamp(debugTimestamp) {
    this.setState({debugTimestamp}, () => this.updateUrl());
  }

  updateNumItems(numItems) {
    this.setState({numItems}, () => this.updateUrl());
  }

  updateLanguages(languages) {
    this.setState({languages}, () => this.updateUrl());
  }

  updateSkipUids(skipUids) {
    this.setState({skipUids}, () => this.updateUrl());
  }

  updateFeedOrMore(feedOrMore) {
    console.log(feedOrMore);
    this.setState({feedOrMore}, () => this.updateUrl());
  }

  updateWithBlogPosts(withBlogPosts) {
    console.log(withBlogPosts);
    this.setState({withBlogPosts}, () => this.updateUrl());
  }

  updateUrl() {
    const {
      debugTimestamp,
      feedOrMore,
      languages,
      numItems,
      skipUids,
      spec,
      withBlogPosts,
    } = this.state;
    const params = [];
    if (spec) {
      params.push(`spec=${spec}`);
    }
    if (debugTimestamp) {
      params.push(`debug_timestamp=${debugTimestamp}`);
    }
    if (numItems) {
      params.push(`num_items=${numItems}`);
    }
    if (languages) {
      params.push(`languages=${languages}`);
    }
    if (skipUids) {
      params.push(`skip_uids=${skipUids}`);
    }
    if (feedOrMore) {
      params.push(`feed_or_more=${feedOrMore}`);
    }
    if (withBlogPosts) {
      params.push('with_blog_posts');
    }
    const url = new URL(window.location.toString());
    url.search = `?${params.join('&')}`;
    window.history.replaceState({}, 'Feed', url.toString());
  }

	resetHandler() {
		this.setState({...this.state, feed: [], items: [], itemsByUid: {}}, () => this.moreHandler());
	}

	moreHandler() {
    this.setState({error: ''}, () => {
      const {
        debugTimestamp,
        feed,
        feedOrMore,
        itemsByUid,
        languages,
        numItems,
        options,
        skipUids,
        spec,
        withBlogPosts,
      } = this.state;

      const parseSpec = (spec) => {
        if (spec) {
          try {
            return [JSON.parse(spec), null];
          } catch(e) {
            return [null, e];
          }
        }
        return [null, null];
      }
      const [specObj/*, specParseErr*/] = parseSpec(spec);
      if (specObj) {
        options.spec = specObj;
      }
      if (debugTimestamp) {
        options.debug_timestamp = Number(debugTimestamp);
      } else {
        delete options.debug_timestamp;
      }
      if (languages) {
        options.languages = languages.split(',');
      } else {
        delete options.languages;
      }
      if (skipUids) {
        options.skip_uids = skipUids.split(',');
      } else {
        delete options.skip_uids;
      }
      options.with_posts = withBlogPosts;

      const feedOrMoreFunc = feedOrMore === 'feed' ? feedFunc : more;
      
      feedOrMoreFunc(feed, itemsByUid, options, Number(numItems))
        .then(({feed, items, itemsByUid}) => this.setState({feed, items, itemsByUid}))
        .catch((error) => this.setState({error : String(error)}));
    });
	}

  render() {
    //console.log('render container');
    const {
      /*fetchingSubscribeCollections, subscribeCollections,*/ 
      debugTimestamp, 
      error, 
      feedOrMore,
      items, 
      languages, 
      numItems, 
      options, 
      skipUids,
      spec, 
      withBlogPosts,
    } = this.state;
    return (
      <Feed
        items={items}
        options={options}
        more={this.moreHandler}
        reset={this.resetHandler}
        updateOptions={this.updateOptions}
        /*fetchingSubscribeCollections={fetchingSubscribeCollections}
        subscribeCollections={subscribeCollections}*/
        error={error}
        spec={spec}
        withBlogPosts={withBlogPosts}
        updateSpec={this.updateSpec}
        debugTimestamp={debugTimestamp}
        updateDebugTimestamp={this.updateDebugTimestamp}
        numItems={numItems}
        updateNumItems={this.updateNumItems}
        languages={languages}
        updateLanguages={this.updateLanguages}
        skipUids={skipUids}
        updateSkipUids={this.updateSkipUids}
        feedOrMore={feedOrMore}
        updateFeedOrMore={this.updateFeedOrMore}
        updateWithBlogPosts={this.updateWithBlogPosts}
      />
    );
  }
}

export default FeedContainer;

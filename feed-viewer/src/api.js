import {
	CONTENT_UNIT_TYPES_SET,
	COLLECTION_TYPES_SET,
  CT_BLOG_POST,
} from './helpers/consts'

const MORE_ITEMS = 10;
const RECOMMEND_ITEMS = 6;

export const paramsToUrl = (params) => {
  return Object.entries(params).reduce((pairs, [param, value]) => {
    if (Array.isArray(value)) {
      value.forEach((v) => pairs.push(`${param}=${v}`));
    } else {
      pairs.push(`${param}=${value}`);
    }
    return pairs;
  }, []).join('&')
};

export const fetchBlogPosts = (contentItems) => {
  const blogPosts = contentItems.filter((contentItem) => contentItem.content_type === CT_BLOG_POST);
  console.log('fetchBlogPosts', contentItems, blogPosts);
  if (blogPosts.length === 0) {
    return Promise.resolve({posts: [], total: 0});
  }
  const params = {
    'page_size': blogPosts.length,
  };
  const ids = blogPosts.map((contentUnit) => `id=${contentUnit.uid}`).join('&')
  return fetch(`https://kabbalahmedia.info/backend/posts?${ids}&${paramsToUrl(params)}`)
    .then((results) => results.json()).then((json) => {
      if (json.posts.length !== blogPosts.length) {
        return Promise.reject(`Expected number of posts ${blogPosts.length} got ${json.posts.length}`);
      };
      blogPosts.forEach((blogPost, index) => {
        json.posts[index].uid = blogPost.uid;
        json.posts[index].id = blogPost.uid;
        json.posts[index].content_type = 'POST';
      });
      return json;
    });
}

export const fetchContentUnits = (contentItems) => {
  const contentUnits = contentItems.filter((contentItem) => CONTENT_UNIT_TYPES_SET.has(contentItem.content_type));
  console.log('fetchContentUnits', contentItems, contentUnits);
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
  return fetch(`https://kabbalahmedia.info/backend/content_units?${ids}&${paramsToUrl(params)}`)
    .then((results) => results.json());
};

export const fetchCollections = (contentItems) => {
  const collections = contentItems.filter((contentItem) => COLLECTION_TYPES_SET.has(contentItem.content_type));
  if (collections.length === 0) {
    return Promise.resolve({collections: [], total: 0});
  }
  
  const params = {
    'page_size': collections.length,
    'language': 'he',
  };
  const ids = collections.map((collection) => `id=${collection.uid}`).join('&')
  return fetch(`https://kabbalahmedia.info/backend/collections?${ids}&${paramsToUrl(params)}`)
    .then((results) => results.json());
};

export const feed = (feed, itemsByUid, options, num_items = MORE_ITEMS) => {
  return moreOrReccomend(feed, itemsByUid, options, 'feed', num_items);
};

export const more = (feed, itemsByUid, options, num_items = MORE_ITEMS) => {
  return moreOrReccomend(feed, itemsByUid, options, 'more', num_items);
};

export const recommend = (feed, itemsByUid, options, num_items = RECOMMEND_ITEMS) => {
  return moreOrReccomend(feed, itemsByUid, options, 'recommend', num_items);
};

export const moreOrReccomend = (feed, itemsByUid, options, handler, numItems) => {
  console.log('moreOrReccomend', JSON.stringify({more_items: numItems, current_feed: feed, options}));
  return fetch(`http://bbdev6.kbb1.com:9590/${handler}`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({more_items: numItems, current_feed: feed, options}),
  }).then(results => {
    console.log(results);
    if (results.status !== 200) {
      return results.text().then((text) => {
        return Promise.reject(`${results.status}: ${results.statusText}. Error from server: ${text}`);
      });
    }
    return results.json();
  }).then(data => {
    console.log('moreOrReccomend', data, feed, itemsByUid, options, handler, numItems);
    const newFeed = (feed || []).concat(data.feed || []);
    const newFeedUids = new Set(newFeed.map((contentItem) => contentItem.uid));
    const fetchItems = newFeed.filter((contentItem) => !(contentItem.uid in itemsByUid));
    console.log('fetchItems', fetchItems);
    const fetchPromises = [
      // Fetch collections.
      fetchCollections(fetchItems).then((data) => {
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
      fetchContentUnits(fetchItems).then((data) => {
        console.log('data', data);
        data.content_units.forEach((content_unit) => {
          itemsByUid[content_unit.id] = content_unit;
        });
      }),
      // Fetch blog posts.
      fetchBlogPosts(fetchItems).then((data) => {
        console.log('blogs posts data', data);
        data.posts.forEach((post) => {
          itemsByUid[post.uid] = post;
        });
      }),
    ];
    return Promise.all(fetchPromises).then(() => {
      // Delete data from non required uids.
      feed.forEach((contentItem) => {
        if (!newFeedUids.has(contentItem.uid)) {
          delete itemsByUid[contentItem.uid];
        }
      });
      const items = newFeed.map((contentItem) => itemsByUid[contentItem.uid]);
      const ret = {feed: newFeed, items, itemsByUid};
      console.log('ret', ret);
      return ret;
    });
  });
};

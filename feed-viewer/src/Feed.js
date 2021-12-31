import React, { PureComponent } from 'react';

import PropTypes from 'prop-types';
import {
  Button,
  // Checkbox,
  Container,
  // Dimmer,
  Grid,
  Icon,
  Input,
  // Loader,
  Radio,
  Segment,
  TextArea,
} from 'semantic-ui-react'

import SpecTree from './SpecTree.js'

import './Feed.css';
import { CT_DAILY_LESSON, CT_VIDEO_PROGRAM } from './helpers/consts';
import { canonicalLink } from './helpers/links';

class Feed extends PureComponent {
  /*constructor() {
    super();
    /*this.optionsChange = this.optionsChange.bind(this);* /
  }*/

	static propTypes = {
		items: PropTypes.arrayOf(PropTypes.shape({})),
        options: PropTypes.shape({}),
		more: PropTypes.func,
		reset: PropTypes.func,
		updateOptions: PropTypes.func,
		/*fetchingSubscribeCollections: PropTypes.bool,
		subscribeCollections: PropTypes.arrayOf(PropTypes.shape({})),*/
    error: PropTypes.string,
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
            rel="noopener noreferrer"
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
		if ([CT_DAILY_LESSON, CT_VIDEO_PROGRAM].includes(item.content_type)) {
			return (
				<Segment key={item.id}>
				  <Container>
					<Container as="h3">
					  <a
						className="search__link"
						href={toLink}
						target="_blank"
                        rel="noopener noreferrer"
					  >
						{ item.content_type === CT_DAILY_LESSON ?
							<span>שיעור בוקר {item.film_date}</span> :
							<span>תוכנית {item.name}, פרק אחרון ב-{item.content_units[0].film_date}</span> }
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
                      rel="noopener noreferrer"
					  style={{'whiteSpace': 'nowrap'}}
					>
					  <Icon name="tasks" size="small" />
					  &nbsp;
					  {`הראה את כל ${item.content_units.length} החלקים`}
					</a>
				  </Container>
			    </Segment>
			);
		} else {
			return (
				<Segment key={item.id}>
				  <Container>
					<h3>{item.film_date}</h3>
					- {item.name} - {item.id} - {item.content_type}
				  </Container>
				</Segment>
			);
		}
	}

  /*
  nextCollectionState(contentType, mid, options) {
      return this.nextOptionState(
          (options.collections &&
          contentType in options.collections &&
          mid in options.collections[contentType] &&
          options.collections[contentType][mid]) || '');
  }

  nextContentTypeState(contentType, options) {
    return this.nextOptionState(
        (options.content_types &&
        contentType in options.content_types &&
          options.content_types[contentType]) || '');
  }

  nextOptionState(option) {
    switch (option) {
      case 'subscribe':
        return 'unsubscribe';
      case 'unsubscribe':
        return 'default';
      case 'default':
        return 'subscribe';
      default:
        return 'subscribe';
    }
  }

	optionsChange(e, elem) {
    //console.log(elem);
    const {label} = elem;
		const {subscribeCollections, updateOptions, options} = this.props;
    const c = subscribeCollections.find((c) => c.name === label);
    const update = c ? {
        collections: {
            [c.content_type]: {
                [c.id]: this.nextCollectionState(c.content_type, c.id, options),
            },
        },
    } : {
        content_types: {
            [label]: this.nextContentTypeState(label, options),
        },
    };
    updateOptions(update);
	}

  subscriptionChecked(value) {
      switch (value) {
          case 'subscribe':
              return true;
          case 'unsubscribe':
              return undefined;
          case 'default':
              return false;
          default:
              return undefined;
      }
  }

  subscriptionIndeterminate(value) {
      switch (value) {
          case 'subscribe':
              return false;
          case 'unsubscribe':
              return true;
          case 'default':
              return false;
          default:
              return false;
      }
  }

  contentTypeChecked(contentType, options) {
      if (!('content_types' in options &&
          contentType in options.content_types)) {
          return undefined;
      }
      return this.subscriptionChecked(options.content_types[contentType]);
  }

  contentTypeIndeterminate(contentType, options) {
      if (!('content_types' in options &&
          contentType in options.content_types)) {
          return false;
      }
      return this.subscriptionIndeterminate(options.content_types[contentType]);
  }

  collectionChecked(contentType, collectionMid, options) {
      if (!('collections' in options &&
          contentType in options.collections &&
          collectionMid in options.collections[contentType])) {
          return false;
      }
      return this.subscriptionChecked(options.collections[contentType][collectionMid]);
  }

  collectionIndeterminate(contentType, collectionMid, options) {
      if (!('collections' in options &&
          contentType in options.collections &&
          collectionMid in options.collections[contentType])) {
          return false;
      }
      return this.subscriptionIndeterminate(options.collections[contentType][collectionMid]);
  }
  */

	render() {
		const {
      debugTimestamp,
      error,
      feedOrMore,
      // fetchingSubscribeCollections,
      items,
      languages,
      more,
      numItems,
      // options,
      reset,
      skipUids,
      spec,
      // subscribeCollections,
      updateDebugTimestamp,
      updateFeedOrMore,
      updateLanguages,
      updateNumItems,
      updateSkipUids,
      updateSpec,
    } = this.props;

    console.log('items', items);
		//console.log(fetchingSubscribeCollections, subscribeCollections);
        //console.log('checked:', this.contentTypeChecked('CT_DAILY_LESSON', options),
        //    'indetermidiate:', this.contentTypeIndeterminate('CT_DAILY_LESSON', options));
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
    const [specObj, specParseErr] = parseSpec(spec);

		return (
			<Grid columns={2}>
				<Grid.Row>
					<Grid.Column>
						<Segment style={{'direction': 'ltr'}}>
              <h3>Options</h3>
              <Segment textAlign='left'>
                <table>
                  <tbody>
                    <tr><td>Feed</td><td><Radio toggle defaultChecked={feedOrMore !== 'feed'} onChange={(event, data) => { console.log(event, data); updateFeedOrMore(data.checked ? 'more' : 'feed')}} /></td><td>More</td></tr>
                    <tr><td>Debug timestamp:</td><td colSpan="2"><Input placeholder='Debug timestamp...' defaultValue={debugTimestamp} onChange={(event, data) => updateDebugTimestamp(data.value)} /></td></tr>
                    <tr><td>Num Items:</td><td colSpan="2"><Input placeholder='Num Items to Recommend' defaultValue={numItems} onChange={(event, data) => updateNumItems(data.value)} /></td></tr>
                    <tr><td>Languages:</td><td colSpan="2"><Input placeholder='List of preffered languages' defaultValue={languages} onChange={(event, data) => updateLanguages(data.value)} /></td></tr>
                    <tr><td>Skip Uids:</td><td colSpan="2"><Input placeholder='List of uids' defaultValue={skipUids} onChange={(event, data) => updateSkipUids(data.value)} /></td></tr>
                  </tbody>
                </table>
              </Segment>
              <h3>Spec Tree</h3>
              <Segment textAlign='left'>
                <SpecTree spec={specObj} onChange={spec => updateSpec(spec ? JSON.stringify(spec, null, 2) : '')} />
              </Segment>
              <Segment textAlign='left'>
                <div>Spec JSON:</div>
                <div>
                  <TextArea placeholder='Spec' rows="10" style={{'width': '100%'}} value={spec} onChange={(event, data) => updateSpec(data.value)} />
                </div>
              </Segment>
              { null /*
                <h3>Subscriptions</h3>
                <Segment textAlign='left'>
                  <Checkbox label='CT_DAILY_LESSON'
                                      checked={this.contentTypeChecked('CT_DAILY_LESSON', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_DAILY_LESSON', options)}
                                      onChange={this.optionsChange} />
                </Segment>
                <Segment textAlign='left'>
                  <Checkbox label='CT_LECTURE'
                                      checked={this.contentTypeChecked('CT_LECTURE', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_LECTURE', options)}
                                      onChange={this.optionsChange} /><br />
                  <Checkbox label='CT_VIRTUAL_LESSON'
                                      checked={this.contentTypeChecked('CT_VIRTUAL_LESSON', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_VIRTUAL_LESSON', options)}
                                      onChange={this.optionsChange} /><br />
                  <Checkbox label='CT_WOMEN_LESSON'
                                      checked={this.contentTypeChecked('CT_WOMEN_LESSON', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_WOMEN_LESSON', options)}
                                      onChange={this.optionsChange} /><br />
                  <Checkbox label='CT_EVENT_PART'
                                      checked={this.contentTypeChecked('CT_EVENT_PART', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_EVENT_PART', options)}
                                      onChange={this.optionsChange} />
                </Segment>
                <Segment textAlign='left'>
                  <Checkbox label='CT_VIDEO_PROGRAM_CHAPTER'
                                      checked={this.contentTypeChecked('CT_VIDEO_PROGRAM_CHAPTER', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_VIDEO_PROGRAM_CHAPTER', options)}
                                      onChange={this.optionsChange} />
                              </Segment>
                              <Segment textAlign='left'>
                                  <Checkbox label='CT_BLOG_POST'
                                      checked={this.contentTypeChecked('CT_BLOG_POST', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_BLOG_POST', options)}
                                      onChange={this.optionsChange} /><br />
                                  <Checkbox label='CT_ARTICLE'
                                      checked={this.contentTypeChecked('CT_ARTICLE', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_ARTICLE', options)}
                                      onChange={this.optionsChange} /><br />
                                  <Checkbox label='CT_PUBLICATION'
                                      checked={this.contentTypeChecked('CT_PUBLICATION', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_PUBLICATION', options)}
                                      onChange={this.optionsChange} /><br />
                              </Segment>
                              <Segment textAlign='left'>
                                  <Checkbox label='CT_FRIENDS_GATHERING'
                                      checked={this.contentTypeChecked('CT_FRIENDS_GATHERING', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_FRIENDS_GATHERING', options)}
                                      onChange={this.optionsChange} /><br />
                                  <Checkbox label='CT_MEAL'
                                      checked={this.contentTypeChecked('CT_MEAL', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_MEAL', options)}
                                      onChange={this.optionsChange} />
                              </Segment>
                              <Segment textAlign='left'>
                                  <Checkbox label='CT_CLIP'
                                      checked={this.contentTypeChecked('CT_CLIP', options)}
                                      indeterminate={this.contentTypeIndeterminate('CT_CLIP', options)}
                                      onChange={this.optionsChange} />
                              </Segment>
                              <h4>Subscribe programs, clips or articles:</h4>
                              <Segment textAlign='left' style={{overflow: 'auto', maxHeight: 200, minHeight: 50}}>
                                  <Dimmer active={fetchingSubscribeCollections} inverted>
                                      <Loader inverted/>
                                  </Dimmer>
                                  {subscribeCollections.map(option => {
                                      return (<Container key={option.id}>
                                          <Checkbox label={option.name}
                                              checked={this.collectionChecked(option.content_type, option.id, options)}
                                              indeterminate={this.collectionIndeterminate(option.content_type, option.id, options)}
                                              onChange={this.optionsChange} />
                                      </Container>);
                                  })}
                              </Segment>
                              */
                }
              </Segment>
					</Grid.Column>
					<Grid.Column>
						<Segment style={{overflow: 'auto', maxHeight: '80vh'}}>
							{items.map((item) => this.renderItem(item))}
						</Segment>
						<Segment>
							<Button onClick={more}>More</Button>
							<Button onClick={reset}>Reset</Button>
              <br />
              {(error || specParseErr) && <span style={{color: 'red'}}>{error || String(specParseErr)}</span>}
						</Segment>
					</Grid.Column>
				</Grid.Row>
			</Grid>
		);
	}
}

export default Feed;


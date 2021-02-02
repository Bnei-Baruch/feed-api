import React, { 
  useEffect,
  useState, 
} from 'react';

import { 
  Button,
  Grid,
  Input,
  Segment,
  TextArea,
} from 'semantic-ui-react';

import { recommend } from './api'

import Item from './Item.js'
import SpecTree from './SpecTree.js'

import './Recommender.css';

const DEFAULT_UID = 'VDtljVgk';

const Recommender = (props) => {
  const [items, setItems] = useState([]);
  const [feed, setFeed] = useState([]);
  const [uid, setUid] = useState(new URLSearchParams(window.location.search).get('uid') || DEFAULT_UID);
  const [spec, setSpec] = useState(new URLSearchParams(window.location.search).get('spec') || '');
  const [numItems, setNumItems] = useState(20);
  const [languages, setLanguages] = useState(['he', 'en']);
  const [itemsByUid, setItemsByUid] = useState({});
  const [recommendError, setRecommendError] = useState('');
  const [specError, setSpecError] = useState('');

  const parseSpec = () => {
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
  if (!specError && specParseErr) {
    setSpecError(`Bad Spec: ${specParseErr.message}`);
  }
  if (specError && !specParseErr) {
    setSpecError('');
  }

  const recommendClicked = () => {
    setRecommendError('');
    const options = {recommend: {uid}, languages};
    if (specObj) {
      options.spec = specObj;
    }
    recommend(/*feed=*/[], itemsByUid, options, numItems).then(({feed, items, itemsByUid}) => {
      setFeed(feed);
      setItems(items);
      setItemsByUid(itemsByUid);
      setRecommendError('');
    }).catch((error) => {
      setRecommendError(error);
    });
  };
  useEffect(recommendClicked, []);

  return (
    <Grid>
      <Grid.Row>
        <Grid.Column width={9}>
          <Segment style={{'direction': 'ltr'}}>
            <h3>Context</h3>
            <Segment textAlign='left'>
              <table>
                <tbody>
                  <tr><td>UID:</td><td><Input placeholder='UID...' defaultValue={uid} onChange={(event, data) => setUid(data.value)} /></td></tr>
                  <tr><td>Num Items:</td><td><Input placeholder='Num Items to Recommend' defaultValue={numItems} onChange={(event, data) => setNumItems(Number(data.value))} /></td></tr>
                  <tr><td>Languages:</td><td><Input placeholder='List of preffered languages' defaultValue={languages.join(',')} onChange={(event, data) => setLanguages(data.value.split(',').filter(language => !!language))} /></td></tr>
                </tbody>
              </table>
            </Segment>
            <Segment textAlign='left'>
              <div>Spec Tree</div>
              <div><SpecTree spec={specObj} onChange={spec => spec ? setSpec(JSON.stringify(spec, null, 2)) : setSpec('')} /></div>
            </Segment>
            <Segment textAlign='left'>
              <div>Spec JSON:</div>
              <div>
                <TextArea placeholder='Spec' rows="10" style={{'width': '100%'}} defaultValue={spec} onChange={(event, data) => setSpec(data.value)} />
              </div>
            </Segment>
          </Segment>
        </Grid.Column>
        <Grid.Column width={7}>
          <Segment style={{overflow: 'auto', maxHeight: '80vh'}}>
            {items.map((item, index) => <Item key={index} item={item} contentItem={feed[index]} />)}
          </Segment>
          <Segment>
            <Button disabled={!!specParseErr} onClick={recommendClicked}>Recommend</Button>
            <br />
            {(recommendError || specError) && <span style={{color: 'red'}}>{specError || recommendError}</span>}
          </Segment>
        </Grid.Column>
      </Grid.Row>
    </Grid>
  );
};


export default Recommender;

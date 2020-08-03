import React, { 
  useState, 
} from 'react';

import { 
  Button,
  Grid,
  Input,
  Segment,
} from 'semantic-ui-react';

import { recommend } from './api'

import Item from './Item.js'

import './Recommender.css';

const Recommender = (props) => {
  const [items, setItems] = useState([]);
  const [feed, setFeed] = useState([]);
  const [uid, setUid] = useState('EPrcYXLv');
  const [itemsByUid, setItemsByUid] = useState({});
  const [error, setError] = useState('');

  const recommendClicked = () => {
    setError('');
    recommend(/*feed=*/[], itemsByUid, /*options=*/{recommend: {uid}}).then(({feed, items, itemsByUid}) => {
      setFeed(feed);
      setItems(items);
      setItemsByUid(itemsByUid);
    }).catch((error) => setError(error));
  };

  console.log(items);

  return (
    <Grid columns={2}>
      <Grid.Row>
        <Grid.Column>
          <Segment style={{'direction': 'ltr'}}>
            <h3>Context</h3>
            <Segment textAlign='left'>
              <Input placeholder='UID...' defaultValue='EPrcYXLv' onChange={(event, data) => setUid(data.value)} />
            </Segment>
          </Segment>
        </Grid.Column>
        <Grid.Column>
          <Segment style={{overflow: 'auto', maxHeight: '80vh'}}>
            {items.map((item, index) => <Item key={index} item={item} contentItem={feed[index]} />)}
          </Segment>
          <Segment>
            <Button onClick={recommendClicked}>Recommend</Button>
            <br />
            {error !== '' ? <span style={{color: 'red'}}>{error}</span> : null}
          </Segment>
        </Grid.Column>
      </Grid.Row>
    </Grid>
  );
};


export default Recommender;

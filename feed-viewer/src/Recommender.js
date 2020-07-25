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
  const [uid, setUid] = useState('ewJZit1t');
  const [itemsByUid, setItemsByUid] = useState({});
  const [error, setError] = useState('');

  const recommendClicked = () => {
    setError('');
    recommend(/*feed=*/[{uid, original_order: [0]}], itemsByUid, /*options=*/{}).then(({items, itemsByUid}) => {
      setItems(items);  // Fix the filter here!!! should not be!
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
              <Input placeholder='UID...' defaultValue='ewJZit1t' onChange={(event, data) => setUid(data.value)} />
            </Segment>
          </Segment>
        </Grid.Column>
        <Grid.Column>
          <Segment style={{overflow: 'auto', maxHeight: '80vh'}}>
            {items.map((item, index) => <Item key={index} item={item} />)}
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
